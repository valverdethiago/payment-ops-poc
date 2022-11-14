package domain

type SyncRequestService interface {
	Find(ID string) (*SyncRequest, error)
	Request(AccountId string, Type SyncType) (*SyncRequest, error)
}
