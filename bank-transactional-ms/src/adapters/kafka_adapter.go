package adapters

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/infra"
	"gopkg.in/mgo.v2/bson"
)

type EventDispatcherImpl struct {
	producer                        *infra.DynamicTopicProducer
	KafkaSyncRequestOutputTopicName string
}

func NewEventDispatcherImpl(ctx context.Context, brokers []string, KafkaSyncRequestOutputTopicName string) domain.EventDispatcher {

	producer := infra.NewDynamicTopicProducer(ctx, brokers)
	return &EventDispatcherImpl{
		producer:                        producer,
		KafkaSyncRequestOutputTopicName: KafkaSyncRequestOutputTopicName,
	}
}

func (dispatcher *EventDispatcherImpl) UpdateSyncRequestStatus(id bson.ObjectId, requestStatus domain.RequestStatus, Message *string) error {
	requestOutput := domain.SyncRequestResult{
		ID:            id,
		RequestStatus: requestStatus,
		Message:       Message,
		SentAt:        time.Now(),
	}
	jsonString, err := json.Marshal(requestOutput)
	if err != nil {
		return err
	}
	return dispatcher.producer.SendMessage(dispatcher.KafkaSyncRequestOutputTopicName, bson.NewObjectId().Hex(), string(jsonString))
}

func (dispatcher *EventDispatcherImpl) RequestProviderSync(providerInputTopic string, providerSyncRequest domain.ProviderSyncRequest) error {
	jsonString, err := json.Marshal(providerSyncRequest)
	if err != nil {
		return err
	}
	return dispatcher.producer.SendMessage(providerInputTopic, bson.NewObjectId().Hex(), string(jsonString))

}
