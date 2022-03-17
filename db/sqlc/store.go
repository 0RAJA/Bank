package db

import (
	"context"
	"database/sql"
	"fmt"
)

/*
	数据库事务ACID
	A:原子性:全部完成，或者全部失败回滚
	C:一致性:写入的数据必须正确
	I:隔离性:并发事务不应该相互影响
	D:持久性:即使出现异常也应该持久化
*/

type Store struct {
	*Queries
	db *sql.DB
}

// NewStore 返回一个查询对象
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

//通过事务执行回调函数
func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) //开启事务
	if err != nil {
		return err
	}
	q := New(tx) //使用开启的事务创建一个查询
	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); err != nil { //回滚
			return fmt.Errorf("tx err:%v,rb err:%v", err, rbErr)
		}
		return err
	}
	return tx.Commit() //提交
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id,omitempty"`
	ToAccountID   int64 `json:"to_account_id,omitempty"`
	Amount        int64 `json:"amount,omitempty"`
}

type TransferTxRequest struct {
	Transfer    Transfer `json:"transfer"`     // 交易记录
	FromAccount Account  `json:"from_account"` //源账户
	ToAccount   Account  `json:"to_account"`   //目标账户
	FromEntry   Entry    `json:"from_entry"`   //源交易后账户
	ToEntry     Entry    `json:"to_entry"`     //目标交易后账户
}

// TransferTx 通过事务进行交易执行
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxRequest, error) {
	var result TransferTxRequest
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		//创建一个转移记录
		result.Transfer, err = store.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		//给源账户创建一个减少的账户条目
		result.FromEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		//给目标账户创建一个增加的账户条目
		result.ToEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//更新账户信息 可能出现死锁现象，原因是交易顺序不同，可能会导致锁无法释放
		//规定更新默认以同样的方式进行获取锁
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addAmount(ctx, queries, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = addAmount(ctx, queries, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

func addAmount(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return Account{}, Account{}, err
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return Account{}, Account{}, err
	}
	return account1, account2, nil
}
