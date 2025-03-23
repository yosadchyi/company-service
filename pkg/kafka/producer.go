package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func InitKafka(brokers []string, topic string) {
	writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func PublishEvent(eventType string, payload any) {
	event := map[string]any{
		"event_type": eventType,
		"payload":    payload,
	}
	data, err := json.Marshal(event)
	if err != nil {
		log.Println("Kafka marshal error:", err)
		return
	}
	if err := writer.WriteMessages(context.Background(), kafka.Message{Value: data}); err != nil {
		log.Println("Kafka write error:", err)
	}
}
