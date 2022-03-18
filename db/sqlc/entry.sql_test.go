package db

import (
	"context"
	"github.com/0RAJA/Bank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueries_CreateEntry(t *testing.T) {
	account := testCreateAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomInt(0, 1000),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)

	require.NotZero(t, entry.ID)
}

func testCreateEntry(t *testing.T) Entry {
	account := testCreateAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomInt(1, 1000),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)

	require.NotZero(t, entry.ID)
	return entry
}

func TestQueries_GetEntry(t *testing.T) {
	entry := testCreateEntry(t)
	require.NotZero(t, entry)
	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotZero(t, entry2)
	require.Equal(t, entry, entry2)
	require.WithinDuration(t, entry.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestQueries_ListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		testCreateEntry(t)
	}
	account := testCreateAccount(t)
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entrys, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	for _, entry := range entrys {
		require.NotZero(t, entry)
	}
}
