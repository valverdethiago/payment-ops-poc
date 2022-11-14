package domain

import (
	"time"
)

type BalanceServiceImpl struct {
	eventDispatcher    EventDispatcher
	accountRepository  AccountRepository
	syncRequestService SyncRequestService
}

func NewBalanceServiceImpl(eventDispatcher EventDispatcher,
	accountRepository AccountRepository,
	syncRequestService SyncRequestService) BalanceService {
	return &BalanceServiceImpl{
		eventDispatcher:    eventDispatcher,
		accountRepository:  accountRepository,
		syncRequestService: syncRequestService,
	}
}

func (service BalanceServiceImpl) UpdateAccountBalance(accountId string, balance float64, currency string) error {
	account, err := service.accountRepository.FindByProviderAccountId(accountId)
	if err != nil {
		return err
	}
	err = service.eventDispatcher.TriggerBalanceUpdateEvent(account.InternalAccountId, balance, currency)
	if err != nil {
		service.updateSyncRequest(account.InternalAccountId, RequestStatusFailed, err)
		return err
	}
	err = service.updateBalanceOnDatabase(*account, balance, currency)
	if err != nil {
		service.updateSyncRequest(account.InternalAccountId, RequestStatusFailed, err)
		return err
	}
	service.updateSyncRequest(account.InternalAccountId, RequestStatusSuccessful, nil)
	return nil
}

func (service BalanceServiceImpl) updateBalanceOnDatabase(account Account, balance float64, currency string) error {
	now := time.Now()
	account.LastBalanceUpdateAt = &now
	account.CurrentBalance = &Balance{
		Amount:   balance,
		Currency: currency,
	}
	_, err := service.accountRepository.Update(&account)
	return err
}

func (service BalanceServiceImpl) updateSyncRequest(internalAccountId string, status RequestStatus, err error) {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}
	service.syncRequestService.UpdateStatusByAccountIdAndSyncType(internalAccountId,
		SyncTypeBalances, status, &errorMessage)
}
