package domain

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type SyncRequestService interface {
	Find(ID bson.ObjectId) (*SyncRequest, error)
	Request(AccountId uuid.UUID, Type *SyncType, time time.Time) (*SyncRequest, error)
	UpdateSyncRequestStatus(request *SyncRequest, status RequestStatus, message *string) error
}
