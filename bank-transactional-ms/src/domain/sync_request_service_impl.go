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

func (syncRequestService *SyncRequestServiceImpl) UpdateSyncRequestStatus(id bson.ObjectId, requestStatus RequestStatus, Message *string) {
	syncRequestService.eventDispatcher.UpdateSyncRequestStatus(id, requestStatus, Message)
}
func (syncRequestService *SyncRequestServiceImpl) ChangeToFailingStatus(ID bson.ObjectId, Message string) {
	syncRequestService.UpdateSyncRequestStatus(ID, REQUEST_STATUS_FAILED, &Message)
}

func (syncRequestService *SyncRequestServiceImpl) ChangeToPendingStatus(ID bson.ObjectId) {
	syncRequestService.UpdateSyncRequestStatus(ID, REQUEST_STATUS_PENDING, nil)
}

func (syncRequestService *SyncRequestServiceImpl) ChangeToSuccessfulStatus(ID bson.ObjectId) {
	syncRequestService.UpdateSyncRequestStatus(ID, REQUEST_STATUS_SUCCESSFUL, nil)
}

func (syncRequestService *SyncRequestServiceImpl) RequestProviderSync(topic string, request ProviderSyncRequest) error {
	return syncRequestService.eventDispatcher.RequestProviderSync(topic, request)
}
