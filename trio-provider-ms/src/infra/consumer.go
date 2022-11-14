package infra

import (
	"context"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/events"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	ctx              context.Context
	brokers          []string
	reader           kafka.Reader
	onMessageReceive events.OnMessageReceive
}

func NewConsumer(ctx context.Context, brokers []string, topic string, groupId string) *Consumer {
	logger := log.New(os.Stdout, "kafka reader: ", 0)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupId,
		Logger:  logger,
	})

	return &Consumer{
		ctx:     ctx,
		brokers: brokers,
		reader:  *reader,
	}

}

func (consumer *Consumer) StartReading(onMessageReceive events.OnMessageReceive) {
	for {
		msg, err := consumer.reader.ReadMessage(consumer.ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		onMessageReceive(string(msg.Value))
	}
}
