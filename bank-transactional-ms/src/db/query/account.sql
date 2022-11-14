-- name: GetAccountByID :one
SELECT * 
  FROM account
 WHERE account_uuid =$1;

-- name: IsAccountEnabled :one
 SELECT * 
  FROM account_activity
 WHERE account_uuid =$1;

 -- name: GetFullAccountInfoByID :one
SELECT * 
  FROM account
  JOIN bank as bank on bank.bank_uuid = account.bank_uuid
  JOIN configuration as config on config.bank_uuid = bank.bank_uuid
 WHERE account_uuid =$1;

 -- name: ListAllAccounts :many
SELECT *
  FROM account;

