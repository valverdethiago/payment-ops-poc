package domain

import (
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/google/uuid"
)

type AccountBalanceServiceImpl struct {
	accountBalanceRepository AccountBalanceRepository
}

func NewAccountBalanceServiceImpl(accountBalanceRepository AccountBalanceRepository) BalanceService {
	return &AccountBalanceServiceImpl{
		accountBalanceRepository: accountBalanceRepository,
	}
}

func (service AccountBalanceServiceImpl) FindAllBalancesByAccount(accountId uuid.UUID) (*[]db.AccountBalance, error) {
	return service.accountBalanceRepository.FindAllBalancesByAccount(accountId)
}

func (service AccountBalanceServiceImpl) FindCurrentBalanceByAccount(accountId uuid.UUID) (*db.AccountBalance, error) {
	return service.accountBalanceRepository.FindCurrentBalanceByAccount(accountId)
}

func (service AccountBalanceServiceImpl) UpdateAccountBalance(accountId uuid.UUID, amount float64,
	currency string) (*db.AccountBalance, error) {
	return service.accountBalanceRepository.UpdateAccountBalance(accountId, amount, currency)
}
