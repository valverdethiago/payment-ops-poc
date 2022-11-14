package domain

type SyncRequestRepository interface {
	Find(id string) (*SyncRequest, error)
	FindPendingRequests(AccountId string, Type SyncType) ([]SyncRequest, error)
	Store(Request *SyncRequest) (*SyncRequest, error)
}
