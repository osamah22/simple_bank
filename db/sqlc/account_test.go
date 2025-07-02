package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/osamah22/simple_bank/utils"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandonOwner(),
		Balance:  utils.RandonAmount(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.NotEqual(t, uuid.Nil, account.ID)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	dbAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)

	require.Equal(t, account.ID, dbAccount.ID)
	require.Equal(t, account.Owner, dbAccount.Owner)
	require.Equal(t, account.Balance, dbAccount.Balance)
	require.Equal(t, account.Currency, dbAccount.Currency)
	require.WithinDuration(t, account.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: account.Balance,
	}

	dbAccount, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, account.ID, dbAccount.ID)
	require.Equal(t, account.Owner, dbAccount.Owner)
	require.Equal(t, args.Balance, dbAccount.Balance)
	require.Equal(t, account.Currency, dbAccount.Currency)
	require.WithinDuration(t, account.CreatedAt, account.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	account, err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	dbAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, dbAccount)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := AccountsListParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.AccountsList(context.Background(), args)

	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
