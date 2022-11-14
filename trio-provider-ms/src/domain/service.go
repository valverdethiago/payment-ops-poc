package domain

type EventDispatcher interface {
	UpdateSyncRequestStatus(id string, requestStatus RequestStatus, Message *string) error
	TriggerBalanceUpdateEvent(accountId string, balance float64, currency string) error
	TriggerTransactionsUpdateEvent(accountId string, transactions []Transaction) error
}

type SyncRequestService interface {
	Insert(Request *SyncRequest) (*SyncRequest, error)
	FindLastRequestByAccountIdAndSyncType(internalAccountId string, syncType SyncType) (*SyncRequest, error)
	UpdateStatusByAccountIdAndSyncType(internalAccountId string, syncType SyncType, status RequestStatus, Message *string) error
	UpdateSyncRequestStatus(internalAccountId string, syncType SyncType, requestStatus RequestStatus, Message *string) error
	ChangeToFailingStatus(internalAccountId string, syncType SyncType, Message *string) error
	ChangeToPendingStatus(internalAccountId string, syncType SyncType) error
	ChangeToSuccessfulStatus(internalAccountId string, syncType SyncType) error
}

type BalanceService interface {
	UpdateAccountBalance(AccountId string, balance float64, currency string) error
}

type TransactionService interface {
	UpdateAccountTransactions(AccountId string) error
}
