package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/leoomi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	params := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), newAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, newAccount.ID, account.ID)
	require.Equal(t, newAccount.Owner, account.Owner)
	require.Equal(t, newAccount.Balance, account.Balance)
	require.Equal(t, newAccount.Currency, account.Currency)
	require.Equal(t, newAccount.CreatedAt, account.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	params := UpdateAccountParams{
		ID:      newAccount.ID,
		Balance: util.RandomMoney(),
	}
	account, err := testQueries.UpdateAccount(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, newAccount.ID, account.ID)
	require.Equal(t, newAccount.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, newAccount.Currency, account.Currency)
	require.Equal(t, newAccount.CreatedAt, account.CreatedAt)
}

func TestDeleteAccont(t *testing.T) {
	newAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), newAccount.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), newAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	params := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), params)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
