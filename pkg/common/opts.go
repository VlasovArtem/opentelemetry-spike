package common

import (
	"github.com/jessevdk/go-flags"
)

var GlobalOpts opts

type opts struct {
	ExporterType string `short:"t" long:"type" required:"true" description:"type of exporter: file or grpc" choice:"file" choice:"grpc" env:"EXPORTER_TYPE"`
	Collector    struct {
		Url      string `short:"u" long:"url" description:"collector url" env:"URL"`
		Insecure bool   `short:"i" long:"insecure" description:"collector grpc insecure" env:"INSECURE"`
	} `group:"collector" namespace:"collector" env-namespace:"COLLECTOR"`
	Kafka struct {
		Address   string `short:"a" long:"address" description:"kafka address" env:"ADDRESS"`
		Topic     string `short:"c" long:"topic" description:"kafka topic" env:"TOPIC"`
		Partition int    `short:"p" long:"partition" description:"kafka partition" env:"PARTITION"`
	} `group:"kafka" namespace:"kafka" env-namespace:"KAFKA"`
	Executor struct {
		BasePath string `short:"b" long:"base-path" description:"base path for executor" env:"BASE_PATH"`
	} `group:"executor" namespace:"executor" env-namespace:"EXECUTOR"`
}

func ParseGlobalOpts() error {
	_, err := flags.Parse(&GlobalOpts)
	return err
}
