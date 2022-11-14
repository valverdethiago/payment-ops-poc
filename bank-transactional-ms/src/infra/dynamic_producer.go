package infra

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type DynamicTopicProducer struct {
	ctx     context.Context
	brokers []string
}

func NewDynamicTopicProducer(ctx context.Context, brokers []string) *DynamicTopicProducer {
	return &DynamicTopicProducer{
		ctx:     ctx,
		brokers: brokers,
	}

}

func (producer *DynamicTopicProducer) SendMessage(topic string, key string, message string) error {
	logger := log.New(os.Stdout, "kafka writer: ", 0)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: producer.brokers,
		Topic:   topic,
		Logger:  logger,
	})
	return writer.WriteMessages(producer.ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(message),
	})
}
