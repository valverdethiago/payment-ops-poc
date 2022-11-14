package domain

type SyncRequestService interface {
	Find(ID string) (*SyncRequest, error)
	Request(AccountId string, SyncType string) (*SyncRequest, error)
}
