package domain

import (
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/google/uuid"
)

type TransactionServiceImpl struct {
	transactionRepository TransactionRepository
}

func NewTransactionServiceImpl(transactionRepository TransactionRepository) TransactionService {
	return &TransactionServiceImpl{
		transactionRepository: transactionRepository,
	}
}

func (service TransactionServiceImpl) FindAllTransactionsByAccount(accountId uuid.UUID) (*[]db.Transaction, error) {
	return service.transactionRepository.FindAllTransactionsByAccount(accountId)
}

func (service TransactionServiceImpl) InsertTransaction(transaction db.Transaction) (*db.Transaction, error) {
	return service.transactionRepository.InsertTransaction(transaction)
}

func (service TransactionServiceImpl) InsertTransactions(transactions []db.Transaction) ([]db.Transaction, error) {
	var result []db.Transaction
	for _, transaction := range transactions {
		dbTransaction, err := service.transactionRepository.InsertTransaction(transaction)
		if err != nil {
			return nil, err
		}
		result = append(result, *dbTransaction)
	}
	return result, nil
}
