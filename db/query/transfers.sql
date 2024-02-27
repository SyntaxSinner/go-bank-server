-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount,
  currency
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY created_at;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;

-- name: UpdateTransfer :one
UPDATE transfers
SET amount = $2
WHERE id = $1
RETURNING *;
