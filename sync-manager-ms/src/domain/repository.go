package domain

type SyncRequestRepository interface {
	Find(id string) (*SyncRequest, error)
	FindPendingRequest(AccountId string, SyncType string) (*SyncRequest, error)
	Store(Request *SyncRequest) (*SyncRequest, error)
}
