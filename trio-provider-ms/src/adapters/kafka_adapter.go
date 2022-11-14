package adapters

import (
	"context"
	"encoding/json"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/infra"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type EventDispatcherImpl struct {
	producer                        *infra.Producer
	KafkaSyncRequestOutputTopicName string
}

func NewEventDispatcherImpl(ctx context.Context, brokers []string, KafkaSyncRequestOutputTopicName string) domain.EventDispatcher {

	producer := infra.NewProducer(ctx, brokers, KafkaSyncRequestOutputTopicName)
	return &EventDispatcherImpl{
		producer:                        producer,
		KafkaSyncRequestOutputTopicName: KafkaSyncRequestOutputTopicName,
	}
}

func (dispatcher EventDispatcherImpl) UpdateSyncRequestStatus(id bson.ObjectId, requestStatus domain.RequestStatus, Message *string) error {
	requestOutput := domain.SyncRequestResult{
		ID:            id,
		RequestStatus: requestStatus,
		Message:       Message,
		SentAt:        time.Now().Unix(),
	}
	jsonString, err := json.Marshal(requestOutput)
	if err != nil {
		return err
	}
	return dispatcher.producer.SendMessage(bson.NewObjectId().Hex(), string(jsonString))
}
