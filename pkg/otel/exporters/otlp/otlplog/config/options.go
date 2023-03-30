package config

import (
	"crypto/tls"
	"fmt"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"path"
	"spike-go-opentelemetry-logging/pkg/otel/exporters/otlp/otlplog/retry"
	"strings"
	"time"
)

const (
	// DefaultCollectorGRPCPort is the default gRPC port of the collector.
	DefaultCollectorGRPCPort uint16 = 4317
	// DefaultCollectorHTTPPort is the default HTTP port of the collector.
	DefaultCollectorHTTPPort uint16 = 4318
	// DefaultCollectorHost is the host address the Exporter will attempt
	// connect to if no collector address is provided.
	DefaultCollectorHost string = "localhost"
	// DefaultLogRecordsPath is a default URL path for endpoint that
	// receives spans.
	DefaultLogRecordsPath string = "/v1/logs"
	// DefaultTimeout is a default max waiting time for the backend to process
	// each span batch.
	DefaultTimeout time.Duration = 10 * time.Second
)

type (
	SignalConfig struct {
		Endpoint    string
		Insecure    bool
		TLSCfg      *tls.Config
		Headers     map[string]string
		Compression Compression
		Timeout     time.Duration
		URLPath     string

		// gRPC configurations
		GRPCCredentials credentials.TransportCredentials
	}

	Config struct {
		// Signal specific configurations
		LogRecords SignalConfig

		RetryConfig retry.Config

		// gRPC configurations
		ReconnectionPeriod time.Duration
		ServiceConfig      string
		DialOptions        []grpc.DialOption
		GRPCConn           *grpc.ClientConn
	}
)

// NewHTTPConfig returns a new Config with all settings applied from opts and
// any unset setting using the default HTTP config values.
func NewHTTPConfig(opts ...HTTPOption) Config {
	cfg := Config{
		LogRecords: SignalConfig{
			Endpoint:    fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorHTTPPort),
			URLPath:     DefaultLogRecordsPath,
			Compression: NoCompression,
			Timeout:     DefaultTimeout,
		},
		RetryConfig: retry.DefaultConfig,
	}
	cfg = ApplyHTTPEnvConfigs(cfg)
	for _, opt := range opts {
		cfg = opt.ApplyHTTPOption(cfg)
	}
	cfg.LogRecords.URLPath = cleanPath(cfg.LogRecords.URLPath, DefaultLogRecordsPath)
	return cfg
}

// NewGRPCConfig returns a new Config with all settings applied from opts and
// any unset setting using the default gRPC config values.
func NewGRPCConfig(opts ...GRPCOption) Config {
	cfg := Config{
		LogRecords: SignalConfig{
			Endpoint:    fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorGRPCPort),
			URLPath:     DefaultLogRecordsPath,
			Compression: NoCompression,
			Timeout:     DefaultTimeout,
		},
		RetryConfig: retry.DefaultConfig,
		DialOptions: []grpc.DialOption{grpc.WithUserAgent(getUserAgentHeader())},
	}
	cfg = ApplyGRPCEnvConfigs(cfg)
	for _, opt := range opts {
		cfg = opt.ApplyGRPCOption(cfg)
	}

	if cfg.ServiceConfig != "" {
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithDefaultServiceConfig(cfg.ServiceConfig))
	}
	// Priroritize GRPCCredentials over Insecure (passing both is an error).
	if cfg.LogRecords.GRPCCredentials != nil {
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithTransportCredentials(cfg.LogRecords.GRPCCredentials))
	} else if cfg.LogRecords.Insecure {
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Default to using the host's root CA.
		creds := credentials.NewTLS(nil)
		cfg.LogRecords.GRPCCredentials = creds
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithTransportCredentials(creds))
	}
	if cfg.LogRecords.Compression == GzipCompression {
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	}
	if len(cfg.DialOptions) != 0 {
		cfg.DialOptions = append(cfg.DialOptions, cfg.DialOptions...)
	}
	if cfg.ReconnectionPeriod != 0 {
		p := grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: cfg.ReconnectionPeriod,
		}
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithConnectParams(p))
	}

	return cfg
}

type (
	// GenericOption applies an option to the HTTP or gRPC driver.
	GenericOption interface {
		ApplyHTTPOption(Config) Config
		ApplyGRPCOption(Config) Config

		// A private method to prevent users implementing the
		// interface and so future additions to it will not
		// violate compatibility.
		private()
	}

	// HTTPOption applies an option to the HTTP driver.
	HTTPOption interface {
		ApplyHTTPOption(Config) Config

		// A private method to prevent users implementing the
		// interface and so future additions to it will not
		// violate compatibility.
		private()
	}

	// GRPCOption applies an option to the gRPC driver.
	GRPCOption interface {
		ApplyGRPCOption(Config) Config

		// A private method to prevent users implementing the
		// interface and so future additions to it will not
		// violate compatibility.
		private()
	}
)

