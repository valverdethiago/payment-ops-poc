package domain

import "github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/restclient"

type FetchData func(string) (*restclient.FetchRequestResponse, error)

type TrioClient interface {
	FetchBalancesFromBank(AccountId string) (*restclient.FetchRequestResponse, error)
	ListBalance(AccountId string) (*restclient.ListBalanceResponse, error)
	FetchTransactionsFromBank(AccountId string) (*restclient.FetchRequestResponse, error)
	ListTransactions(AccountId string) (*restclient.ListTransactionsResponse, error)
}
