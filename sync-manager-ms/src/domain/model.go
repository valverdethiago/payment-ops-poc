package domain

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type SyncType string
type RequestStatus string

const (
	SyncTypeBalances        SyncType      = "BALANCES"
	SyncTypeTransactions    SyncType      = "TRANSACTIONS"
	RequestStatusCreated    RequestStatus = "CREATED"
	RequestStatusPending    RequestStatus = "PENDING"
	RequestStatusFailed     RequestStatus = "FAILED"
	RequestStatusSuccessful RequestStatus = "SUCCESSFUL"
)

var (
	SyncTypes = map[SyncType]struct{}{
		SyncTypeBalances:     {},
		SyncTypeTransactions: {},
	}
	RequestStatuses = map[RequestStatus]struct{}{
		RequestStatusCreated:    {},
		RequestStatusPending:    {},
		RequestStatusFailed:     {},
		RequestStatusSuccessful: {},
	}
)

type ProviderSyncRequest struct {
	AccountId uuid.UUID `bson:"account_id" json:"account_id"`
	SyncType  SyncType  `bson:"sync_type" json:"sync_type"`
}

type SyncRequest struct {
	ID            bson.ObjectId `bson:"_id" json:"id,omitempty" `
	AccountId     string        `bson:"account_id" json:"account_id"`
	SyncType      *SyncType     `bson:"sync_type" json:"sync_type"`
	RequestStatus RequestStatus `bson:"request_status" json:"request_status"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
}

type SyncRequestResult struct {
	ID            bson.ObjectId `bson:"_id" json:"id,omitempty" `
	RequestStatus RequestStatus `bson:"request_status" json:"request_status"`
	Message       string        `bson:"message" json:"message"`
}

type SyncRequestEvent struct {
	Id            string        `json:"id"`
	RequestStatus RequestStatus `json:"request_status"`
	Message       *string       `json:"message"`
	SentAt        time.Time     `json:"sent_at"`
}
