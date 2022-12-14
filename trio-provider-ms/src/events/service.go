package events

import (
	"encoding/json"
	"errors"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/restclient"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type OnMessageReceive func(string) error

type EventSubscriberService interface {
	OnReceiveSyncRequest(value string) error
}

type EventSubscriberServiceImpl struct {
	syncRequestService domain.SyncRequestService
	accountRepository  domain.AccountRepository
	trioClient         domain.TrioClient
}

func NewEventSubscriberServiceImpl(syncRequestService domain.SyncRequestService,
	accountRepository domain.AccountRepository,
	trioClient domain.TrioClient) EventSubscriberService {
	return &EventSubscriberServiceImpl{
		syncRequestService: syncRequestService,
		accountRepository:  accountRepository,
		trioClient:         trioClient,
	}
}

func (subscriberService *EventSubscriberServiceImpl) OnReceiveSyncRequest(value string) error {
	providerSyncRequest, err := ParseJson(value)
	if err != nil {
		return err
	}
	SyncRequest, err := buildSyncRequest(providerSyncRequest)
	if err != nil {
		return err
	}
	SyncRequest.RequestStatus = domain.RequestStatusPending
	SyncRequest, err = subscriberService.persistSyncRequest(SyncRequest)
	switch SyncRequest.SyncType {
	case domain.SyncTypeBalances:
		err = subscriberService.synchronizeBalances(SyncRequest)
	case domain.SyncTypeTransactions:
		err = subscriberService.synchronizeTransactions(SyncRequest)
	}
	if err != nil {
		errorMessage := err.Error()
		subscriberService.syncRequestService.ChangeToFailingStatus(SyncRequest.AccountId, SyncRequest.SyncType, &errorMessage)
	}
	return err
}

func (subscriberService *EventSubscriberServiceImpl) persistSyncRequest(syncRequest *domain.SyncRequest) (*domain.SyncRequest, error) {
	return subscriberService.syncRequestService.Insert(syncRequest)
}

func (subscriberService *EventSubscriberServiceImpl) synchronizeBalances(Request *domain.SyncRequest) error {
	return subscriberService.synchronizeWithTrio(Request, subscriberService.trioClient.FetchBalancesFromBank)
}

func (subscriberService *EventSubscriberServiceImpl) synchronizeTransactions(Request *domain.SyncRequest) error {
	return subscriberService.synchronizeWithTrio(Request, subscriberService.trioClient.FetchTransactionsFromBank)
}

func (subscriberService *EventSubscriberServiceImpl) synchronizeWithTrio(Request *domain.SyncRequest, FetchFunction domain.FetchData) error {
	account, err := subscriberService.accountRepository.FindByInternalAccountId(Request.AccountId)
	if err != nil {
		return errors.New("invalid account id")
	}
	response, err := FetchFunction(*account)
	if err != nil {
		errorMessage := err.Error()
		subscriberService.syncRequestService.UpdateSyncRequestStatus(Request.AccountId, Request.SyncType,
			domain.RequestStatusFailed, &errorMessage)
		return err
	}
	if response.Data.Status == string(restclient.FailedFetchRequestStatus) {
		subscriberService.syncRequestService.UpdateSyncRequestStatus(Request.AccountId, Request.SyncType,
			domain.RequestStatusFailed, nil)
		return errors.New("fetch operation failed")
	}
	subscriberService.syncRequestService.UpdateSyncRequestStatus(Request.AccountId, Request.SyncType,
		domain.RequestStatusPending, nil)
	return nil
}

func buildSyncRequest(providerSyncRequest domain.ProviderSyncRequest) (*domain.SyncRequest, error) {
	AccountId, err := domain.ParseUUID(providerSyncRequest.AccountId)
	SyncType, err := domain.ParseSyncType(providerSyncRequest.SyncType)
	if err != nil {
		log.Println("Unable to parse ProviderSyncRequest")
		return nil, err
	}
	return &domain.SyncRequest{
		ID:            bson.NewObjectId(),
		OriginalId:    providerSyncRequest.ID,
		RequestStatus: domain.RequestStatusCreated,
		CreatedAt:     time.Now(),
		SyncType:      SyncType,
		AccountId:     AccountId.String(),
	}, nil
}

func ParseJson(value string) (domain.ProviderSyncRequest, error) {
	providerSyncRequest := domain.ProviderSyncRequest{}
	err := json.Unmarshal([]byte(value), &providerSyncRequest)
	return providerSyncRequest, err
}
