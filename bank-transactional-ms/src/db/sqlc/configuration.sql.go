// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: configuration.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const getConfigurationByBankID = `-- name: GetConfigurationByBankID :one
SELECT configuration_uuid, bank_uuid, kafka_input_topic_name 
  FROM configuration
 WHERE bank_uuid =$1
`

func (q *Queries) GetConfigurationByBankID(ctx context.Context, bankUuid uuid.UUID) (Configuration, error) {
	row := q.db.QueryRowContext(ctx, getConfigurationByBankID, bankUuid)
	var i Configuration
	err := row.Scan(&i.ConfigurationUuid, &i.BankUuid, &i.KafkaInputTopicName)
	return i, err
}

const getConfigurationByID = `-- name: GetConfigurationByID :one
SELECT configuration_uuid, bank_uuid, kafka_input_topic_name 
  FROM configuration
 WHERE configuration_uuid =$1
`

func (q *Queries) GetConfigurationByID(ctx context.Context, configurationUuid uuid.UUID) (Configuration, error) {
	row := q.db.QueryRowContext(ctx, getConfigurationByID, configurationUuid)
	var i Configuration
	err := row.Scan(&i.ConfigurationUuid, &i.BankUuid, &i.KafkaInputTopicName)
	return i, err
}
