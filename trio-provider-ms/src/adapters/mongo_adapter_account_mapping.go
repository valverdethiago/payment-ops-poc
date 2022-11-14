package adapters

import (
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	accountsMappingCollectionName = "accounts"
)

type AccountMappingMongoDbRepositoryImpl struct {
	database   *mgo.Database
	collection *mgo.Collection
}

func NewAccountMappingMongoDbRepositoryImpl(database *mgo.Database) domain.AccountRepository {
	repository := &AccountMappingMongoDbRepositoryImpl{
		database: database,
	}
	repository.connect(accountsMappingCollectionName)
	return repository
}

func (repository *AccountMappingMongoDbRepositoryImpl) connect(collectionName string) {
	repository.collection = repository.database.C(collectionName)
}

func (repository *AccountMappingMongoDbRepositoryImpl) Find(id bson.ObjectId) (*domain.AccountMapping, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *AccountMappingMongoDbRepositoryImpl) Store(AccountMapping *domain.AccountMapping) (*domain.AccountMapping, error) {
	AccountMapping.ID = bson.NewObjectId()
	err := repository.collection.Insert(&AccountMapping)
	return AccountMapping, err
}

func (repository *AccountMappingMongoDbRepositoryImpl) FindByAccountId(AccountId string) (*domain.AccountMapping, error) {
	//var accountMapping domain.AccountMapping
	//filter := bson.D{
	//	{"internalaccountid", AccountId.String()},
	//}
	//err := repository.collection.Find(filter).One(&accountMapping)
	//return &accountMapping, err

	return &domain.AccountMapping{
		InternalAccountId: "9b2272fe-53ad-4d5d-bfaf-ce2f1cf27ccf",
		ProviderAccountId: "6501663c-3a19-47df-9a2a-bc0796b702fd",
	}, nil
}
