package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type routineResult struct {
	err    error
	result TransferTxResult
}

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(5)

	resultChan := make(chan routineResult)

	n := 5
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			resultChan <- routineResult{err, result}
		}()
	}

	existed := map[int]bool{}
	for i := 0; i < n; i++ {
		result := <-resultChan
		require.NoError(t, result.err)
		require.NotEmpty(t, result.result)

		transfer := result.result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		_, err := store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diffFrom := account1.Balance - fromAccount.Balance
		diffTo := toAccount.Balance - account2.Balance
		require.Equal(t, diffFrom, diffTo)
		require.True(t, diffTo > 0)
		require.True(t, diffTo%amount == 0)
		require.True(t, diffFrom > 0)
		require.True(t, diffFrom%amount == 0)

		totalTransfers := int(diffFrom / amount)
		require.True(t, totalTransfers >= 1 && totalTransfers <= n)
		require.NotContains(t, existed, totalTransfers)
		existed[totalTransfers] = true
	}

	updatedFromAccount, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTransferTxAlternateSourceDestination(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(10)

	errs := make(chan error)

	n := 10
	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID
		if i%2 == 1 {
			fromAccountId, toAccountId = toAccountId, fromAccountId
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
