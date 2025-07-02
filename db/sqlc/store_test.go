package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/osamah22/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	// Initialize
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(10)
	fmt.Println(">>before", account1.Balance, account2.Balance)

	// tests
	// run n concurrent transfer transactions
	n := 50
	errs := make(chan error)
	results := make(chan TransferFxResult)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(t.Context(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// validate
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results
		require.NoError(t, err)
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		toEntry := result.ToEntry
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		require.NotEqual(t, fromEntry.AccountID, toEntry.AccountID)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts balance
		fmt.Println(">>tx", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2 > 0)
	}

	// check accounts final balance
	updatedFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFromAccount)
	updatedToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedToAccount)
	fmt.Println(">>after", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, account1.Balance-(int64(n)*amount), updatedFromAccount.Balance)
	require.Equal(t, account2.Balance+(int64(n)*amount), updatedToAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	// Initialize
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(10)
	fmt.Println(">>before", account1.Balance, account2.Balance)

	// tests
	// run n concurrent transfer transactions
	n := 10 // this should be an even number
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 0 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// validate
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check accounts final balance
	updatedFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFromAccount)
	updatedToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedToAccount)
	fmt.Println(">>after", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, account1.Balance, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance, updatedToAccount.Balance)
}

func TestTransferTxRandomAmounts(t *testing.T) {
	// Initialize
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>before", account1.Balance, account2.Balance)

	// tests
	// run n concurrent transfer transactions
	n := 50
	amounts := make([]int64, n)

	ops := make(chan int) // operation number (for go routines and concurrent testing)
	errs := make(chan error)
	results := make(chan TransferFxResult)
	for i := 0; i < n; i++ {
		amounts[i] = utils.RandomInt(0, 250)
		go func(op int, amount int64) {
			result, err := store.TransferTx(t.Context(), TransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amounts[i],
			})
			errs <- err
			results <- result
			ops <- op
		}(i, amounts[i])
	}

	// validate
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results
		op := <-ops
		require.NoError(t, err)
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amounts[op], transfer.Amount)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amounts[op], fromEntry.Amount)

		toEntry := result.ToEntry
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amounts[op], toEntry.Amount)

		require.NotEqual(t, fromEntry.AccountID, toEntry.AccountID)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts balance
		fmt.Println(">>tx", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2 > 0)
	}

	// check accounts final balance
	var totalTransferred int64 = 0
	// Iterate each element, add an element  to a variable
	for i := 0; i < len(amounts); i++ {
		totalTransferred += amounts[i]
	}

	updatedFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFromAccount)
	updatedToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedToAccount)
	fmt.Println(">>after", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, account1.Balance-totalTransferred, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance+totalTransferred, updatedToAccount.Balance)

}
