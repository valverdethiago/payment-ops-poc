package domain

import (
	"encoding/json"
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/infra"
	"github.com/google/uuid"
)

var (
	ErrorInvalidValueForSyncType = errors.New("Invalid value for SyncType")
)

type syncRequestService struct {
	syncRequestRepository SyncRequestRepository
	producer              *infra.Producer
}

func NewSyncRequestService(syncRequestRepository SyncRequestRepository,
	producer *infra.Producer) SyncRequestService {
	return &syncRequestService{
		syncRequestRepository,
		producer,
	}
}

func (service *syncRequestService) Find(id bson.ObjectId) (*SyncRequest, error) {
	return service.syncRequestRepository.Find(id)
}

func (service *syncRequestService) Request(AccountId uuid.UUID, Type SyncType) (*SyncRequest, error) {
	requests, err := service.syncRequestRepository.FindPendingRequests(AccountId, Type)
	if err != nil || (requests != nil && len(requests) > 0) {
		if err == mgo.ErrNotFound {
			return service.createSyncRequest(AccountId, Type)
		}
		return &requests[0], err
	}
	return service.createSyncRequest(AccountId, Type)
}

func (service *syncRequestService) createSyncRequest(AccountId uuid.UUID, Type SyncType) (_ *SyncRequest, err error) {

	request := &SyncRequest{
		ID:            bson.NewObjectId(),
		RequestStatus: REQUEST_STATUS_CREATED,
		CreatedAt:     time.Now().Unix(),
		SyncType:      Type,
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
	err = service.producer.SendMessage(bson.NewObjectId().Hex(), string(jsonString))
	if err != nil {
		return nil, err
	}
	return entity, err

}
