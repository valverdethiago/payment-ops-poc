package adapters

import (
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/domain"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	syncRequestCollectionName = "requests"
)

type MongoDbStore struct {
	database   *mgo.Database
	collection *mgo.Collection
}

func NewMongoDbStore(database *mgo.Database) domain.SyncRequestRepository {
	store := &MongoDbStore{
		database: database,
	}
	store.connect(syncRequestCollectionName)
	return store
}

func (store *MongoDbStore) connect(collectionName string) {
	store.collection = store.database.C(collectionName)
}

func (store *MongoDbStore) Find(ID bson.ObjectId) (*domain.SyncRequest, error) {
	var syncRequest domain.SyncRequest
	err := store.collection.FindId(ID).One(&syncRequest)
	return &syncRequest, err
}

func (store *MongoDbStore) FindPendingRequests(AccountId uuid.UUID, Type *domain.SyncType) ([]domain.SyncRequest, error) {
	var syncRequest []domain.SyncRequest
	filter := bson.M{
		"account_id": AccountId.String(),
		"sync_type":  Type,
		"$or": []interface{}{
			bson.D{{"request_status", "CREATED"}},
			bson.D{{"request_status", "PENDING"}},
		},
	}
	err := store.collection.Find(filter).All(&syncRequest)
	return syncRequest, err
}

func (store *MongoDbStore) Insert(Request *domain.SyncRequest) (*domain.SyncRequest, error) {
	Request.ID = bson.NewObjectId()
	err := store.collection.Insert(&Request)
	return Request, err
}

func (store *MongoDbStore) Update(Request *domain.SyncRequest) (*domain.SyncRequest, error) {
	err := store.collection.UpdateId(Request.ID, Request)
	return Request, err
}
