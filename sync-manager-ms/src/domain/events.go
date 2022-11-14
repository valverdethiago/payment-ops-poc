package domain

import (
	"encoding/json"
)

type OnMessageReceive func(string) error

type EventDispatcher interface {
	CreateSyncRequest(syncRequest SyncRequest) error
}

type EventSubscriberService interface {
	OnReceiveSyncRequestUpdate(value string) error
}

type EventSubscriberServiceImpl struct {
	syncRequestService SyncRequestService
}

func NewEventSubscriberServiceImpl(syncRequestService SyncRequestService) EventSubscriberService {
	return &EventSubscriberServiceImpl{
		syncRequestService: syncRequestService,
	}
}

func (e EventSubscriberServiceImpl) OnReceiveSyncRequestUpdate(value string) error {
	syncRequestEvent, err := ParseSyncRequestJson(value)
	if err != nil {
		return err
	}
	ID := ParseBson(syncRequestEvent.Id)
	syncRequest, err := e.syncRequestService.Find(ID)
	if err != nil {
		return err
	}
	return e.syncRequestService.UpdateSyncRequestStatus(syncRequest,
		syncRequestEvent.RequestStatus, syncRequestEvent.Message)
}

func ParseSyncRequestJson(value string) (SyncRequestEvent, error) {
	syncRequest := SyncRequestEvent{}
	err := json.Unmarshal([]byte(value), &syncRequest)
	return syncRequest, err
}
