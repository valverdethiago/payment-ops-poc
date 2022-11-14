package domain

import (
	"encoding/json"
	"fmt"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/restclient"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type OnMessageReceive func(string) error

type EventDispatcher interface {
	UpdateSyncRequestStatus(id bson.ObjectId, requestStatus RequestStatus, Message *string) error
}

type EventSubscriberService interface {
	OnMessageReceive(value string) error
}

type EventSubscriberServiceImpl struct {
	syncRequestRepository SyncRequestRepository
	accountRepository     AccountRepository
	trioClient            TrioClient
	eventDispatcher       EventDispatcher
}

func NewEventSubscriberServiceImpl(syncRequestRepository SyncRequestRepository,
	accountRepository AccountRepository,
	trioClient TrioClient,
	dispatcher EventDispatcher) EventSubscriberService {
	return &EventSubscriberServiceImpl{
		syncRequestRepository: syncRequestRepository,
		accountRepository:     accountRepository,
		trioClient:            trioClient,
		eventDispatcher:       dispatcher,
	}
}

func (subscriberService *EventSubscriberServiceImpl) OnMessageReceive(value string) error {
	providerSyncRequest, err := ParseJson(value)
	if err != nil {
		return err
	}
	fmt.Println("received at callback: ", value)
	fmt.Println("received object", providerSyncRequest.ID, providerSyncRequest.AccountId, providerSyncRequest.SyncType)
	SyncRequest := buildSyncRequest(providerSyncRequest)
	SyncRequest, err = subscriberService.persistSyncRequest(SyncRequest)
	switch SyncRequest.SyncType {
	case SyncTypeBalances:
		subscriberService.synchronizeBalances(SyncRequest)
	case SyncTypeTransactions:
		subscriberService.synchronizeTransactions(SyncRequest)
	}
	return err
}

func (subscriberService *EventSubscriberServiceImpl) persistSyncRequest(syncRequest *SyncRequest) (*SyncRequest, error) {
	return subscriberService.syncRequestRepository.Insert(syncRequest)
}

func (subscriberService *EventSubscriberServiceImpl) updateSyncRequestStatus(ID bson.ObjectId, Status RequestStatus, Message *string) (*SyncRequest, error) {
	syncRequest, err := subscriberService.syncRequestRepository.Find(ID)
	if err != nil {
		return nil, err
	}
	syncRequest.RequestStatus = Status
	syncRequest.ErrorMessage = Message
	subscriberService.eventDispatcher.UpdateSyncRequestStatus(ID, Status, Message)
	return subscriberService.syncRequestRepository.Update(syncRequest)
}

func (subscriberService *EventSubscriberServiceImpl) synchronizeBalances(Request *SyncRequest) {
	subscriberService.synchronizeWithTrio(Request, subscriberService.trioClient.FetchBalancesFromBank)
}

func (subscriberService *EventSubscriberServiceImpl) synchronizeTransactions(Request *SyncRequest) {
	subscriberService.synchronizeWithTrio(Request, subscriberService.trioClient.FetchTransactionsFromBank)
}

func (subscriberService *EventSubscriberServiceImpl) synchronizeWithTrio(Request *SyncRequest, FetchFunction FetchData) {
	accountMapping, err := subscriberService.accountRepository.FindByAccountId(Request.AccountId)
	if err != nil {
		errorMessage := fmt.Sprintf("Invalid account id %s", Request.AccountId)
		subscriberService.updateSyncRequestStatus(Request.ID, RequestStatusPending, &errorMessage)
	}
	response, err := FetchFunction(accountMapping.ProviderAccountId)
	if err != nil {
		subscriberService.updateSyncRequestStatus(Request.ID, RequestStatusPending, nil)
	}
	if response.Event.Status == restclient.FAILED {
		subscriberService.updateSyncRequestStatus(Request.ID, RequestStatusFailed, &response.Error.Message)
	}
	subscriberService.updateSyncRequestStatus(Request.ID, RequestStatusPending, nil)
}

func buildSyncRequest(providerSyncRequest ProviderSyncRequest) *SyncRequest {
	AccountId, err := parseUUID(providerSyncRequest.AccountId)
	SyncType, err := parseSyncType(providerSyncRequest.SyncType)
	if err != nil {
		log.Println("Unable to parse ProviderSyncRequest")
		return nil
	}
	return &SyncRequest{
		ID:            bson.NewObjectId(),
		RequestStatus: RequestStatusCreated,
		CreatedAt:     time.Now().Unix(),
		SyncType:      SyncType,
		AccountId:     AccountId.String(),
	}
}

func ParseJson(value string) (ProviderSyncRequest, error) {
	providerSyncRequest := ProviderSyncRequest{}
	err := json.Unmarshal([]byte(value), &providerSyncRequest)
	return providerSyncRequest, err
}
