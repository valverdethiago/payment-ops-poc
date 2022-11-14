package domain

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

func ParseUUID(ID string) (uuid.UUID, error) {
	log.Printf("Trying to parse id %s", ID)
	var result uuid.UUID
	result, err := uuid.Parse(ID)
	if err != nil {
		return result, fmt.Errorf(`cannot parse:[%s] as valid uuid`, ID)
	}
	return result, nil
}

func ParseBson(ID string) bson.ObjectId {
	log.Printf("Trying to parse id %s", ID)
	return bson.ObjectIdHex(ID)
}

func ParseSyncType(Type string) (*SyncType, error) {
	result := SyncType(Type)
	_, ok := SyncTypes[result]
	if !ok {
		return nil, fmt.Errorf(`cannot parse:[%s] as SyncType`, Type)
	}
	return &result, nil
}

func ParseSRequestStatus(Status string) (*RequestStatus, error) {
	result := RequestStatus(Status)
	_, ok := RequestStatuses[result]
	if !ok {
		return nil, fmt.Errorf(`cannot parse:[%s] as RequestStatus`, Status)
	}
	return &result, nil
}
