package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type SyncRequestRepository interface {
	Find(id bson.ObjectId) (*SyncRequest, error)
	FindPendingRequests(internalAccountId string, Type SyncType) ([]SyncRequest, error)
	FindLastRequest(internalAccountId string, Type SyncType) (*SyncRequest, error)
	Insert(Request *SyncRequest) (*SyncRequest, error)
	Update(Request *SyncRequest) (*SyncRequest, error)
	FindPendingRequestByAccountIdAndSyncType(accountId string, syncType SyncType) (*SyncRequest, error)
}

type AccountRepository interface {
	Find(id bson.ObjectId) (*Account, error)
	FindByInternalAccountId(internalAccountId string) (*Account, error)
	FindByProviderAccountId(providerAccountId string) (*Account, error)
	Insert(account *Account) (*Account, error)
	Update(account *Account) (*Account, error)
}

type TransactionRepository interface {
	Find(id bson.ObjectId) (*Transaction, error)
	FindByProviderIdAndAccountId(providerId string, accountId bson.ObjectId) (*Transaction, error)
	FindByIdentification(providerId string) (*Transaction, error)
	Insert(transaction *Transaction) (*Transaction, error)
}
