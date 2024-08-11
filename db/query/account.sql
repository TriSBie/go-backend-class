-- name: CreateAccounts :one
INSERT INTO accounts (owner, balance, currency) 
VALUES ($1, $2, $3) 
RETURNING *;

-- name: GetAccountById :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1;

-- name: GetAccounts :many
SELECT * FROM accounts
ORDER BY owner
LIMIT $1
OFFSET $2;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: UpdateAccountBalance :one
UPDATE accounts SET balance = $1 WHERE id = $2
RETURNING *;

-- name: AddMoneyToAccount :one
UPDATE accounts SET balance = balance + $1 WHERE id = $2 RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;