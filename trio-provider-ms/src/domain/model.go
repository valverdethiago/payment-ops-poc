package domain

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type RequestStatus string
type SyncType string

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
	ID            string `json:"id,omitempty" `
	AccountId     string `json:"AccountId"`
	SyncType      string `json:"SyncType"`
	RequestStatus string `json:"RequestStatus"`
	CreatedAt     int64  `json:"CreatedAt"`
}

type SyncRequest struct {
	ID            bson.ObjectId `bson:"_id" json:"id,omitempty" `
	AccountId     string        `json:"AccountId"`
	SyncType      SyncType      `json:"SyncType"`
	RequestStatus RequestStatus `json:"RequestStatus"`
	ErrorMessage  *string       `json:"ErrorMessage"`
	CreatedAt     int64         `json:"CreatedAt"`
}

type SyncRequestResult struct {
	ID            bson.ObjectId `json:"id,omitempty" `
	RequestStatus RequestStatus `json:"RequestStatus"`
	Message       *string       `json:"Message"`
	SentAt        int64         `json:"CreatedAt"`
}

type Balance struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Transaction struct {
	ID              bson.ObjectId `bson:"_id" json:"id,omitempty" `
	AccountID       bson.ObjectId `bson:"account_id" json:"id,omitempty" `
	Description     string        `json:"description"`
	DescriptionType string        `json:"description_type"`
	ProviderId      string        `json:"id"`
	Identification  string        `json:"identification"`
	InsertedAt      string        `json:"inserted_at"`
	Status          string        `json:"status"`
	Timestamp       time.Time     `json:"timestamp"`
	UpdatedAt       string        `json:"updated_at"`
	Balance         Balance       `json:"balance"`
}

type Account struct {
	ID                       bson.ObjectId `bson:"_id" json:"id,omitempty" `
	InternalAccountId        string        `json:"InternalAccountId"`
	ProviderAccountId        string        `json:"ProviderAccountId"`
	LastBalanceUpdateAt      *time.Time    `json:"LastBalanceUpdateAt"`
	LastTransactionsUpdateAt *time.Time    `json:"LastBalanceUpdateAt"`
	CurrentBalance           Balance       `json:"CurrentBalance"`
}
