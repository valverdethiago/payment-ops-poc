// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AccountActivityType string

const (
	AccountActivityTypeCREATED     AccountActivityType = "CREATED"
	AccountActivityTypeUPDATED     AccountActivityType = "UPDATED"
	AccountActivityTypeDISABLED    AccountActivityType = "DISABLED"
	AccountActivityTypeENABLED     AccountActivityType = "ENABLED"
	AccountActivityTypeINVALIDATED AccountActivityType = "INVALIDATED"
)

func (e *AccountActivityType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccountActivityType(s)
	case string:
		*e = AccountActivityType(s)
	default:
		return fmt.Errorf("unsupported scan type for AccountActivityType: %T", src)
	}
	return nil
}

type AccountType string

const (
	AccountTypeSAVINGS  AccountType = "SAVINGS"
	AccountTypeCHECKING AccountType = "CHECKING"
)

func (e *AccountType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccountType(s)
	case string:
		*e = AccountType(s)
	default:
		return fmt.Errorf("unsupported scan type for AccountType: %T", src)
	}
	return nil
}

type Account struct {
	AccountUuid   uuid.UUID   `json:"account_uuid"`
	AccountNumber string      `json:"account_number"`
	AccountType   AccountType `json:"account_type"`
	BankUuid      uuid.UUID   `json:"bank_uuid"`
}

type AccountActivity struct {
	AccountActivityUuid uuid.UUID           `json:"account_activity_uuid"`
	AccountUuid         uuid.UUID           `json:"account_uuid"`
	ActivityType        AccountActivityType `json:"activity_type"`
	DateTime            sql.NullTime        `json:"date_time"`
}

type AccountBalance struct {
	AccountBalanceUuid uuid.UUID    `json:"account_balance_uuid"`
	AccountUuid        uuid.UUID    `json:"account_uuid"`
	Amount             string       `json:"amount"`
	Currency           string       `json:"currency"`
	DateTime           sql.NullTime `json:"date_time"`
}

type Bank struct {
	BankUuid    uuid.UUID `json:"bank_uuid"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	CountryCode string    `json:"country_code"`
}

type Configuration struct {
	ConfigurationUuid   uuid.UUID `json:"configuration_uuid"`
	BankUuid            uuid.UUID `json:"bank_uuid"`
	KafkaInputTopicName string    `json:"kafka_input_topic_name"`
}

type Transaction struct {
	TransactionUuid   uuid.UUID      `json:"transaction_uuid"`
	AccountUuid       uuid.UUID      `json:"account_uuid"`
	ProviderAccountID string         `json:"provider_account_id"`
	Description       sql.NullString `json:"description"`
	DescriptionType   sql.NullString `json:"description_type"`
	Identification    string         `json:"identification"`
	Status            string         `json:"status"`
	Amount            float64        `json:"amount"`
	Currency          string         `json:"currency"`
	DateTime          time.Time      `json:"date_time"`
	CreatedAt         sql.NullTime   `json:"created_at"`
}
