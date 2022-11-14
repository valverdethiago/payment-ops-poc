package domain

import (
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type AccountService interface {
	FindAccountInformation(id uuid.UUID) (*db.Account, *db.Bank, *db.Configuration, error)
	IsAccountInValidState(account *db.Account) bool
	ListAll() ([]db.Account, error)
	GetAccountSnapshot(ID uuid.UUID) (*AccountResponse, error)
}

type BalanceService interface {
	FindAllBalancesByAccount(accountId uuid.UUID) (*[]db.AccountBalance, error)
	FindCurrentBalanceByAccount(accountId uuid.UUID) (*db.AccountBalance, error)
	UpdateAccountBalance(accountId uuid.UUID, amount float64, currency string) (*db.AccountBalance, error)
}

type TransactionService interface {
	FindAllTransactionsByAccount(accountId uuid.UUID) (*[]db.Transaction, error)
	InsertTransaction(transaction db.Transaction) (*db.Transaction, error)
	InsertTransactions(transactions []db.Transaction) ([]db.Transaction, error)
	FindByAccountIdAndTransactionId(accountId uuid.UUID, transactionId uuid.UUID) (*db.Transaction, error)
}

type SyncRequestService interface {
	UpdateSyncRequestStatus(id bson.ObjectId, requestStatus RequestStatus, Message *string)
	ChangeToFailingStatus(ID bson.ObjectId, Message string)
	ChangeToPendingStatus(ID bson.ObjectId)
	ChangeToSuccessfulStatus(ID bson.ObjectId)
	RequestProviderSync(name string, request ProviderSyncRequest) error
}
