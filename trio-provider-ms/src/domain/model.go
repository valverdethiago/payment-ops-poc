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
	ID            string    `bson:"_id" json:"id,omitempty" `
	AccountId     string    `bson:"account_id" json:"account_id"`
	SyncType      string    `bson:"sync_type" json:"sync_type"`
	RequestStatus string    `bson:"request_status" json:"request_status"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}

type SyncRequest struct {
	ID            bson.ObjectId `bson:"_id" json:"id,omitempty" `
	AccountId     string        `bson:"account_id" json:"account_id"`
	SyncType      SyncType      `bson:"sync_type" json:"sync_type"`
	RequestStatus RequestStatus `bson:"request_status" json:"request_status"`
	ErrorMessage  *string       `bson:"error_message" json:"error_message"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
}

type SyncRequestResult struct {
	ID            bson.ObjectId `bson:"_id" json:"id,omitempty" `
	RequestStatus RequestStatus `bson:"request_status" json:"RequestStatus"`
	Message       *string       `bson:"message" json:"message"`
	SentAt        time.Time     `bson:"sent_at" json:"sent_at"`
}

type Balance struct {
	Amount   float64 `bson:"amount" json:"amount"`
	Currency string  `bson:"currency" json:"currency"`
}

type Transaction struct {
	ID              bson.ObjectId `bson:"_id" json:"id,omitempty" `
	AccountID       bson.ObjectId `bson:"account_id" json:"account_id,omitempty" `
	Description     string        `bson:"description" json:"description"`
	DescriptionType string        `bson:"description_type" json:"description_type"`
	ProviderId      string        `bson:"provider_id" json:"provider_id"`
	Identification  string        `bson:"identification" json:"identification"`
	InsertedAt      string        `bson:"inserted_at" json:"inserted_at"`
	Status          string        `bson:"status" json:"status"`
	Timestamp       time.Time     `bson:"timestamp" json:"timestamp"`
	UpdatedAt       string        `bson:"updated_at" json:"updated_at"`
	Balance         Balance       `bson:"balance" json:"balance"`
}

type Account struct {
	ID                       bson.ObjectId `bson:"_id" json:"id,omitempty" `
	InternalAccountId        string        `bson:"internal_account_id" json:"internal_account_id"`
	ProviderAccountId        string        `bson:"provider_account_id" json:"provider_account_id"`
	LastBalanceUpdateAt      *time.Time    `bson:"last_balance_update_at" json:"last_balance_update_at"`
	LastTransactionsUpdateAt *time.Time    `bson:"last_transactions_update_at" json:"last_transactions_update_at"`
	CurrentBalance           *Balance      `bson:"current_balance" json:"current_balance"`
}
