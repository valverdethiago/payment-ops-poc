-- name: GetAccountByID :one
SELECT * 
  FROM account
 WHERE account_uuid =$1;

-- name: IsAccountEnabled :one
 SELECT * 
  FROM account_activity
 WHERE account_uuid =$1

