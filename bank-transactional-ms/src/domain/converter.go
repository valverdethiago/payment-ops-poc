package domain

import (
	"database/sql"
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/google/uuid"
	"time"
)

func ConvertTransactionsFromRestPayload(payload TransactionsUpdateEvent) ([]db.Transaction, error) {
	var transactions []db.Transaction
	accountId, err := parseUUID(payload.AccountId)
	if err != nil {
		return nil, err
	}
	for _, transaction := range payload.Transactions {
		transactions = append(transactions, ConvertTransactionFromPayload(accountId, transaction))
	}
	return transactions, nil
}

func ConvertTransactionFromPayload(accountUuid uuid.UUID, transaction Transaction) db.Transaction {
	return db.Transaction{
		TransactionUuid:   uuid.UUID{},
		AccountUuid:       accountUuid,
		ProviderAccountID: transaction.ProviderId,
		Description: sql.NullString{
			String: transaction.Description,
			Valid:  true,
		},
		DescriptionType: sql.NullString{
			String: transaction.DescriptionType,
			Valid:  true,
		},
		Identification: transaction.Identification,
		Status:         transaction.Status,
		Amount:         transaction.Balance.Amount,
		Currency:       transaction.Balance.Currency,
		DateTime:       transaction.Timestamp,
		CreatedAt:      sql.NullTime{Time: time.Now(), Valid: true},
	}
}
