package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type OnMessageReceive func(string) error

type EventSubscriberService interface {
	OnMessageReceive(value string) error
}

type EventDispatcher interface {
	UpdateSyncRequestStatus(id bson.ObjectId, requestStatus RequestStatus, Message *string) error
	RequestProviderSync(providerInputTopic string, providerSyncRequest ProviderSyncRequest) error
}

type EventSubscriberServiceImpl struct {
	accountService     AccountService
	syncRequestService SyncRequestService
}

func NewEventSubscriberServiceImpl(accountService AccountService, syncRequestService SyncRequestService) EventSubscriberService {
	return &EventSubscriberServiceImpl{
		accountService:     accountService,
		syncRequestService: syncRequestService,
	}
}

func (subscriberService *EventSubscriberServiceImpl) OnMessageReceive(value string) error {
	syncRequest, err := ParseJson(value)
	if err != nil {
		return err
	}
	fmt.Println("received at callback: ", value)
	fmt.Println("received object", syncRequest.ID, syncRequest.AccountId, syncRequest.SyncType)
	ID := parseBson(syncRequest.ID)
	AccountID, err := parseUUID(syncRequest.AccountId)
	if err != nil {
		subscriberService.syncRequestService.ChangeToFailingStatus(ID, fmt.Sprintf("Error parsing account  ID [%s]", syncRequest.AccountId))
		return err
	}
	account, _, configuration, err := subscriberService.accountService.FindAccountInformation(AccountID)
	if err != nil {
		subscriberService.syncRequestService.ChangeToFailingStatus(ID, fmt.Sprintf("Error on finding account with ID [%s]", syncRequest.AccountId))
		return err
	}
	if !subscriberService.accountService.IsAccountInValidState(account) {
		subscriberService.syncRequestService.ChangeToFailingStatus(ID, "Account is in invalid state")
		return errors.New(fmt.Sprintf("Account %s is in invalid state", AccountID))
	}
	subscriberService.syncRequestService.ChangeToPendingStatus(ID)
	providerSyncRequest := ProviderSyncRequest{
		AccountId: AccountID,
		SyncType:  syncRequest.SyncType,
	}
	return subscriberService.syncRequestService.RequestProviderSync(configuration.KafkaInputTopicName, providerSyncRequest)
}

func ParseJson(value string) (SyncRequest, error) {
	syncRequest := SyncRequest{}
	err := json.Unmarshal([]byte(value), &syncRequest)
	return syncRequest, err

}
