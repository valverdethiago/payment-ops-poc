package adapters

import (
	"fmt"

	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/domain"
	"github.com/Shopify/sarama"
)

type EventNotifierServiceImpl struct {
	syncProducer sarama.SyncProducer
	topic        string
}

func NewEventNotifierServiceImpl(syncProducer sarama.SyncProducer, topic string) domain.EventNotifierService {
	return &EventNotifierServiceImpl{
		syncProducer: syncProducer,
		topic:        topic,
	}
}

func (publisherService *EventNotifierServiceImpl) Send(value []byte) error {
	// defer publisherService.syncProducer.Close()
	msg := &sarama.ProducerMessage{
		Topic: publisherService.topic,
		Value: sarama.StringEncoder(value),
	}
	partition, offset, err := publisherService.syncProducer.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)", publisherService.topic, partition, offset)
	return nil
}
