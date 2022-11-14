package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type OnMessageReceive func(string) error

type EventSubscriberService interface {
	OnReceiveSyncRequest(value string) error
	OnReceiveBalanceUpdate(value string) error
	OnReceiveTransactionsUpdate(value string) error
}

type EventDispatcher interface {
	UpdateSyncRequestStatus(id bson.ObjectId, requestStatus RequestStatus, Message *string) error
	RequestProviderSync(providerInputTopic string, providerSyncRequest ProviderSyncRequest) error
}

type EventSubscriberServiceImpl struct {
	accountService        AccountService
	accountBalanceService AccountBalanceService
	transactionService    TransactionService
	syncRequestService    SyncRequestService
}

func NewEventSubscriberServiceImpl(accountService AccountService, accountBalanceService AccountBalanceService,
	transactionService TransactionService, syncRequestService SyncRequestService) EventSubscriberService {
	return &EventSubscriberServiceImpl{
		accountService:        accountService,
		accountBalanceService: accountBalanceService,
		transactionService:    transactionService,
		syncRequestService:    syncRequestService,
	}
}

func (subscriberService *EventSubscriberServiceImpl) OnReceiveSyncRequest(value string) error {
	syncRequest, err := ParseSyncRequestJson(value)
	if err != nil {
		return err
	}
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
		ID:            ID,
		AccountId:     AccountID,
		SyncType:      syncRequest.SyncType,
		RequestStatus: RequestStatusCreated,
		CreatedAt:     time.Now(),
	}
	return subscriberService.syncRequestService.RequestProviderSync(configuration.KafkaInputTopicName, providerSyncRequest)
}

func (subscriberService *EventSubscriberServiceImpl) OnReceiveBalanceUpdate(value string) error {
	payload, err := ParseBalanceUpdateEventJson(value)
	if err != nil {
		return err
	}
	AccountID, err := parseUUID(payload.AccountID)
	if err != nil {
		return err
	}
	account, _, _, err := subscriberService.accountService.FindAccountInformation(AccountID)
	if err != nil {
		return err
	}
	if !subscriberService.accountService.IsAccountInValidState(account) {
		return errors.New(fmt.Sprintf("Account %s is in invalid state", AccountID))
	}
	_, err = subscriberService.accountBalanceService.UpdateAccountBalance(AccountID, payload.Balance, payload.Currency)
	return err
}

func (subscriberService *EventSubscriberServiceImpl) OnReceiveTransactionsUpdate(value string) error {
	payload, err := ParseTransactionsUpdateEventJson(value)
	if err != nil {
		return err
	}
	AccountID, err := parseUUID(payload.AccountId)
	if err != nil {
		return err
	}
	account, _, _, err := subscriberService.accountService.FindAccountInformation(AccountID)
	if err != nil {
		return err
	}
	if !subscriberService.accountService.IsAccountInValidState(account) {
		return errors.New(fmt.Sprintf("Account %s is in invalid state", AccountID))
	}
	transactions, err := ConvertTransactionsFromRestPayload(payload)
	if err != nil {
		return err
	}
	_, err = subscriberService.transactionService.InsertTransactions(transactions)
	return err
}

func ParseSyncRequestJson(value string) (SyncRequest, error) {
	syncRequest := SyncRequest{}
	err := json.Unmarshal([]byte(value), &syncRequest)
	return syncRequest, err
}

func ParseBalanceUpdateEventJson(value string) (BalanceUpdateEvent, error) {
	event := BalanceUpdateEvent{}
	err := json.Unmarshal([]byte(value), &event)
	return event, err
}

func ParseTransactionsUpdateEventJson(value string) (TransactionsUpdateEvent, error) {
	event := TransactionsUpdateEvent{}
	err := json.Unmarshal([]byte(value), &event)
	return event, err
}
