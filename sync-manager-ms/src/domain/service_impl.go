package domain

import (
	"encoding/json"
	"errors"
	"time"

	errs "github.com/pkg/errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/async"
)

var (
	ErrorInvalidValueForSyncType = errors.New("Invalid value for SyncType")
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
	requests, err := service.syncRequestRepository.FindPendingRequests(AccountId, SyncType)
	if err != nil || (requests != nil && len(requests) > 0) {
		if err == mgo.ErrNotFound {
			return service.createSyncRequest(AccountId, SyncType)
		}
		return &requests[0], err
	}
	return service.createSyncRequest(AccountId, SyncType)
}

func (service *syncRequestService) createSyncRequest(AccountId string, SynchronizationType string) (*SyncRequest, error) {
	sync_type, err := ScanSyncType(SynchronizationType)
	if err != nil {
		return nil, errs.Wrap(ErrorInvalidValueForSyncType, "service.Redirect.Store")
	}

	request := &SyncRequest{
		ID:            bson.NewObjectId(),
		RequestStatus: REQUEST_STATUS_CREATED,
		CreatedAt:     time.Now().Unix(),
		SyncType:      *sync_type,
		AccountId:     AccountId,
	}
	entity, err := service.syncRequestRepository.Store(request)
	if err != nil {
		return nil, err
	}
	jsonString, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	err = service.syncRequestPublisherService.Send(jsonString)
	if err != nil {
		return nil, err
	}
	return entity, err

}
