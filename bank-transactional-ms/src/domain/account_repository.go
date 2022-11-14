package domain

import (
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/db/sqlc"
	"github.com/google/uuid"
)

type AccountRepository interface {
	Find(id uuid.UUID) (*db.Account, *db.Bank, *db.Configuration, error)
	FindAccountStatuses(id uuid.UUID) (db.FindAccountStatusesRow, error)
}
