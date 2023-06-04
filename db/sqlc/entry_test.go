package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/leoomi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	return createEntryForAccount(t, account.ID)
}

func createEntryForAccount(t *testing.T, accountID int64) Entry {
	params := CreateEntryParams{
		AccountID: accountID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, params.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func createMultipleRandomEntries(t *testing.T, n int) (Account, []Entry) {
	account := createRandomAccount(t)
	var entries []Entry

	for i := 0; i < n; i++ {
		e := createEntryForAccount(t, account.ID)
		entries = append(entries, e)
	}

	return account, entries
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	newEntry := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), newEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, newEntry.ID, entry.ID)
	require.Equal(t, newEntry.AccountID, entry.AccountID)
	require.Equal(t, newEntry.Amount, entry.Amount)
	require.Equal(t, newEntry.CreatedAt, entry.CreatedAt)
}

func TestUpdateEntry(t *testing.T) {
	newEntry := createRandomEntry(t)

	params := UpdateEntryParams{
		ID:     newEntry.ID,
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.UpdateEntry(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, newEntry.ID, entry.ID)
	require.Equal(t, newEntry.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)
	require.Equal(t, newEntry.CreatedAt, entry.CreatedAt)
}

func TestDeleteEntry(t *testing.T) {
	newEntry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), newEntry.ID)
	require.NoError(t, err)

	entry, err := testQueries.GetEntry(context.Background(), newEntry.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry)
}

func TestListEntries(t *testing.T) {
	createMultipleRandomEntries(t, 10)

	params := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(), params)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, e := range entries {
		require.NotEmpty(t, e)
	}
}

func TestListEntriesFromAccount(t *testing.T) {
	account, _ := createMultipleRandomEntries(t, 10)

	params := ListEntriesFromAccountParams{
		AccountID: account.ID,
		Limit:     10,
		Offset:    0,
	}
	entries, err := testQueries.ListEntriesFromAccount(context.Background(), params)

	require.NoError(t, err)
	require.Len(t, entries, 10)

	for _, e := range entries {
		require.NotEmpty(t, e)
		require.Equal(t, account.ID, e.AccountID)
	}
}
