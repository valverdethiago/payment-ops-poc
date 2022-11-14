package domain

import (
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/restclient"
)

type FetchData func(account Account) (*restclient.FetchRequestResponse, error)

type TrioClient interface {
	FetchBalancesFromBank(account Account) (*restclient.FetchRequestResponse, error)
	ListBalance(AccountId string) (*restclient.ListBalanceResponse, error)
	FetchTransactionsFromBank(account Account) (*restclient.FetchRequestResponse, error)
	ListTransactions(account Account) (*restclient.ListTransactionsResponse, error)
}
