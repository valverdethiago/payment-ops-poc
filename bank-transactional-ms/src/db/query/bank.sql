-- name: GetBankByID :one
SELECT * 
  FROM bank
 WHERE bank_uuid =$1;