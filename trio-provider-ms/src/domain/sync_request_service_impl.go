package domain

type SyncRequestServiceImpl struct {
	eventDispatcher       EventDispatcher
	syncRequestRepository SyncRequestRepository
}

func NewSyncRequestServiceImpl(eventsDispatcher EventDispatcher,
	syncRequestRepository SyncRequestRepository) SyncRequestService {
	return &SyncRequestServiceImpl{
		eventDispatcher:       eventsDispatcher,
		syncRequestRepository: syncRequestRepository,
	}
}

func (syncRequestService *SyncRequestServiceImpl) Insert(Request *SyncRequest) (*SyncRequest, error) {
	return syncRequestService.syncRequestRepository.Insert(Request)
}

func (syncRequestService *SyncRequestServiceImpl) FindLastRequestByAccountIdAndSyncType(internalAccountId string,
	syncType SyncType) (*SyncRequest, error) {
	syncRequest, err := syncRequestService.syncRequestRepository.FindLastRequest(internalAccountId, syncType)
	if err != nil {
		return nil, err
	}
	return syncRequest, nil
}

func (syncRequestService *SyncRequestServiceImpl) UpdateStatusByAccountIdAndSyncType(internalAccountId string,
	syncType SyncType, status RequestStatus, Message *string) error {
	syncRequest, err := syncRequestService.FindLastRequestByAccountIdAndSyncType(internalAccountId, syncType)
	if err != nil {
		return err
	}
	syncRequest.RequestStatus = status
	syncRequest.ErrorMessage = Message
	_, err = syncRequestService.syncRequestRepository.Update(syncRequest)
	if err != nil {
		return err
	}
	return syncRequestService.eventDispatcher.UpdateSyncRequestStatus(syncRequest.OriginalId, status, Message)
}

func (syncRequestService *SyncRequestServiceImpl) UpdateSyncRequestStatus(internalAccountId string,
	syncType SyncType,
	requestStatus RequestStatus, Message *string) error {
	syncRequest, err := syncRequestService.FindLastRequestByAccountIdAndSyncType(internalAccountId, syncType)
	if err != nil {
		return err
	}
	if syncRequest == nil || syncRequest.RequestStatus == requestStatus {
		return nil
	}
	syncRequest.RequestStatus = requestStatus
	syncRequest, err = syncRequestService.syncRequestRepository.Update(syncRequest)
	if err != nil {
		return err
	}
	return syncRequestService.eventDispatcher.UpdateSyncRequestStatus(syncRequest.OriginalId, requestStatus, Message)
}

func (syncRequestService *SyncRequestServiceImpl) ChangeToFailingStatus(internalAccountId string,
	syncType SyncType, Message *string) error {
	return syncRequestService.UpdateSyncRequestStatus(internalAccountId, syncType, RequestStatusFailed, Message)
}

func (syncRequestService *SyncRequestServiceImpl) ChangeToPendingStatus(internalAccountId string, syncType SyncType) error {
	return syncRequestService.UpdateSyncRequestStatus(internalAccountId, syncType, RequestStatusPending, nil)
}

func (syncRequestService *SyncRequestServiceImpl) ChangeToSuccessfulStatus(internalAccountId string, syncType SyncType) error {
	return syncRequestService.UpdateSyncRequestStatus(internalAccountId, syncType, RequestStatusSuccessful, nil)
}
