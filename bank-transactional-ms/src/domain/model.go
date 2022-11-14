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
)

type SyncRequest struct {
	ID            string        `bson:"_id" json:"id,omitempty" `
	AccountId     string        `bson:"account_id" json:"account_id"`
	SyncType      SyncType      `bson:"sync_type" json:"sync_type"`
	RequestStatus RequestStatus `bson:"request_status" json:"request_status"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
}

type BalanceUpdateEvent struct {
	AccountID string  `json:"AccountID"`
	Balance   float64 `json:"Balance"`
	Currency  string  `json:"Currency"`
}

type Balance struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Transaction struct {
	AccountId       string    `json:"account_id"`
	Description     string    `json:"description"`
	DescriptionType string    `json:"description_type"`
	ProviderId      string    `json:"provider_id"`
	Identification  string    `json:"identification"`
	InsertedAt      string    `json:"inserted_at"`
	Status          string    `json:"status"`
	Timestamp       time.Time `json:"timestamp"`
	UpdatedAt       string    `json:"updated_at"`
	Balance         Balance   `json:"balance"`
}

type TransactionsUpdateEvent struct {
	AccountId    string        `json:"accountId"`
	Transactions []Transaction `json:"transactions"`
}

type SyncRequestResult struct {
	ID            bson.ObjectId `bson:"_id" json:"id,omitempty" `
	RequestStatus RequestStatus `bson:"request_status" json:"request_status"`
	Message       *string       `bson:"message" json:"message"`
	SentAt        time.Time     `bson:"sent_at" json:"sent_at"`
}

type ProviderSyncRequest struct {
	ID            bson.ObjectId `bson:"_id" json:"id,omitempty" `
	AccountId     uuid.UUID     `bson:"account_id" json:"account_id"`
	SyncType      SyncType      `bson:"sync_type" json:"sync_type"`
	RequestStatus RequestStatus `bson:"request_status" json:"request_status"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
}
