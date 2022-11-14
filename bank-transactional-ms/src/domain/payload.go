package domain

import (
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/google/uuid"
)

type AccountStatus string

const (
	AccountStatusEnabled     AccountStatus = "ENABLED"
	AccountStatusDisabled    AccountStatus = "DISABLED"
	AccountStatusInvalidated AccountStatus = "INVALIDATED"
)

type AccountResponse struct {
	AccountUuid    uuid.UUID      `json:"account_uuid"`
	AccountNumber  string         `json:"account_number"`
	AccountType    db.AccountType `json:"account_type"`
	Status         AccountStatus  `json:"status"`
	CurrentBalance *Balance       `json:"balance,omitempty"`
	Bank           db.Bank        `json:"bank"`
}
