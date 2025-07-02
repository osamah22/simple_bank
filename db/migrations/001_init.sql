-- +goose Up

BEGIN;
CREATE TABLE accounts (
  id uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  owner VARCHAR NOT NULL,
  balance BIGINT NOT NULL,
  currency VARCHAR NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT (now()),
  updated_at TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE entries (
  id uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  account_id uuid NOT NULL,
  amount BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE transfers (
  id uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  from_account_id uuid NOT NULL,
  to_account_id uuid NOT NULL,
  amount BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE INDEX ON accounts ("owner");

CREATE INDEX ON entries ("account_id");

CREATE INDEX ON transfers ("from_account_id");

CREATE INDEX ON transfers ("to_account_id");

CREATE INDEX ON transfers ("from_account_id", "to_account_id");

COMMENT ON COLUMN entries.amount IS 'can be negative or positive';

COMMENT ON COLUMN transfers.amount IS 'must be positive';

ALTER TABLE entries ADD FOREIGN KEY ("account_id") REFERENCES accounts ("id");

ALTER TABLE transfers ADD FOREIGN KEY ("from_account_id") REFERENCES accounts ("id");

ALTER TABLE transfers ADD FOREIGN KEY ("to_account_id") REFERENCES accounts ("id");

COMMIT;

-- +goose Down
BEGIN;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS accounts;
COMMIT;