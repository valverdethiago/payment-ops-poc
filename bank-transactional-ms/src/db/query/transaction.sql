-- name: FindTransactionsByAccount :many
  select transaction.*
    from transaction as transaction
   where transaction.account_uuid = $1
order by transaction.date_time desc;

-- name: InsertTransaction :one
INSERT INTO transaction (account_uuid, provider_account_id, description,
                         description_type, identification, status, amount,
                         currency, date_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING *;

-- name: FindByAccountIdAndTransactionId :one
select transaction.*
 from transaction as transaction
where transaction.account_uuid = $1
  and transaction.transaction_uuid = $2;

