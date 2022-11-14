-- name: GetConfigurationByID :one
SELECT * 
  FROM configuration
 WHERE configuration_uuid =$1;

-- name: GetConfigurationByBankID :one
SELECT * 
  FROM configuration
 WHERE bank_uuid =$1;
