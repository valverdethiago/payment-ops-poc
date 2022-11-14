// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: account.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const getAccountByID = `-- name: GetAccountByID :one
SELECT account_uuid, account_number, account_type, bank_uuid 
  FROM account
 WHERE account_uuid =$1
`

func (q *Queries) GetAccountByID(ctx context.Context, accountUuid uuid.UUID) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountByID, accountUuid)
	var i Account
	err := row.Scan(
		&i.AccountUuid,
		&i.AccountNumber,
		&i.AccountType,
		&i.BankUuid,
	)
	return i, err
}
