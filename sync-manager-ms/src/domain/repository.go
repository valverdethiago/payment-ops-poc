package domain

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type SyncRequestRepository interface {
	Find(id bson.ObjectId) (*SyncRequest, error)
	FindPendingRequests(AccountId uuid.UUID, Type *SyncType) ([]SyncRequest, error)
	Store(Request *SyncRequest) (*SyncRequest, error)
}
