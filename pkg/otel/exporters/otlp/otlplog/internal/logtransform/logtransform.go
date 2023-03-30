package logtransform

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logpb "go.opentelemetry.io/proto/otlp/logs/v1"
	"spike-go-opentelemetry-logging/pkg/otel/sdk/log"
)

func LogRecords(records []log.ReadableLogRecord) []*logpb.ResourceLogs {
	if len(records) == 0 {
		return nil
	}

	rlm := make(map[attribute.Distinct]*logpb.ResourceLogs)

	type key struct {
		r  attribute.Distinct
		is instrumentation.Scope
	}
	slm := make(map[key]*logpb.ScopeLogs)

	var resources int
	for _, lr := range records {
		if lr == nil {
			continue
		}

		rKey := lr.Resource().Equivalent()
		k := key{
			r:  rKey,
			is: lr.InstrumentationScope(),
		}
		scopeLogs, iOk := slm[k]
		if !iOk {
			// Either the resource or instrumentation scope were unknown.
			scopeLogs = &logpb.ScopeLogs{
				Scope:      InstrumentationScope(lr.InstrumentationScope()),
				LogRecords: []*logpb.LogRecord{},
				SchemaUrl:  lr.InstrumentationScope().SchemaURL,
			}
		}
		scopeLogs.LogRecords = append(scopeLogs.LogRecords, logRecord(lr))
		slm[k] = scopeLogs

		rs, rOk := rlm[rKey]
		if !rOk {
			resources++
			// The resource was unknown.
			rs = &logpb.ResourceLogs{
				Resource:  Resource(lr.Resource()),
				ScopeLogs: []*logpb.ScopeLogs{scopeLogs},
				SchemaUrl: lr.Resource().SchemaURL(),
			}
			rlm[rKey] = rs
			continue
		}

		// The resource has been seen before. Check if the instrumentation
		// library lookup was unknown because if so we need to add it to the
		// ResourceSpans. Otherwise, the instrumentation library has already
		// been seen and the append we did above will be included it in the
		// ScopeLogs reference.
		if !iOk {
			rs.ScopeLogs = append(rs.ScopeLogs, scopeLogs)
		}
	}

	// Transform the categorized map into a slice
	rss := make([]*logpb.ResourceLogs, 0, resources)
	for _, rs := range rlm {
		rss = append(rss, rs)
	}
	return rss
}

func logRecord(record log.ReadableLogRecord) *logpb.LogRecord {
	tid := record.Context().TraceID()
	sid := record.Context().SpanID()

	return &logpb.LogRecord{
		TimeUnixNano:         uint64(record.Timestamp().UnixNano()),
		ObservedTimeUnixNano: uint64(record.ObservedTimestamp().UnixNano()),
		SeverityNumber:       logpb.SeverityNumber(record.SeverityNumber()),
		SeverityText:         record.SeverityText(),
		Body:                 value(record.Body()),
		Attributes:           attributes(record.Attributes()),
		TraceId:              tid[:],
		SpanId:               sid[:],
		Flags:                record.Context().TraceFlags(),
	}
}

func value(v attribute.Value) *commonpb.AnyValue {
	av := new(commonpb.AnyValue)
	switch v.Type() {
	case attribute.BOOL:
		av.Value = &commonpb.AnyValue_BoolValue{
			BoolValue: v.AsBool(),
		}
	case attribute.BOOLSLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: boolSliceValues(v.AsBoolSlice()),
			},
		}
	case attribute.INT64:
		av.Value = &commonpb.AnyValue_IntValue{
			IntValue: v.AsInt64(),
		}
	case attribute.INT64SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: int64SliceValues(v.AsInt64Slice()),
			},
		}
	case attribute.FLOAT64:
		av.Value = &commonpb.AnyValue_DoubleValue{
			DoubleValue: v.AsFloat64(),
		}
	case attribute.FLOAT64SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: float64SliceValues(v.AsFloat64Slice()),
			},
		}
	case attribute.STRING:
		av.Value = &commonpb.AnyValue_StringValue{
			StringValue: v.AsString(),
		}
	case attribute.STRINGSLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: stringSliceValues(v.AsStringSlice()),
			},
		}
	default:
		av.Value = &commonpb.AnyValue_StringValue{
			StringValue: "INVALID",
		}
	}
	return av
}

func attributes(attributes []attribute.KeyValue) []*commonpb.KeyValue {
	if len(attributes) == 0 {
		return nil
	}

	kvs := make([]*commonpb.KeyValue, len(attributes))

	for _, kv := range attributes {
		kvs = append(kvs, &commonpb.KeyValue{
			Key:   string(kv.Key),
			Value: value(kv.Value),
		})
	}

	return kvs
}
