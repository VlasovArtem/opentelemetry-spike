package otlplog

import (
	"context"
	logpb "go.opentelemetry.io/proto/otlp/logs/v1"
)

type Client interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	UploadLogRecords(ctx context.Context, records []*logpb.ResourceLogs) error
}
