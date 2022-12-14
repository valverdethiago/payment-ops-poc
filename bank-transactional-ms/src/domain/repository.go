package domain

import (
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/google/uuid"
)

type AccountRepository interface {
	Find(id uuid.UUID) (*db.Account, *db.Bank, *db.Configuration, error)
	FindAccountStatuses(id uuid.UUID) (db.FindAccountStatusesRow, error)
	ListAll() ([]db.Account, error)
}

type AccountBalanceRepository interface {
	FindAllBalancesByAccount(accountId uuid.UUID) (*[]db.AccountBalance, error)
	FindCurrentBalanceByAccount(accountId uuid.UUID) (*db.AccountBalance, error)
	UpdateAccountBalance(accountId uuid.UUID, amount float64, currency string) (*db.AccountBalance, error)
}

type TransactionRepository interface {
	FindAllTransactionsByAccount(accountId uuid.UUID) (*[]db.Transaction, error)
	InsertTransaction(transaction db.Transaction) (*db.Transaction, error)
	FindByAccountIdAndTransactionId(accountId uuid.UUID, transactionId uuid.UUID) (*db.Transaction, error)
}
