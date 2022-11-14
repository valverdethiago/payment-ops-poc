package domain

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type SyncType string
type RequestStatus string

const (
	SYNC_TYPE_BALANCES        SyncType      = "BALANCES"
	SYNC_TYPE_TRANSACTIONS    SyncType      = "TRANSACTIONS"
	REQUEST_STATUS_CREATED    RequestStatus = "CREATED"
	REQUEST_STATUS_PENDING    RequestStatus = "PENDING"
	REQUEST_STATUS_FAILED     RequestStatus = "FAILED"
	REQUEST_STATUS_SUCCESSFUL RequestStatus = "SUCCESSFUL"
)

var (
	SYNC_TYPES map[SyncType]struct{} = map[SyncType]struct{}{
		SYNC_TYPE_BALANCES:     {},
		SYNC_TYPE_TRANSACTIONS: {},
	}
	REQUEST_STATUSES map[RequestStatus]struct{} = map[RequestStatus]struct{}{
		REQUEST_STATUS_CREATED:    {},
		REQUEST_STATUS_PENDING:    {},
		REQUEST_STATUS_FAILED:     {},
		REQUEST_STATUS_SUCCESSFUL: {},
	}
)

type ProviderSyncRequest struct {
	AccountId uuid.UUID `json:"AccountId"`
	SyncType  SyncType  `json:"SyncType"`
}

type SyncRequest struct {
	ID            bson.ObjectId `bson:"_id" json:"id,omitempty" `
	AccountId     string        `json:"AccountId"`
	SyncType      *SyncType     `json:"SyncType"`
	RequestStatus RequestStatus `json:"RequestStatus"`
	CreatedAt     int64         `json:"CreatedAt"`
}

type SyncRequestResult struct {
	ID            bson.ObjectId `json:"id,omitempty" `
	RequestStatus RequestStatus `json:"RequestStatus"`
	Message       string        `json:"Message"`
}
