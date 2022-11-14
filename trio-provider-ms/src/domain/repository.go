package domain

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type SyncRequestRepository interface {
	Find(id bson.ObjectId) (*SyncRequest, error)
	FindPendingRequests(AccountId uuid.UUID, Type *SyncType) ([]SyncRequest, error)
	Insert(Request *SyncRequest) (*SyncRequest, error)
	Update(Request *SyncRequest) (*SyncRequest, error)
}

type AccountRepository interface {
	Find(id bson.ObjectId) (*AccountMapping, error)
	FindByAccountId(AccountId string) (*AccountMapping, error)
	Store(AccountMapping *AccountMapping) (*AccountMapping, error)
}
