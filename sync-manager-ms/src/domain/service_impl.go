package domain

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/google/uuid"
)

type syncRequestService struct {
	syncRequestRepository SyncRequestRepository
	eventDispatcher       EventDispatcher
}

func NewSyncRequestService(syncRequestRepository SyncRequestRepository,
	eventDispatcher EventDispatcher) SyncRequestService {
	return &syncRequestService{
		syncRequestRepository,
		eventDispatcher,
	}
}

func (service *syncRequestService) Find(id bson.ObjectId) (*SyncRequest, error) {
	return service.syncRequestRepository.Find(id)
}

func (service *syncRequestService) Request(AccountId uuid.UUID, Type *SyncType) (*SyncRequest, error) {
	requests, err := service.syncRequestRepository.FindPendingRequests(AccountId, Type)
	if err != nil || (requests != nil && len(requests) > 0) {
		if err == mgo.ErrNotFound {
			return service.createSyncRequest(AccountId, Type)
		}
		return &requests[0], err
	}
	return service.createSyncRequest(AccountId, Type)
}

func (service *syncRequestService) UpdateSyncRequestStatus(syncRequest *SyncRequest, status RequestStatus) error {
	if syncRequest.RequestStatus == status {
		return nil
	}
	syncRequest.RequestStatus = status
	syncRequest, err := service.syncRequestRepository.Update(syncRequest)
	return err
}

func (service *syncRequestService) createSyncRequest(AccountId uuid.UUID, Type *SyncType) (_ *SyncRequest, err error) {

	request := &SyncRequest{
		ID:            bson.NewObjectId(),
		RequestStatus: RequestStatusCreated,
		CreatedAt:     time.Now(),
		SyncType:      Type,
		AccountId:     AccountId.String(),
	}
	entity, err := service.syncRequestRepository.Insert(request)
	if err != nil {
		return nil, err
	}
	err = service.eventDispatcher.CreateSyncRequest(*request)
	if err != nil {
		return nil, err
	}
	return entity, err
}
