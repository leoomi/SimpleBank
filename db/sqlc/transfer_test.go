package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/leoomi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	return createTransferForAccount(t, account1.ID, account2.ID)
}

func createTransferForAccount(t *testing.T, fromID, toID int64) Transfer {
	params := CreateTransferParams{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, params.Amount, transfer.Amount)
	require.Equal(t, params.FromAccountID, transfer.FromAccountID)
	require.Equal(t, params.ToAccountID, transfer.ToAccountID)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func createMultipleRandomTransfers(t *testing.T, n int) (Account, Account, []Transfer) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	var transfers []Transfer

	for i := 0; i < n; i++ {
		e := createTransferForAccount(t, account1.ID, account2.ID)
		transfers = append(transfers, e)
	}

	return account1, account2, transfers
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	newTransfer := createRandomTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), newTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, newTransfer.ID, transfer.ID)
	require.Equal(t, newTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, newTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, newTransfer.Amount, transfer.Amount)
	require.Equal(t, newTransfer.CreatedAt, transfer.CreatedAt)
}

func TestUpdateTransfer(t *testing.T) {
	newTransfer := createRandomTransfer(t)

	params := UpdateTransferParams{
		ID:     newTransfer.ID,
		Amount: util.RandomMoney(),
	}
	transfer, err := testQueries.UpdateTransfer(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, newTransfer.ID, transfer.ID)
	require.Equal(t, newTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, newTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, params.Amount, transfer.Amount)
	require.Equal(t, newTransfer.CreatedAt, transfer.CreatedAt)
}

func TestDeleteTransfer(t *testing.T) {
	newTransfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), newTransfer.ID)
	require.NoError(t, err)

	entry, err := testQueries.GetTransfer(context.Background(), newTransfer.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry)
}

func TestListTransfers(t *testing.T) {
	createMultipleRandomTransfers(t, 10)

	params := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	entries, err := testQueries.ListTransfers(context.Background(), params)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, e := range entries {
		require.NotEmpty(t, e)
	}
}

func TestListTransfersFromAcount(t *testing.T) {
	account, _, _ := createMultipleRandomTransfers(t, 10)

	params := ListTransfersFromParams{
		FromAccountID: account.ID,
		Limit:         10,
		Offset:        0,
	}
	transfers, err := testQueries.ListTransfersFrom(context.Background(), params)

	require.NoError(t, err)
	require.Len(t, transfers, 10)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, account.ID, transfer.FromAccountID)
	}
}

func TestListTransfersToAcount(t *testing.T) {
	_, account, _ := createMultipleRandomTransfers(t, 10)

	params := ListTransfersToParams{
		ToAccountID: account.ID,
		Limit:       10,
		Offset:      0,
	}
	transfers, err := testQueries.ListTransfersTo(context.Background(), params)

	require.NoError(t, err)
	require.Len(t, transfers, 10)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, account.ID, transfer.ToAccountID)
	}
}
