package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rollbackErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"fromAccountID"`
	ToAccountID   int64 `json:"toAccountID"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"To_entry"`
}

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			var changeBalancesResult ChangeBalancesResult
			changeBalancesResult, err = ChangeBalances(ctx, q, ChangeBalancesParams{
				Account1ID: arg.FromAccountID,
				Amount1:    -arg.Amount,
				Account2ID: arg.ToAccountID,
				Amount2:    arg.Amount,
			})
			result.FromAccount = changeBalancesResult.Account1
			result.ToAccount = changeBalancesResult.Account2
		} else {
			var changeBalancesResult ChangeBalancesResult
			changeBalancesResult, err = ChangeBalances(ctx, q, ChangeBalancesParams{
				Account2ID: arg.FromAccountID,
				Amount2:    -arg.Amount,
				Account1ID: arg.ToAccountID,
				Amount1:    arg.Amount,
			})
			result.FromAccount = changeBalancesResult.Account2
			result.ToAccount = changeBalancesResult.Account1
		}

		if err != nil {
			return err
		}

		return err
	})

	return result, err
}

type ChangeBalancesParams struct {
	Account1ID int64
	Amount1    int64
	Account2ID int64
	Amount2    int64
}
type ChangeBalancesResult struct {
	Account1 Account
	Account2 Account
}

func ChangeBalances(ctx context.Context, q *Queries, params ChangeBalancesParams) (ChangeBalancesResult, error) {
	result := ChangeBalancesResult{}
	var err error

	result.Account1, err = q.AddToAccountBalance(ctx, AddToAccountBalanceParams{
		ID:     params.Account1ID,
		Amount: params.Amount1,
	})
	if err != nil {
		return ChangeBalancesResult{}, err
	}

	result.Account2, err = q.AddToAccountBalance(ctx, AddToAccountBalanceParams{
		ID:     params.Account2ID,
		Amount: params.Amount2,
	})

	if err != nil {
		return ChangeBalancesResult{}, err
	}

	return result, nil
}
