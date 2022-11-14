package domain

import (
	"encoding/json"
	"fmt"
)

type OnMessageReceive func(string) error

type EventSubscriberService interface {
	OnMessageReceive(value string)
}

type EventSubscriberServiceImpl struct{}

func NewEventSubscriberServiceImpl() EventSubscriberService {
	return &EventSubscriberServiceImpl{}
}

func (subscriberService *EventSubscriberServiceImpl) OnMessageReceive(value string) error {
	syncRequest, err := subscriberService.parseJson(value)
	if err != nil {
		return err
	}
	account :=
		fmt.Println("received at callback: ", value)
	fmt.Println("received object", syncRequest.ID, syncRequest.AccountId, syncRequest.SyncType)
}

func (subscriberService *EventSubscriberServiceImpl) parseJson(value string) (SyncRequest, error) {
	syncRequest := SyncRequest{}
	err := json.Unmarshal([]byte(value), &syncRequest)
	return syncRequest, err

}
