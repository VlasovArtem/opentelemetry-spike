package opentelemetry

import (
	"github.com/segmentio/kafka-go"
)

type MessageCarrier struct {
	msg *kafka.Message
}

func NewMessageCarrier(msg *kafka.Message) MessageCarrier {
	return MessageCarrier{msg: msg}
}

func (m MessageCarrier) Get(key string) string {
	for _, header := range m.msg.Headers {
		if header.Key == key {
			return string(header.Value)
		}
	}
	return ""
}

func (m MessageCarrier) Set(key string, value string) {
	m.msg.Headers = append(m.msg.Headers, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (m MessageCarrier) Keys() []string {
	out := make([]string, len(m.msg.Headers))
	for _, header := range m.msg.Headers {
		out = append(out, header.Key)
	}
	return out
}
