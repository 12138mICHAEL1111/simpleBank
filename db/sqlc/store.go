package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct{
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store{
	return &Store{
		db:db,
		Queries: New(db),
	}
}

// encapulate the rollback and commit method
func (store *Store) execTx(ctx context.Context,fn func(*Queries) error) error{
	tx, err := store.db.BeginTx(ctx,nil)
	if err != nil{
		return err 
	}

	q:= New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, rb err: %v",err,rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxResult struct{
	Transfer Transfer 
	FromAccount Account
	ToAccount Account
	FromEntry Entry
	ToEntry Entry
}


func (store *Store) TransferTx(ctx context.Context,arg CreateTransferParams) (TransferTxResult,error){
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx,arg)
		if err != nil{
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		if arg.FromAccountID < arg.ToAccountID { //防止出现死锁，例如事务1 ，account1更新，account2更新，事务2，account2更新，account1更新就会出现死锁
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err

	})
	return result, err
}


func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}