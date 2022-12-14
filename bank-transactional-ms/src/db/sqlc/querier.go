// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	FindAccountStatuses(ctx context.Context, accountUuid uuid.UUID) (FindAccountStatusesRow, error)
	FindAllBalancesByAccount(ctx context.Context, accountUuid uuid.UUID) ([]AccountBalance, error)
	FindByAccountIdAndTransactionId(ctx context.Context, arg FindByAccountIdAndTransactionIdParams) (Transaction, error)
	FindCurrentBalanceByAccount(ctx context.Context, accountUuid uuid.UUID) (AccountBalance, error)
	FindLastActivityByAccount(ctx context.Context, accountUuid uuid.UUID) (AccountActivity, error)
	FindTransactionsByAccount(ctx context.Context, accountUuid uuid.UUID) ([]Transaction, error)
	GetAccountByID(ctx context.Context, accountUuid uuid.UUID) (Account, error)
	GetBankByID(ctx context.Context, bankUuid uuid.UUID) (Bank, error)
	GetConfigurationByBankID(ctx context.Context, bankUuid uuid.UUID) (Configuration, error)
	GetConfigurationByID(ctx context.Context, configurationUuid uuid.UUID) (Configuration, error)
	GetFullAccountInfoByID(ctx context.Context, accountUuid uuid.UUID) (GetFullAccountInfoByIDRow, error)
	InsertTransaction(ctx context.Context, arg InsertTransactionParams) (Transaction, error)
	IsAccountEnabled(ctx context.Context, accountUuid uuid.UUID) (AccountActivity, error)
	ListAllAccounts(ctx context.Context) ([]Account, error)
	UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (AccountBalance, error)
}

var _ Querier = (*Queries)(nil)
