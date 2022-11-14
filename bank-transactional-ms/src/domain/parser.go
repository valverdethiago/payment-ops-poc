package domain

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

func parseUUID(ID string) (uuid.UUID, error) {
	log.Printf("Trying to parse id %s", ID)
	var result uuid.UUID
	result, err := uuid.Parse(ID)
	if err != nil {
		return result, fmt.Errorf(`cannot parse:[%s] as valid uuid`, ID)
	}
	return result, nil
}

func parseBson(ID string) bson.ObjectId {
	log.Printf("Trying to parse id %s", ID)
	return bson.ObjectIdHex(ID)
}

func parseSyncType(Type string) (*SyncType, error) {
	result := SyncType(Type)
	_, ok := SYNC_TYPES[result]
	if !ok {
		return nil, fmt.Errorf(`cannot parse:[%s] as SyncType`, Type)
	}
	return &result, nil
}
