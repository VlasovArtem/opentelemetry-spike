package log

import (
	"context"
	"spike-go-opentelemetry-logging/pkg/otel/log"
)

type LogExporter interface {
	ExportLogRecords(ctx context.Context, records []log.ReadableLogRecord) error
	Shutdown(ctx context.Context) error
}
