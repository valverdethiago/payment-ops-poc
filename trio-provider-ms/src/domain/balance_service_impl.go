package domain

import (
	"time"
)

type BalanceServiceImpl struct {
	eventDispatcher       EventDispatcher
	accountRepository     AccountRepository
	syncRequestRepository SyncRequestRepository
}

func NewBalanceServiceImpl(eventDispatcher EventDispatcher,
	accountRepository AccountRepository, syncRequestRepository SyncRequestRepository) BalanceService {
	return &BalanceServiceImpl{
		eventDispatcher:       eventDispatcher,
		accountRepository:     accountRepository,
		syncRequestRepository: syncRequestRepository,
	}
}

func (service BalanceServiceImpl) UpdateAccountBalance(accountId string, balance float64, currency string) error {
	account, err := service.accountRepository.FindByProviderAccountId(accountId)
	if err != nil {
		return err
	}
	err = service.eventDispatcher.TriggerBalanceUpdateEvent(account.InternalAccountId, balance, currency)
	if err != nil {
		return err
	}
	return service.updateBalanceOnDatabase(*account, balance, currency)
}

func (service BalanceServiceImpl) updateBalanceOnDatabase(account Account, balance float64, currency string) error {
	now := time.Now()
	account.LastBalanceUpdateAt = &now
	account.CurrentBalance = Balance{
		Amount:   balance,
		Currency: currency,
	}
	_, err := service.accountRepository.Update(&account)
	return err
}
