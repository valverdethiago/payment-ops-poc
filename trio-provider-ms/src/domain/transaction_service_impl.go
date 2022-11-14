package domain

import (
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/restclient"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TransactionServiceImpl struct {
	eventDispatcher       EventDispatcher
	syncRequestService    SyncRequestService
	accountRepository     AccountRepository
	transactionRepository TransactionRepository
	trioClient            TrioClient
}

func NewTransactionServiceImpl(eventDispatcher EventDispatcher, syncRequestService SyncRequestService,
	accountRepository AccountRepository, transactionRepository TransactionRepository, trioClient TrioClient) TransactionService {
	return &TransactionServiceImpl{
		eventDispatcher:       eventDispatcher,
		syncRequestService:    syncRequestService,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		trioClient:            trioClient,
	}
}

func (service TransactionServiceImpl) UpdateAccountTransactions(accountId string) error {
	account, err := service.accountRepository.FindByProviderAccountId(accountId)
	if err != nil {
		return err
	}
	newTransactions, err := service.GetTransactionsList(*account)
	if err != nil {
		service.updateSyncRequest(account.InternalAccountId, RequestStatusFailed, err)
		return err
	}
	err = service.eventDispatcher.TriggerTransactionsUpdateEvent(account.InternalAccountId, newTransactions)
	if err != nil {
		service.updateSyncRequest(account.InternalAccountId, RequestStatusFailed, err)
		return err
	}
	err = service.updateTransactionsOnDatabase(*account, newTransactions)
	if err != nil {
		service.updateSyncRequest(account.InternalAccountId, RequestStatusFailed, err)
		return err
	}
	return service.updateSyncRequest(account.InternalAccountId, RequestStatusSuccessful, nil)
}

func (service TransactionServiceImpl) updateTransactionsOnDatabase(account Account, transactions []Transaction) error {
	now := time.Now()
	for _, transaction := range transactions {
		err := service.MergeTransaction(transaction)
		if err != nil {
			return err
		}
	}
	account.LastTransactionsUpdateAt = &now
	_, err := service.accountRepository.Update(&account)
	return err
}

func (service TransactionServiceImpl) GetTransactionsList(account Account) ([]Transaction, error) {
	response, err := service.trioClient.ListTransactions(account)
	if err != nil {
		return nil, err
	}
	var transactionsResult []Transaction

	for _, transaction := range response.Data {
		transactionsResult = append(transactionsResult, service.buildDomainObjectFromPayload(transaction, account.ID))
	}

	return transactionsResult, nil

}

func (service TransactionServiceImpl) buildDomainObjectFromPayload(transaction restclient.Transaction, accountID bson.ObjectId) Transaction {
	return Transaction{
		AccountID:       accountID,
		Description:     transaction.Description,
		DescriptionType: transaction.DescriptionType,
		ProviderId:      transaction.Id,
		Identification:  transaction.Identification,
		InsertedAt:      transaction.InsertedAt,
		Status:          transaction.Status,
		Timestamp:       transaction.Timestamp,
		UpdatedAt:       transaction.UpdatedAt,
		Balance: Balance{
			Amount:   transaction.Amount.Amount,
			Currency: transaction.Amount.Currency,
		},
	}

}

func (service TransactionServiceImpl) MergeTransaction(transaction Transaction) error {
	_, err := service.transactionRepository.FindByProviderIdAndAccountId(transaction.ProviderId,
		transaction.AccountID)
	if err != nil && err == mgo.ErrNotFound {
		_, err = service.transactionRepository.Insert(&transaction)
		return err
	}
	return err
}

func (service TransactionServiceImpl) updateSyncRequest(internalAccountId string, status RequestStatus, err error) error {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}
	return service.syncRequestService.UpdateStatusByAccountIdAndSyncType(internalAccountId,
		SyncTypeTransactions, status, &errorMessage)
}
