package domain

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type SyncRequestService interface {
	Find(ID bson.ObjectId) (*SyncRequest, error)
	Request(AccountId uuid.UUID, Type SyncType) (*SyncRequest, error)
}
