-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id,
                       to_account_id,
                       amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTransfer :one
SELECT t.id, t.from_account_id, t.to_account_id, t.amount, t.created_at
FROM transfers t,
     accounts a
WHERE t.id = $1
  and ((a.owner = @username::text and t.from_account_id = a.id)
    or (a.owner = @username::text and t.to_account_id = a.id))
LIMIT 1;

-- name: ListTransfers :many
SELECT t.id, t.from_account_id, t.to_account_id, t.amount, t.created_at
FROM transfers t,
     accounts a
WHERE (from_account_id = $1
    OR to_account_id = $2)
  and ((a.owner = @username::text and t.from_account_id = a.id)
    or (a.owner = @username::text and t.to_account_id = a.id))
ORDER BY id
LIMIT $3 OFFSET $4;
