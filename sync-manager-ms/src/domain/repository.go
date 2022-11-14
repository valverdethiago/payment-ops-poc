package domain

type SyncRequestRepository interface {
	Find(id string) (*SyncRequest, error)
	FindPendingRequests(AccountId string, SyncType string) ([]SyncRequest, error)
	Store(Request *SyncRequest) (*SyncRequest, error)
}
