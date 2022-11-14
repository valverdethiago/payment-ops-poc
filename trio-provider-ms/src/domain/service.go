package domain

type SyncRequestService interface {
	NotifyErrorOnSyncRequest(AccountId string, ErrorMessage *string) error
}

type BalanceService interface {
	UpdateAccountBalance(AccountId string, balance float64, currency string) error
}

type TransactionService interface {
	UpdateAccountTransactions(AccountId string) error
}
