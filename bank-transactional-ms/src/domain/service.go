package domain

import (
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/db/sqlc"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type AccountService interface {
	FindAccountInformation(id uuid.UUID) (*db.Account, *db.Bank, *db.Configuration, error)
	IsAccountInValidState(account *db.Account) bool
}

type SyncRequestService interface {
	UpdateSyncRequestStatus(id bson.ObjectId, requestStatus RequestStatus, Message *string)
	ChangeToFailingStatus(ID bson.ObjectId, Message string)
	ChangeToPendingStatus(ID bson.ObjectId)
	ChangeToSuccessfulStatus(ID bson.ObjectId)
	RequestProviderSync(name string, request ProviderSyncRequest) error
}
