package db

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Store provides all funcations to excute db queries and transaction
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx excutes a funcation within a transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %s, rbErr: %s", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferFxParams struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}

type TransferFxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, args TransferParams) (TransferFxResult, error) {
	var result TransferFxResult

	err :=
		store.execTx(ctx, func(q *Queries) error {
			var err error

			txName := ctx.Value(txKey)

			// making a transfer record
			fmt.Println(txName, "create transfer")
			result.Transfer, err = q.Transfer(ctx, TransferParams{
				FromAccountID: args.FromAccountID,
				ToAccountID:   args.ToAccountID,
				Amount:        args.Amount,
			})
			if err != nil {
				return err
			}

			// making an entries records
			fmt.Println(txName, "create entry 1")
			result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
				AccountID: args.FromAccountID,
				Amount:    -args.Amount,
			})
			if err != nil {
				return err
			}

			fmt.Println(txName, "create entry 2")
			result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
				AccountID: args.ToAccountID,
				Amount:    args.Amount,
			})
			if err != nil {
				return err
			}

			// updating accounts balances

			// the following condition is to prevent a deadlock from happening by comapring the uuid
			account1Params := AddToAccountBalanceParams{
				ID:    args.FromAccountID,
				AddBy: -args.Amount,
			}
			account2Params := AddToAccountBalanceParams{
				ID:    args.ToAccountID,
				AddBy: args.Amount,
			}

			if bytes.Compare(args.FromAccountID[:], args.ToAccountID[:]) > 0 {
				result.FromAccount, err = q.AddToAccountBalance(ctx, account1Params)
				if err != nil {
					return err
				}
				result.ToAccount, err = q.AddToAccountBalance(ctx, account2Params)
				if err != nil {
					return err
				}
			} else {
				result.ToAccount, err = q.AddToAccountBalance(ctx, account2Params)
				if err != nil {
					return err
				}
				result.FromAccount, err = q.AddToAccountBalance(ctx, account1Params)
				if err != nil {
					return err
				}
			}

			return nil
		})

	return result, err
}
