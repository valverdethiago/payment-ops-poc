package domain

import "gopkg.in/mgo.v2/bson"

type SyncRequest struct {
	ID            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	AccountId     string        `json: "AccountId" bson: "AccountId" msgpack: "AccountId"`
	SyncType      string        `json: "SyncType" bson: "SyncType" msgpack: "SyncType"`
	RequestStatus string        `json: "RequestStatus" bson: "RequestStatus" msgpack: "RequestStatus"`
	CreatedAt     int64         `json: "CreatedAt" bson: "CreatedAt" msgpack: "CreatedAt"`
}
