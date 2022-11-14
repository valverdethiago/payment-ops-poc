package domain

import (
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/db/sqlc"
	"github.com/google/uuid"
)

type AccountServiceImpl struct {
	accountRepository AccountRepository
}

func NewAccountServiceImpl(accountRepository AccountRepository) AccountService {
	return &AccountServiceImpl{
		accountRepository: accountRepository,
	}
}

func (accountService *AccountServiceImpl) FindAccountInformation(id uuid.UUID) (*db.Account, *db.Bank, *db.Configuration, error) {
	return accountService.accountRepository.Find(id)
}

func (accountService *AccountServiceImpl) IsAccountInValidState(account *db.Account) bool {
	statuses, err := accountService.accountRepository.FindAccountStatuses(account.AccountUuid)
	if err != nil {
		return false
	}
	return statuses.IsEnabled && !statuses.IsDisabled && !statuses.IsInvalidated
}
