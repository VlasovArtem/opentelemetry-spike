package env

import (
	"os"
	"spike-go-opentelemetry-logging/pkg/otel/internal/global"
	"strconv"
)

// Environment variable names.
const (
	// BatchLogRecordProcessorScheduleDelayKey is the delay interval between two
	// consecutive exports (i.e. 5000).
	BatchLogRecordProcessorScheduleDelayKey = "OTEL_BLRP_SCHEDULE_DELAY"
	// BatchLogRecordProcessorExportTimeoutKey is the maximum allowed time to
	// export data (i.e. 3000).
	BatchLogRecordProcessorExportTimeoutKey = "OTEL_BLRP_EXPORT_TIMEOUT"
	// BatchLogRecordProcessorMaxQueueSizeKey is the maximum queue size (i.e. 2048).
	BatchLogRecordProcessorMaxQueueSizeKey = "OTEL_BLRP_MAX_QUEUE_SIZE"
	// BatchLogRecordProcessorMaxExportBatchSizeKey is the maximum batch size (i.e.
	// 512). Note: it must be less than or equal to
	// EnvBatchLogRecordProcessorMaxQueueSize.
	BatchLogRecordProcessorMaxExportBatchSizeKey = "OTEL_BLRP_MAX_EXPORT_BATCH_SIZE"

	// AttributeValueLengthKey is the maximum allowed attribute value size.
	AttributeValueLengthKey = "OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT"

	// AttributeCountKey is the maximum allowed span attribute count.
	AttributeCountKey = "OTEL_ATTRIBUTE_COUNT_LIMIT"

	// SpanAttributeValueLengthKey is the maximum allowed attribute value size
	// for a span.
	SpanAttributeValueLengthKey = "OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT"

	// SpanAttributeCountKey is the maximum allowed span attribute count for a
	// span.
	SpanAttributeCountKey = "OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT"

	// SpanEventCountKey is the maximum allowed span event count.
	SpanEventCountKey = "OTEL_SPAN_EVENT_COUNT_LIMIT"

	// SpanEventAttributeCountKey is the maximum allowed attribute per span
	// event count.
	SpanEventAttributeCountKey = "OTEL_EVENT_ATTRIBUTE_COUNT_LIMIT"

	// SpanLinkCountKey is the maximum allowed span link count.
	SpanLinkCountKey = "OTEL_SPAN_LINK_COUNT_LIMIT"

	// SpanLinkAttributeCountKey is the maximum allowed attribute per span
	// link count.
	SpanLinkAttributeCountKey = "OTEL_LINK_ATTRIBUTE_COUNT_LIMIT"
)

// firstInt returns the value of the first matching environment variable from
// keys. If the value is not an integer or no match is found, defaultValue is
// returned.
func firstInt(defaultValue int, keys ...string) int {
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if !ok {
			continue
		}

		intValue, err := strconv.Atoi(value)
		if err != nil {
			global.Info("Got invalid value, number value expected.", key, value)
			return defaultValue
		}

		return intValue
	}

	return defaultValue
}

// IntEnvOr returns the int value of the environment variable with name key if
// it exists and the value is an int. Otherwise, defaultValue is returned.
func IntEnvOr(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		global.Info("Got invalid value, number value expected.", key, value)
		return defaultValue
	}

	return intValue
}

// BatchLogRecordProcessorScheduleDelay returns the environment variable value for
// the OTEL_BSP_SCHEDULE_DELAY key if it exists, otherwise defaultValue is
// returned.
func BatchLogRecordProcessorScheduleDelay(defaultValue int) int {
	return IntEnvOr(BatchLogRecordProcessorScheduleDelayKey, defaultValue)
}

// BatchLogRecordProcessorExportTimeout returns the environment variable value for
// the OTEL_BSP_EXPORT_TIMEOUT key if it exists, otherwise defaultValue is
// returned.
func BatchLogRecordProcessorExportTimeout(defaultValue int) int {
	return IntEnvOr(BatchLogRecordProcessorExportTimeoutKey, defaultValue)
}

// BatchLogRecordProcessorMaxQueueSize returns the environment variable value for
// the OTEL_BSP_MAX_QUEUE_SIZE key if it exists, otherwise defaultValue is
// returned.
func BatchLogRecordProcessorMaxQueueSize(defaultValue int) int {
	return IntEnvOr(BatchLogRecordProcessorMaxQueueSizeKey, defaultValue)
}

// BatchLogRecordProcessorMaxExportBatchSize returns the environment variable value for
// the OTEL_BSP_MAX_EXPORT_BATCH_SIZE key if it exists, otherwise defaultValue
// is returned.
func BatchLogRecordProcessorMaxExportBatchSize(defaultValue int) int {
	return IntEnvOr(BatchLogRecordProcessorMaxExportBatchSizeKey, defaultValue)
}

// SpanAttributeValueLength returns the environment variable value for the
// OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT key if it exists. Otherwise, the
// environment variable value for OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT is
// returned or defaultValue if that is not set.
func SpanAttributeValueLength(defaultValue int) int {
	return firstInt(defaultValue, SpanAttributeValueLengthKey, AttributeValueLengthKey)
}

// SpanAttributeCount returns the environment variable value for the
// OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT key if it exists. Otherwise, the
// environment variable value for OTEL_ATTRIBUTE_COUNT_LIMIT is returned or
// defaultValue if that is not set.
func SpanAttributeCount(defaultValue int) int {
	return firstInt(defaultValue, SpanAttributeCountKey, AttributeCountKey)
}

// SpanEventCount returns the environment variable value for the
// OTEL_SPAN_EVENT_COUNT_LIMIT key if it exists, otherwise defaultValue is
// returned.
func SpanEventCount(defaultValue int) int {
	return IntEnvOr(SpanEventCountKey, defaultValue)
}

// SpanEventAttributeCount returns the environment variable value for the
// OTEL_EVENT_ATTRIBUTE_COUNT_LIMIT key if it exists, otherwise defaultValue
// is returned.
func SpanEventAttributeCount(defaultValue int) int {
	return IntEnvOr(SpanEventAttributeCountKey, defaultValue)
}

// SpanLinkCount returns the environment variable value for the
// OTEL_SPAN_LINK_COUNT_LIMIT key if it exists, otherwise defaultValue is
// returned.
func SpanLinkCount(defaultValue int) int {
	return IntEnvOr(SpanLinkCountKey, defaultValue)
}

// SpanLinkAttributeCount returns the environment variable value for the
// OTEL_LINK_ATTRIBUTE_COUNT_LIMIT key if it exists, otherwise defaultValue is
// returned.
func SpanLinkAttributeCount(defaultValue int) int {
	return IntEnvOr(SpanLinkAttributeCountKey, defaultValue)
}
