package domain

import (
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
	REQIEST_STATUS_SUCCESSFUL RequestStatus = "SUCCESSFUL"
)

type SyncRequest struct {
	ID            bson.ObjectId `json:"id,omitempty" `
	AccountId     string        `json: "AccountId"`
	SyncType      SyncType      `json: "SyncType"`
	RequestStatus RequestStatus `json: "RequestStatus"`
	CreatedAt     int64         `json: "CreatedAt"`
}
