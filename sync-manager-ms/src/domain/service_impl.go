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

func (service *syncRequestService) Request(AccountId uuid.UUID, Type *SyncType, time time.Time) (*SyncRequest, error) {
	requests, err := service.syncRequestRepository.FindPendingRequests(AccountId, Type)
	if err != nil || (requests != nil && len(requests) > 0) {
		if err == mgo.ErrNotFound {
			return service.createSyncRequest(AccountId, Type, time)
		}
		return &requests[0], err
	}
	return service.createSyncRequest(AccountId, Type, time)
}

func (service *syncRequestService) UpdateSyncRequestStatus(request *SyncRequest, status RequestStatus, message *string) error {
	if request.RequestStatus == status {
		return nil
	}
	request.RequestStatus = status
	request.Message = message
	request, err := service.syncRequestRepository.Update(request)
	return err
}

func (service *syncRequestService) createSyncRequest(AccountId uuid.UUID, Type *SyncType, time time.Time) (_ *SyncRequest, err error) {

	request := &SyncRequest{
		ID:            bson.NewObjectId(),
		RequestStatus: RequestStatusCreated,
		CreatedAt:     time,
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
