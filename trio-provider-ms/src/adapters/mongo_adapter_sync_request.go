package adapters

import (
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
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

func (repository *SyncRequestMongoDbRepositoryImpl) FindPendingRequests(AccountId string, Type domain.SyncType) ([]domain.SyncRequest, error) {
	var syncRequest []domain.SyncRequest
	filter := bson.D{
		{"account_id", AccountId},
		{"sync_type", Type},
		{"$or", []interface{}{
			bson.D{{"request_status", domain.RequestStatusCreated}},
			bson.D{{"request_status", domain.RequestStatusPending}},
		}},
	}
	err := repository.collection.Find(filter).All(&syncRequest)
	return syncRequest, err
}

func (repository *SyncRequestMongoDbRepositoryImpl) FindLastRequest(internalAccountId string, Type domain.SyncType) (*domain.SyncRequest, error) {
	var syncRequest []domain.SyncRequest
	filter := bson.D{
		{"account_id", internalAccountId},
		{"sync_type", Type},
		{"$or", []interface{}{
			bson.D{{"request_status", domain.RequestStatusCreated}},
			bson.D{{"request_status", domain.RequestStatusPending}},
		}},
	}
	err := repository.collection.Find(filter).Sort("-created_at").All(&syncRequest)
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}
	if len(syncRequest) > 0 {
		return &syncRequest[0], nil
	}
	return nil, nil
}

func (repository *SyncRequestMongoDbRepositoryImpl) Insert(Request *domain.SyncRequest) (*domain.SyncRequest, error) {
	Request.ID = bson.NewObjectId()
	err := repository.collection.Insert(&Request)
	return Request, err
}

func (repository *SyncRequestMongoDbRepositoryImpl) Update(Request *domain.SyncRequest) (*domain.SyncRequest, error) {
	err := repository.collection.UpdateId(Request.ID, &Request)
	return Request, err
}

func (repository *SyncRequestMongoDbRepositoryImpl) FindPendingRequestByAccountIdAndSyncType(accountId string, syncType domain.SyncType) (*domain.SyncRequest, error) {
	syncRequests, err := repository.FindPendingRequests(accountId, syncType)
	if err != nil || len(syncRequests) <= 0 {
		return nil, err
	}
	return &syncRequests[0], nil
}
