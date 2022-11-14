package domain

import (
	"fmt"

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
	ID            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	AccountId     string        `json: "AccountId" bson: "AccountId"`
	SyncType      SyncType      `json: "SyncType" bson: "SyncType"`
	RequestStatus RequestStatus `json: "RequestStatus" bson: "RequestStatus"`
	CreatedAt     int64         `json: "CreatedAt" bson: "CreatedAt"`
}

func ScanSyncType(src interface{}) (*SyncType, error) {
	var e *SyncType
	switch s := src.(type) {
	case []byte:
		*e = SyncType(s)
	case string:
		*e = SyncType(s)
	default:
		return nil, fmt.Errorf("Unsupported type for SyncType %T", src)
	}
	return e, nil
}

func (e *RequestStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RequestStatus(s)
	case string:
		*e = RequestStatus(s)
	default:
		return fmt.Errorf("Unsupported type for RequestStatus %T", src)
	}
	return nil
}
