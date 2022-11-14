package adapters

import (
	"context"
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/domain"
	"github.com/google/uuid"
)

type TransactionRepositoryImpl struct {
	queries db.Querier
	ctx     context.Context
}

func NewTransactionRepositoryImpl(queries db.Querier, ctx context.Context) domain.TransactionRepository {
	return &TransactionRepositoryImpl{
		queries: queries,
		ctx:     ctx,
	}
}

func (repository TransactionRepositoryImpl) FindAllTransactionsByAccount(accountId uuid.UUID) (*[]db.Transaction, error) {
	transactions, err := repository.queries.FindTransactionsByAccount(repository.ctx, accountId)
	return &transactions, err
}

func (repository TransactionRepositoryImpl) InsertTransaction(transaction db.Transaction) (*db.Transaction, error) {
	params := &db.InsertTransactionParams{
		AccountUuid:       transaction.AccountUuid,
		ProviderAccountID: transaction.ProviderAccountID,
		Description:       transaction.Description,
		DescriptionType:   transaction.DescriptionType,
		Identification:    transaction.Identification,
		Status:            transaction.Status,
		Amount:            transaction.Amount,
		Currency:          transaction.Currency,
		DateTime:          transaction.DateTime,
	}
	transaction, err := repository.queries.InsertTransaction(repository.ctx, *params)
	return &transaction, err

}
