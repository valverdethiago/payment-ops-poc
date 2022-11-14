package adapters

import (
	"context"
	"encoding/json"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/events"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/infra"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type EventDispatcherImpl struct {
	producer                         *infra.DynamicTopicProducer
	KafkaSyncRequestOutputTopicName  string
	KafkaBalanceUpdateTopicName      string
	KafkaTransactionsUpdateTopicName string
}

func NewEventDispatcherImpl(ctx context.Context, brokers []string, KafkaSyncRequestOutputTopicName string,
	KafkaBalanceUpdateTopicName string, KafkaTransactionsUpdateTopicName string) domain.EventDispatcher {

	producer := infra.NewDynamicTopicProducer(ctx, brokers)
	return &EventDispatcherImpl{
		producer:                         producer,
		KafkaSyncRequestOutputTopicName:  KafkaSyncRequestOutputTopicName,
		KafkaBalanceUpdateTopicName:      KafkaBalanceUpdateTopicName,
		KafkaTransactionsUpdateTopicName: KafkaTransactionsUpdateTopicName,
	}
}

func (dispatcher EventDispatcherImpl) UpdateSyncRequestStatus(id string,
	requestStatus domain.RequestStatus, Message *string) error {
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
	return dispatcher.producer.SendMessage(dispatcher.KafkaSyncRequestOutputTopicName,
		bson.NewObjectId().Hex(), string(jsonString))
}

func (dispatcher EventDispatcherImpl) TriggerBalanceUpdateEvent(accountId string, balance float64, currency string) error {
	eventPayload := events.BalanceUpdateEvent{
		AccountID: accountId,
		Amount:    balance,
		Currency:  currency,
	}
	jsonString, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}
	return dispatcher.producer.SendMessage(dispatcher.KafkaBalanceUpdateTopicName,
		bson.NewObjectId().Hex(), string(jsonString))
}

func (dispatcher EventDispatcherImpl) TriggerTransactionsUpdateEvent(accountId string,
	transactions []domain.Transaction) error {
	eventPayload := events.TransactionsUpdateEvent{
		AccountId:    accountId,
		Transactions: transactions,
	}
	jsonString, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}
	return dispatcher.producer.SendMessage(dispatcher.KafkaTransactionsUpdateTopicName,
		bson.NewObjectId().Hex(), string(jsonString))
}
