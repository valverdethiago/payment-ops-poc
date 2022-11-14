package infra

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	ctx     context.Context
	brokers []string
	writter kafka.Writer
}

func NewProducer(ctx context.Context, brokers []string, topic string) *Producer {

	logger := log.New(os.Stdout, "kafka writer: ", 0)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
		Logger:  logger,
	})

	return &Producer{
		ctx:     ctx,
		brokers: brokers,
		writter: *writer,
	}

}

func (producer *Producer) SendMessage(key string, message string) error {
	return producer.writter.WriteMessages(producer.ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(message),
	})
}
