package adapters

import (
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	syncRequestCollectionName = "requests"
)

type SyncRequestMongoDbRepositoryImpl struct {
	database   *mgo.Database
	collection *mgo.Collection
}

func NewSyncRepositoryMongoDbImpl(database *mgo.Database) domain.SyncRequestRepository {
	repository := &SyncRequestMongoDbRepositoryImpl{
		database: database,
	}
	repository.connect(syncRequestCollectionName)
	return repository
}

func (repository *SyncRequestMongoDbRepositoryImpl) connect(collectionName string) {
	repository.collection = repository.database.C(collectionName)
}

func (repository *SyncRequestMongoDbRepositoryImpl) Find(ID bson.ObjectId) (*domain.SyncRequest, error) {
	var syncRequest domain.SyncRequest
	err := repository.collection.FindId(ID).One(&syncRequest)
	return &syncRequest, err
}

func (repository *SyncRequestMongoDbRepositoryImpl) FindPendingRequests(AccountId uuid.UUID, Type *domain.SyncType) ([]domain.SyncRequest, error) {
	var syncRequest []domain.SyncRequest
	filter := bson.D{
		{"accountid", AccountId},
		{"synctype", Type},
		{"$or", []interface{}{
			bson.D{{"requeststatus", "PENDING"}},
			bson.D{{"requeststatus", "CREATED"}},
		}},
	}
	err := repository.collection.Find(filter).All(&syncRequest)
	return syncRequest, err
}

func (repository *SyncRequestMongoDbRepositoryImpl) Insert(Request *domain.SyncRequest) (*domain.SyncRequest, error) {
	Request.ID = bson.NewObjectId()
	err := repository.collection.Insert(&Request)
	return Request, err
}

func (repository *SyncRequestMongoDbRepositoryImpl) Update(Request *domain.SyncRequest) (*domain.SyncRequest, error) {
	err := repository.collection.Update(Request.ID, &Request)
	return Request, err
}
