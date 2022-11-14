package adapters

import (
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	transactionsCollectionName = "transactions"
)

type TransactionMongoDbRepositoryImpl struct {
	database   *mgo.Database
	collection *mgo.Collection
}

func NewTransactionMongoDbRepositoryImpl(database *mgo.Database) domain.TransactionRepository {
	repository := &TransactionMongoDbRepositoryImpl{
		database: database,
	}
	repository.connect(transactionsCollectionName)
	return repository
}

func (repository *TransactionMongoDbRepositoryImpl) connect(collectionName string) {
	repository.collection = repository.database.C(collectionName)
}

func (repository *TransactionMongoDbRepositoryImpl) Find(ID bson.ObjectId) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := repository.collection.FindId(ID).One(&transaction)
	return &transaction, err
}

func (repository *TransactionMongoDbRepositoryImpl) Insert(account *domain.Transaction) (*domain.Transaction, error) {
	account.ID = bson.NewObjectId()
	err := repository.collection.Insert(&account)
	return account, err
}

func (repository *TransactionMongoDbRepositoryImpl) FindByProviderIdAndAccountId(providerId string,
	accountId bson.ObjectId) (*domain.Transaction, error) {
	var transaction *domain.Transaction
	filter := bson.D{
		{"providerid", providerId},
		{"account_id", accountId},
	}
	err := repository.collection.Find(filter).One(&transaction)
	return transaction, err
}

func (repository *TransactionMongoDbRepositoryImpl) FindByIdentification(providerId string) (*domain.Transaction, error) {
	var transaction *domain.Transaction
	filter := bson.D{
		{"identification", providerId},
	}
	err := repository.collection.Find(filter).One(&transaction)
	return transaction, err
}
