package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)
	accout1 := testCreateAccount(t)
	accout2 := testCreateAccount(t)
	type AT struct {
		Account
		TransferTxResult
	}
	amout := int64(100)
	errorChan := make(chan error)
	resultChan := make(chan AT)
	var n = 10
	for i := 0; i < n; i++ {
		go func() {
			arg := TransferTxParams{
				FromAccountID: accout1.ID,
				ToAccountID:   accout2.ID,
				Amount:        amout,
			}
			result, err := store.TransferTx(context.Background(), arg)
			errorChan <- err
			resultChan <- AT{
				Account:          accout1,
				TransferTxResult: result,
			}
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errorChan
		require.NoError(t, err)
		result := <-resultChan
		require.NotEmpty(t, result)

		//check Transfer
		require.Equal(t, result.Transfer.FromAccountID, accout1.ID)
		require.Equal(t, result.Transfer.ToAccountID, accout2.ID)
		require.Equal(t, result.Transfer.Amount, amout)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		arg := GetTransferParams{
			ID:       result.Transfer.ID,
			Username: result.Owner,
		}
		_, err = store.GetTransfer(context.Background(), arg)
		require.NoError(t, err)
		//check fromEntity
		fromEntity := result.FromEntry
		require.Equal(t, fromEntity.AccountID, accout1.ID)
		require.Equal(t, fromEntity.Amount, -amout)
		require.NotZero(t, fromEntity.ID)
		require.NotZero(t, fromEntity.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntity.ID)
		require.NoError(t, err)
		//check toEntity
		toEntity := result.ToEntry
		require.Equal(t, toEntity.AccountID, accout2.ID)
		require.Equal(t, toEntity.Amount, amout)
		require.NotZero(t, toEntity.ID)
		require.NotZero(t, toEntity.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntity.ID)
		require.NoError(t, err)

		//check Account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromEntity)
		require.Equal(t, fromAccount.ID, accout1.ID)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, accout2.ID)

		diff1 := accout1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - accout2.Balance
		require.Equal(t, diff1, diff2)
		require.NotZero(t, diff1)
		require.True(t, diff1%amout == 0)
	}
	resultAccount1, err := testQueries.GetAccountForUpdate(context.Background(), accout1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, resultAccount1)
	resultAccount2, err := testQueries.GetAccountForUpdate(context.Background(), accout2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, resultAccount2)
	require.Equal(t, resultAccount1.Balance, accout1.Balance-int64(n)*amout)
	require.Equal(t, resultAccount2.Balance, accout2.Balance+int64(n)*amout)
}

//测试死锁情况
func TestTransferTxDeadBlock(t *testing.T) {
	store := NewStore(testDB)
	accout1 := testCreateAccount(t)
	accout2 := testCreateAccount(t)
	amout := int64(100)
	errorChan := make(chan error)
	var n = 20
	for i := 0; i < n; i++ {
		fromAccountID := accout1.ID
		toAccountID := accout2.ID
		if i%2 != 0 {
			fromAccountID, toAccountID = toAccountID, fromAccountID
		}
		go func() {
			arg := TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amout,
			}
			_, err := store.TransferTx(context.Background(), arg)
			errorChan <- err
		}()
	}
	resultAccount1, err := testQueries.GetAccountForUpdate(context.Background(), accout1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, resultAccount1)
	resultAccount2, err := testQueries.GetAccountForUpdate(context.Background(), accout2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, resultAccount2)
	require.Equal(t, resultAccount1.Balance, accout1.Balance)
	require.Equal(t, resultAccount2.Balance, accout2.Balance)
}
