package domain

import "gopkg.in/mgo.v2/bson"

type EventDispatcher interface {
	UpdateSyncRequestStatus(id bson.ObjectId, requestStatus RequestStatus, Message *string) error
	TriggerBalanceUpdateEvent(accountId string, balance float64, currency string) error
	TriggerTransactionsUpdateEvent(accountId string, transactions []Transaction) error
}

type SyncRequestService interface {
	UpdateSyncRequestStatus(ID bson.ObjectId, requestStatus RequestStatus, Message *string) error
	ChangeToFailingStatus(ID bson.ObjectId, Message *string) error
	ChangeToPendingStatus(ID bson.ObjectId) error
	ChangeToSuccessfulStatus(ID bson.ObjectId) error
}

type BalanceService interface {
	UpdateAccountBalance(AccountId string, balance float64, currency string) error
}

type TransactionService interface {
	UpdateAccountTransactions(AccountId string) error
}
