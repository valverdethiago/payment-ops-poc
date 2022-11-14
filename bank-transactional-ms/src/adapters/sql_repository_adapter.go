package adapters

import (
	"context"

	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/domain"
	"github.com/google/uuid"
)

type AccountRepositoryImpl struct {
	queries db.Querier
	ctx     context.Context
}

func NewAccountRepositoryImpl(queries db.Querier, ctx context.Context) domain.AccountRepository {
	return &AccountRepositoryImpl{
		queries: queries,
		ctx:     ctx,
	}
}

func (repository AccountRepositoryImpl) Find(id uuid.UUID) (*db.Account, *db.Bank, *db.Configuration, error) {
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

func (repository AccountRepositoryImpl) FindAccountStatuses(id uuid.UUID) (db.FindAccountStatusesRow, error) {
	return repository.queries.FindAccountStatuses(repository.ctx, id)
}
