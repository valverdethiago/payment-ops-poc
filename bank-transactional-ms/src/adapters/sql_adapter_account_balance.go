package adapters

import (
	"context"
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/domain"
	"github.com/google/uuid"
)

type AccountBalanceRepositoryImpl struct {
	queries db.Querier
	ctx     context.Context
}

func NewAccountBalanceRepositoryImpl(queries db.Querier, ctx context.Context) domain.AccountBalanceRepository {
	return &AccountBalanceRepositoryImpl{
		queries: queries,
		ctx:     ctx,
	}
}

func (repository AccountBalanceRepositoryImpl) FindAllBalancesByAccount(accountId uuid.UUID) (*[]db.AccountBalance, error) {
	balances, err := repository.queries.FindAllBalancesByAccount(repository.ctx, accountId)
	return &balances, err
}

func (repository AccountBalanceRepositoryImpl) FindCurrentBalanceByAccount(accountId uuid.UUID) (*db.AccountBalance, error) {
	balance, err := repository.queries.FindCurrentBalanceByAccount(repository.ctx, accountId)
	if err != nil {
		return nil, err
	}
	return &balance, err
}

func (repository AccountBalanceRepositoryImpl) Find(id uuid.UUID) (*db.Account, *db.Bank, *db.Configuration, error) {
	row, err := repository.queries.GetFullAccountInfoByID(repository.ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	return &db.Account{
			AccountUuid:   row.AccountUuid,
			AccountNumber: row.AccountNumber,
			AccountType:   row.AccountType,
			BankUuid:      row.BankUuid,
		}, &db.Bank{
			BankUuid:    row.BankUuid,
			Code:        row.Code,
			Name:        row.Name,
			CountryCode: row.CountryCode,
		}, &db.Configuration{
			ConfigurationUuid:   row.ConfigurationUuid,
			BankUuid:            row.BankUuid,
			KafkaInputTopicName: row.KafkaInputTopicName,
		}, nil
}

func (repository AccountBalanceRepositoryImpl) FindAccountStatuses(id uuid.UUID) (db.FindAccountStatusesRow, error) {
	return repository.queries.FindAccountStatuses(repository.ctx, id)
}

func (repository AccountBalanceRepositoryImpl) UpdateAccountBalance(accountId uuid.UUID, amount float64, currency string) (*db.AccountBalance, error) {
	params := db.UpdateAccountBalanceParams{
		AccountUuid: accountId,
		Amount:      amount,
		Currency:    currency,
	}
	balance, err := repository.queries.UpdateAccountBalance(repository.ctx, params)
	return &balance, err
}
