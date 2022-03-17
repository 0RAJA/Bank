-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1 for no key update;
--不会更改就爱你支队

-- name: ListAccounts :many
SELECT *
FROM accounts
ORDER BY ID
LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount) --自定义参数生成的go字段名
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE id = $1;
