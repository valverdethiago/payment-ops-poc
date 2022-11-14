package adapters

import (
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	accountsCollectionName = "accounts"
)

type AccountMongoDbRepositoryImpl struct {
	database   *mgo.Database
	collection *mgo.Collection
}

func NewAccountMongoDbRepositoryImpl(database *mgo.Database) domain.AccountRepository {
	repository := &AccountMongoDbRepositoryImpl{
		database: database,
	}
	repository.connect(accountsCollectionName)
	return repository
}

func (repository *AccountMongoDbRepositoryImpl) connect(collectionName string) {
	repository.collection = repository.database.C(collectionName)
}

func (repository *AccountMongoDbRepositoryImpl) Find(ID bson.ObjectId) (*domain.Account, error) {
	var account domain.Account
	err := repository.collection.FindId(ID).One(&account)
	return &account, err
}

func (repository *AccountMongoDbRepositoryImpl) Insert(account *domain.Account) (*domain.Account, error) {
	account.ID = bson.NewObjectId()
	err := repository.collection.Insert(&account)
	return account, err
}

func (repository *AccountMongoDbRepositoryImpl) Update(account *domain.Account) (*domain.Account, error) {
	err := repository.collection.UpdateId(account.ID, &account)
	return account, err
}

func (repository *AccountMongoDbRepositoryImpl) FindByInternalAccountId(AccountId string) (*domain.Account, error) {
	var account *domain.Account
	filter := bson.D{
		{"internal_account_id", AccountId},
	}
	err := repository.collection.Find(filter).One(&account)
	return account, err
}
func (repository *AccountMongoDbRepositoryImpl) FindByProviderAccountId(AccountId string) (*domain.Account, error) {
	var account *domain.Account
	filter := bson.D{
		{"provider_account_id", AccountId},
	}
	err := repository.collection.Find(filter).One(&account)
	return account, err
}