// genericOption is an option that applies the same logic
// for both gRPC and HTTP.
type genericOption struct {
	fn func(Config) Config
}

func (g *genericOption) ApplyGRPCOption(cfg Config) Config {
	return g.fn(cfg)
}

func (g *genericOption) ApplyHTTPOption(cfg Config) Config {
	return g.fn(cfg)
}

func (genericOption) private() {}

func newGenericOption(fn func(cfg Config) Config) GenericOption {
	return &genericOption{fn: fn}
}

// splitOption is an option that applies different logics
// for gRPC and HTTP.
type splitOption struct {
	httpFn func(Config) Config
	grpcFn func(Config) Config
}

func (g *splitOption) ApplyGRPCOption(cfg Config) Config {
	return g.grpcFn(cfg)
}

func (g *splitOption) ApplyHTTPOption(cfg Config) Config {
	return g.httpFn(cfg)
}

func (splitOption) private() {}

func newSplitOption(httpFn func(cfg Config) Config, grpcFn func(cfg Config) Config) GenericOption {
	return &splitOption{httpFn: httpFn, grpcFn: grpcFn}
}

// httpOption is an option that is only applied to the HTTP driver.
type httpOption struct {
	fn func(Config) Config
}

func (h *httpOption) ApplyHTTPOption(cfg Config) Config {
	return h.fn(cfg)
}

func (httpOption) private() {}

func NewHTTPOption(fn func(cfg Config) Config) HTTPOption {
	return &httpOption{fn: fn}
}

// grpcOption is an option that is only applied to the gRPC driver.
type grpcOption struct {
	fn func(Config) Config
}

func (h *grpcOption) ApplyGRPCOption(cfg Config) Config {
	return h.fn(cfg)
}

func (grpcOption) private() {}

func NewGRPCOption(fn func(cfg Config) Config) GRPCOption {
	return &grpcOption{fn: fn}
}

// Generic Options

func WithEndpoint(endpoint string) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.LogRecords.Endpoint = endpoint
		return cfg
	})
}

func WithCompression(compression Compression) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.LogRecords.Compression = compression
		return cfg
	})
}

func WithURLPath(urlPath string) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.LogRecords.URLPath = urlPath
		return cfg
	})
}

func WithRetry(rc retry.Config) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.RetryConfig = rc
		return cfg
	})
}

func WithTLSClientConfig(tlsCfg *tls.Config) GenericOption {
	return newSplitOption(func(cfg Config) Config {
		cfg.LogRecords.TLSCfg = tlsCfg.Clone()
		return cfg
	}, func(cfg Config) Config {
		cfg.LogRecords.GRPCCredentials = credentials.NewTLS(tlsCfg)
		return cfg
	})
}

func WithInsecure() GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.LogRecords.Insecure = true
		return cfg
	})
}

func WithSecure() GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.LogRecords.Insecure = false
		return cfg
	})
}

func WithHeaders(headers map[string]string) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.LogRecords.Headers = headers
		return cfg
	})
}

func WithTimeout(duration time.Duration) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.LogRecords.Timeout = duration
		return cfg
	})
}

// Compression describes the compression used for payloads sent to the
// collector.
type Compression int

const (
	// NoCompression tells the driver to send payloads without
	// compression.
	NoCompression Compression = iota
	// GzipCompression tells the driver to send payloads after
	// compressing them with gzip.
	GzipCompression
)

// Marshaler describes the kind of message format sent to the collector.
type Marshaler int

const (
	// MarshalProto tells the driver to send using the protobuf binary format.
	MarshalProto Marshaler = iota
	// MarshalJSON tells the driver to send using json format.
	MarshalJSON
)

// CleanPath returns a path with all spaces trimmed and all redundancies removed. If urlPath is empty or cleaning it results in an empty string, defaultPath is returned instead.
func cleanPath(urlPath string, defaultPath string) string {
	tmp := path.Clean(strings.TrimSpace(urlPath))
	if tmp == "." {
		return defaultPath
	}
	if !path.IsAbs(tmp) {
		tmp = fmt.Sprintf("/%s", tmp)
	}
	return tmp
}

func getUserAgentHeader() string {
	return "OTel OTLP Exporter Go/" + otel.Version()
}
