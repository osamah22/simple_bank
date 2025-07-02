-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1
FOR no key UPDATE;

-- name: AccountsList :many
SELECT * FROM ACCOUNTS
LIMIT $1
OFFSET $2;


-- name: UpdateAccount :one
UPDATE accounts
SET 
    balance = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: AddToAccountBalance :one
UPDATE accounts
SET 
    balance = balance + @add_by,
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: DeleteAccount :one
DELETE FROM accounts
WHERE id = $1
RETURNING *;

