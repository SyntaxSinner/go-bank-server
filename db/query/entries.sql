-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY created_at;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;

-- name: UpdateEntry :one
UPDATE entries
SET balance = $2
WHERE id = $1
RETURNING *;
