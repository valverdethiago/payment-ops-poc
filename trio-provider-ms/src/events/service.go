package events

import (
	"encoding/json"
	"fmt"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/restclient"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type OnMessageReceive func(string) error

type EventSubscriberService interface {
	OnMessageReceive(value string) error
}

type EventSubscriberServiceImpl struct {
	syncRequestRepository domain.SyncRequestRepository
	accountRepository     domain.AccountRepository
	trioClient            domain.TrioClient
	eventDispatcher       domain.EventDispatcher
}

func NewEventSubscriberServiceImpl(syncRequestRepository domain.SyncRequestRepository,
	accountRepository domain.AccountRepository,
	trioClient domain.TrioClient,
	dispatcher domain.EventDispatcher) EventSubscriberService {
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
	case domain.SyncTypeBalances:
		subscriberService.synchronizeBalances(SyncRequest)
	case domain.SyncTypeTransactions:
		subscriberService.synchronizeTransactions(SyncRequest)
	}
	return err
}

func (subscriberService *EventSubscriberServiceImpl) persistSyncRequest(syncRequest *domain.SyncRequest) (*domain.SyncRequest, error) {
	return subscriberService.syncRequestRepository.Insert(syncRequest)
}

func (subscriberService *EventSubscriberServiceImpl) updateSyncRequestStatus(ID bson.ObjectId, Status domain.RequestStatus, Message *string) (*domain.SyncRequest, error) {
	syncRequest, err := subscriberService.syncRequestRepository.Find(ID)
	if err != nil {
		return nil, err
	}
	syncRequest.RequestStatus = Status
	syncRequest.ErrorMessage = Message
	subscriberService.eventDispatcher.UpdateSyncRequestStatus(ID, Status, Message)
	return subscriberService.syncRequestRepository.Update(syncRequest)
}

func (subscriberService *EventSubscriberServiceImpl) synchronizeBalances(Request *domain.SyncRequest) {
	subscriberService.synchronizeWithTrio(Request, subscriberService.trioClient.FetchBalancesFromBank)
}

func (subscriberService *EventSubscriberServiceImpl) synchronizeTransactions(Request *domain.SyncRequest) {
	subscriberService.synchronizeWithTrio(Request, subscriberService.trioClient.FetchTransactionsFromBank)
}

func (subscriberService *EventSubscriberServiceImpl) synchronizeWithTrio(Request *domain.SyncRequest, FetchFunction domain.FetchData) {
	subscriberService.updateSyncRequestStatus(Request.ID, domain.RequestStatusPending, nil)
	account, err := subscriberService.accountRepository.FindByInternalAccountId(Request.AccountId)
	if err != nil {
		errorMessage := fmt.Sprintf("Invalid account id %s", Request.AccountId)
		subscriberService.updateSyncRequestStatus(Request.ID, domain.RequestStatusPending, &errorMessage)
	}
	response, err := FetchFunction(*account)
	if err != nil {
		subscriberService.updateSyncRequestStatus(Request.ID, domain.RequestStatusPending, nil)
	}
	if response.Data.Status == string(restclient.FailedFetchRequestStatus) {
		subscriberService.updateSyncRequestStatus(Request.ID, domain.RequestStatusFailed, nil)

	}
	subscriberService.updateSyncRequestStatus(Request.ID, domain.RequestStatusPending, nil)
}

func buildSyncRequest(providerSyncRequest domain.ProviderSyncRequest) *domain.SyncRequest {
	AccountId, err := domain.ParseUUID(providerSyncRequest.AccountId)
	SyncType, err := domain.ParseSyncType(providerSyncRequest.SyncType)
	if err != nil {
		log.Println("Unable to parse ProviderSyncRequest")
		return nil
	}
	return &domain.SyncRequest{
		ID:            bson.NewObjectId(),
		RequestStatus: domain.RequestStatusCreated,
		CreatedAt:     time.Now().Unix(),
		SyncType:      SyncType,
		AccountId:     AccountId.String(),
	}
}

func ParseJson(value string) (domain.ProviderSyncRequest, error) {
	providerSyncRequest := domain.ProviderSyncRequest{}
	err := json.Unmarshal([]byte(value), &providerSyncRequest)
	return providerSyncRequest, err
}
