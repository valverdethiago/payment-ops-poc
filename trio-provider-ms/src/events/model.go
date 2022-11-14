package events

import "github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"

type BalanceUpdateEvent struct {
	AccountID string  `json:"AccountID"`
	Amount    float64 `json:"Balance"`
	Currency  string  `json:"Currency"`
}

type TransactionsUpdateEvent struct {
	AccountId    string               `json:"accountId"`
	Transactions []domain.Transaction `json:"transactions"`
}
