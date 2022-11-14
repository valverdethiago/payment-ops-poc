package domain

import (
	"database/sql"
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/google/uuid"
)

type AccountServiceImpl struct {
	accountRepository AccountRepository
	balanceService    BalanceService
}

func NewAccountServiceImpl(accountRepository AccountRepository, balanceService BalanceService) AccountService {
	return &AccountServiceImpl{
		accountRepository: accountRepository,
		balanceService:    balanceService,
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

func (accountService *AccountServiceImpl) ListAll() ([]db.Account, error) {
	return accountService.accountRepository.ListAll()
}

func (accountService *AccountServiceImpl) GetAccountSnapshot(ID uuid.UUID) (*AccountResponse, error) {
	account, bank, _, err := accountService.FindAccountInformation(ID)
	if err != nil {
		return nil, err
	}
	statuses, err := accountService.accountRepository.FindAccountStatuses(account.AccountUuid)
	if err != nil {
		return nil, err
	}

	balance, err := accountService.balanceService.FindCurrentBalanceByAccount(ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	response := &AccountResponse{
		AccountUuid:    account.AccountUuid,
		AccountNumber:  account.AccountNumber,
		AccountType:    account.AccountType,
		Status:         accountService.GetAccountStatus(statuses),
		CurrentBalance: nil,
		Bank:           *bank,
	}
	if balance != nil {
		response.CurrentBalance = &Balance{
			Amount:   balance.Amount,
			Currency: balance.Currency,
		}
	}
	return response, nil
}

func (accountService *AccountServiceImpl) GetAccountStatus(statuses db.FindAccountStatusesRow) AccountStatus {
	if statuses.IsEnabled {
		return AccountStatusEnabled
	} else if statuses.IsInvalidated {
		return AccountStatusInvalidated
	} else {
		return AccountStatusDisabled
	}

}
