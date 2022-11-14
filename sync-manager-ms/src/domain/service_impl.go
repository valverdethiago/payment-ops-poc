package domain

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/async"
)

type syncRequestService struct {
	syncRequestRepository       SyncRequestRepository
	syncRequestPublisherService async.SyncRequestPublisherService
}

func NewSyncRequestService(syncRequestRepository SyncRequestRepository,
	asyncRequestPublisherService async.SyncRequestPublisherService) SyncRequestService {
	return &syncRequestService{
		syncRequestRepository,
		asyncRequestPublisherService,
	}
}

func (service *syncRequestService) Find(id string) (*SyncRequest, error) {
	return service.syncRequestRepository.Find(id)
}

func (service *syncRequestService) Request(AccountId string, SyncType string) (*SyncRequest, error) {
	if request, err := service.syncRequestRepository.FindPendingRequest(AccountId, SyncType); err != nil && request != nil {
		if err == mgo.ErrNotFound {
			return service.createSyncRequest(AccountId, SyncType)
		}
		return request, err
	}
	return service.createSyncRequest(AccountId, SyncType)
}

func (service *syncRequestService) createSyncRequest(AccountId string, SyncType string) (*SyncRequest, error) {

	request := &SyncRequest{
		ID:            bson.NewObjectId(),
		RequestStatus: "CREATED",
		CreatedAt:     time.Now().Unix(),
		SyncType:      SyncType,
		AccountId:     AccountId,
	}
	request, err := service.syncRequestRepository.Store(request)
	if err != nil {
		return nil, err
	}
	jsonString, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	err = service.syncRequestPublisherService.Send(jsonString)
	if err != nil {
		return nil, err
	}
	return request, err

}
