package adapters

import (
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/domain"
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

func (store *MongoDbStore) Find(ID string) (*domain.SyncRequest, error) {
	var syncRequest domain.SyncRequest
	err := store.collection.FindId(bson.ObjectIdHex(ID)).One(&syncRequest)
	return &syncRequest, err
}

func (store *MongoDbStore) FindPendingRequests(AccountId string, SyncType string) ([]domain.SyncRequest, error) {
	var syncRequest []domain.SyncRequest
	filter := bson.D{
		{"accountid", AccountId},
		{"$or", []interface{}{
			bson.D{{"requeststatus", "PENDING"}},
			bson.D{{"requeststatus", "CREATED"}},
		}},
	}
	err := store.collection.Find(filter).All(&syncRequest)
	return syncRequest, err
}

func (store *MongoDbStore) Store(Request *domain.SyncRequest) (*domain.SyncRequest, error) {
	Request.ID = bson.NewObjectId()
	err := store.collection.Insert(&Request)
	return Request, err
}
