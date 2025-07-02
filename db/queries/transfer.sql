-- TABLE transfers (
--   id uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
--   from_account_id uuid NOT NULL,
--   to_account_id uuid NOT NULL,
--   amount BIGINT NOT NULL,
--   created_at TIMESTAMP NOT NULL DEFAULT (now())
-- );

-- name: Transfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount)
values($1, $2, $3)
returning *;

-- name: GetTransfer :one 
select * from transfers
where id = $1
limit 1;