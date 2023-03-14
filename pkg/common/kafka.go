package common

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func NewKafkaConnection(network string, address string, topic string, partition int) *kafka.Conn {
	conn, err := kafka.DialLeader(context.Background(), network, address, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	return conn
}

func WriteMessages(conn *kafka.Conn, messages []kafka.Message) {
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := conn.WriteMessages(messages...)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}

//func test() {
//	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
//	if err != nil {
//		log.Fatal("failed to dial leader:", err)
//	}
//	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
//	_, err = conn.WriteMessages(
//		kafka.Message{Value: []byte("one!")},
//		kafka.Message{Value: []byte("two!")},
//		kafka.Message{Value: []byte("three!")},
//	)
//	if err != nil {
//		log.Fatal("failed to write messages:", err)
//	}
//
//	if err := conn.Close(); err != nil {
//		log.Fatal("failed to close writer:", err)
//	}
//}
