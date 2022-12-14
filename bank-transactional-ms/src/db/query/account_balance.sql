-- name: FindCurrentBalanceByAccount :one
  select balance.*
    from account_balance as balance
   where balance.account_uuid = $1
order by balance.date_time desc
   limit 1;

-- name: FindAllBalancesByAccount :many
  select balance.*
    from account_balance as balance
   where balance.account_uuid = $1
order by balance.date_time desc;

-- name: UpdateAccountBalance :one
INSERT INTO account_balance (account_uuid, amount, currency)
      VALUES ($1, $2, $3)
RETURNING *;
