package domain

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
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
	SyncTypes map[SyncType]struct{} = map[SyncType]struct{}{
		SyncTypeBalances:     {},
		SyncTypeTransactions: {},
	}
	RequestStatuses map[RequestStatus]struct{} = map[RequestStatus]struct{}{
		RequestStatusCreated:    {},
		RequestStatusPending:    {},
		RequestStatusFailed:     {},
		RequestStatusSuccessful: {},
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
