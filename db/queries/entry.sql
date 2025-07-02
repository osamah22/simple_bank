-- CREATE TABLE entries (
--   id uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
--   account_id uuid,
--   amount VARCHAR NOT NULL,
--   created_at TIMESTAMP NOT NULL DEFAULT (now())
-- );

-- name: CreateEntry :one
INSERT INTO entries (account_id, amount)
values($1, $2)
returning *;

-- name: GetEntry :one
select * from entries where id = $1;

-- name: ListEntries :many
select * from entries
limit $1
offset $2;