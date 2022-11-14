package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type SyncRequestServiceImpl struct {
	eventDispatcher EventDispatcher
}

func NewSyncRequestServiceImpl(eventsDispatcher EventDispatcher) SyncRequestService {
	return &SyncRequestServiceImpl{
		eventDispatcher: eventsDispatcher,
	}
}

func (syncRequestService *SyncRequestServiceImpl) UpdateSyncRequestStatus(id bson.ObjectId, requestStatus RequestStatus, Message *string) error {
	return syncRequestService.eventDispatcher.UpdateSyncRequestStatus(id, requestStatus, Message)
}

func (syncRequestService *SyncRequestServiceImpl) ChangeToFailingStatus(ID bson.ObjectId, Message *string) error {
	return syncRequestService.UpdateSyncRequestStatus(ID, RequestStatusFailed, Message)
}

func (syncRequestService *SyncRequestServiceImpl) ChangeToPendingStatus(ID bson.ObjectId) error {
	return syncRequestService.UpdateSyncRequestStatus(ID, RequestStatusPending, nil)
}

func (syncRequestService *SyncRequestServiceImpl) ChangeToSuccessfulStatus(ID bson.ObjectId) error {
	return syncRequestService.UpdateSyncRequestStatus(ID, RequestStatusSuccessful, nil)
}
