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

const getFullAccountInfoByID = `-- name: GetFullAccountInfoByID :one
SELECT account_uuid, account_number, account_type, account.bank_uuid, bank.bank_uuid, code, name, country_code, configuration_uuid, config.bank_uuid, kafka_input_topic_name 
  FROM account
  JOIN bank as bank on bank.bank_uuid = account.bank_uuid
  JOIN configuration as config on config.bank_uuid = bank.bank_uuid
 WHERE account_uuid =$1
`

type GetFullAccountInfoByIDRow struct {
	AccountUuid         uuid.UUID   `json:"account_uuid"`
	AccountNumber       string      `json:"account_number"`
	AccountType         AccountType `json:"account_type"`
	BankUuid            uuid.UUID   `json:"bank_uuid"`
	BankUuid_2          uuid.UUID   `json:"bank_uuid_2"`
	Code                string      `json:"code"`
	Name                string      `json:"name"`
	CountryCode         string      `json:"country_code"`
	ConfigurationUuid   uuid.UUID   `json:"configuration_uuid"`
	BankUuid_3          uuid.UUID   `json:"bank_uuid_3"`
	KafkaInputTopicName string      `json:"kafka_input_topic_name"`
}

func (q *Queries) GetFullAccountInfoByID(ctx context.Context, accountUuid uuid.UUID) (GetFullAccountInfoByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getFullAccountInfoByID, accountUuid)
	var i GetFullAccountInfoByIDRow
	err := row.Scan(
		&i.AccountUuid,
		&i.AccountNumber,
		&i.AccountType,
		&i.BankUuid,
		&i.BankUuid_2,
		&i.Code,
		&i.Name,
		&i.CountryCode,
		&i.ConfigurationUuid,
		&i.BankUuid_3,
		&i.KafkaInputTopicName,
	)
	return i, err
}

const isAccountEnabled = `-- name: IsAccountEnabled :one
 SELECT account_activity_uuid, account_uuid, activity_type, date_time 
  FROM account_activity
 WHERE account_uuid =$1
`

func (q *Queries) IsAccountEnabled(ctx context.Context, accountUuid uuid.UUID) (AccountActivity, error) {
	row := q.db.QueryRowContext(ctx, isAccountEnabled, accountUuid)
	var i AccountActivity
	err := row.Scan(
		&i.AccountActivityUuid,
		&i.AccountUuid,
		&i.ActivityType,
		&i.DateTime,
	)
	return i, err
}

const listAllAccounts = `-- name: ListAllAccounts :many
SELECT account_uuid, account_number, account_type, bank_uuid
  FROM account
`

func (q *Queries) ListAllAccounts(ctx context.Context) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAllAccounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.AccountUuid,
			&i.AccountNumber,
			&i.AccountType,
			&i.BankUuid,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
