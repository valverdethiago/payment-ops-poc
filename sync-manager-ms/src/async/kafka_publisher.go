package async

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type SyncRequestPublisherServiceImpl struct {
	syncProducer     sarama.SyncProducer
	syncRequestTopic string
}

type SyncRequestPublisherService interface {
	Send(value []byte) error
}

func NewSyncRequestPublisherServiceImpl(syncProducer sarama.SyncProducer, syncRequestTopic string) SyncRequestPublisherService {
	return &SyncRequestPublisherServiceImpl{
		syncProducer:     syncProducer,
		syncRequestTopic: syncRequestTopic,
	}
}

func (publisherService *SyncRequestPublisherServiceImpl) Send(value []byte) error {
	// defer publisherService.syncProducer.Close()
	msg := &sarama.ProducerMessage{
		Topic: publisherService.syncRequestTopic,
		Value: sarama.StringEncoder(value),
	}
	partition, offset, err := publisherService.syncProducer.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)", publisherService.syncRequestTopic, partition, offset)
	return nil
}
