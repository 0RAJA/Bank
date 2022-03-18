package db

import (
	"context"
	"github.com/0RAJA/Bank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueries_CreateTransfer(t *testing.T) {
	account1 := testCreateAccount(t)
	account2 := testCreateAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomInt(0, 1000),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, transfer)
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)
}

func testCreateTransfer(t *testing.T) (Account, Account, Transfer) {
	account1 := testCreateAccount(t)
	account2 := testCreateAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomInt(0, 1000),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.FromAccountID, account1.ID)
	require.Equal(t, transfer.ToAccountID, account2.ID)
	require.Equal(t, transfer.Amount, arg.Amount)
	return account1, account2, transfer
}

func TestQueries_GetTransfer(t *testing.T) {
	account1, _, transfer := testCreateTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), GetTransferParams{
		ID:       transfer.ID,
		Username: account1.Owner,
	})
	require.NoError(t, err)
	require.NotZero(t, transfer2)
	require.Equal(t, transfer, transfer2)
	require.WithinDuration(t, transfer.CreatedAt, transfer2.CreatedAt, time.Second)
}
