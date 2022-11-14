package adapters

import (
	"encoding/json"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/infra"
	"gopkg.in/mgo.v2/bson"
)

type EventDispatcherImpl struct {
	producer infra.Producer
}

func NewEventDispatcherImpl(producer infra.Producer) domain.EventDispatcher {
	return &EventDispatcherImpl{
		producer: producer,
	}
}

func (e EventDispatcherImpl) CreateSyncRequest(syncRequest domain.SyncRequest) error {
	jsonString, err := json.Marshal(syncRequest)
	if err != nil {
		return err
	}
	return e.producer.SendMessage(bson.NewObjectId().Hex(), string(jsonString))
}
