package transactions

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SyntaxSinner/BankCRUD_API/db/sqlc"
)

// In order to execute the queries, we need to create a store that will hold the queries and the database connection
type Store struct {
	*sqlc.Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: sqlc.New(db),
		db:      db,
	}
}

// ExecTx executes a function within a transaction

func (store *Store) ExecTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		rollback_error := tx.Rollback()
		return fmt.Errorf("transaction error: %v, rollback error: %v", err, rollback_error)
	}
	// If the function returns no error, commit the transaction
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	FromEntry   sqlc.Entry    `json:"from_entry"`
	ToEntry     sqlc.Entry    `json:"to_entry"`
	Transfer    sqlc.Transfer `json:"transfer"`
	FromAccount sqlc.Account  `json:"from_account"`
	ToAccount   sqlc.Account  `json:"to_account"`
}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.ExecTx(ctx, func(q *sqlc.Queries) error {

		var err error
		//retrieve the sender account from the database to check if the account exists
		from_account, err := q.GetAccount(ctx, args.FromAccountID)
		if err != nil {
			return err
		}
		//retrieve the receiver account from the database to check if the account exists
		to_account, err := q.GetAccount(ctx, args.ToAccountID)
		if err != nil {
			return err
		}
		// Check wether we can move forward with the transaction in case the sender has the amount to send
		amount := int64(args.Amount)
		if from_account.Balance < amount {
			return fmt.Errorf("insufficient balance")
		}

		//once the condition is passed we're going to create the transfer record in the database , 2 entries one for each account and update the balance of each acccounts
		// 1) create 2 entries

		result.ToEntry, err = q.CreateEntry(ctx, sqlc.CreateEntryParams{
			AccountID: sql.NullInt64{Int64: args.ToAccountID, Valid: true},
			Balance:   args.Amount,
			Currency:  from_account.Currency})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, sqlc.CreateEntryParams{
			AccountID: sql.NullInt64{Int64: args.FromAccountID, Valid: true},
			Balance:   -args.Amount,
			Currency:  from_account.Currency})
		if err != nil {
			return err
		}

		// 2) create the transfer record
		result.Transfer, err = q.CreateTransfer(ctx, sqlc.CreateTransferParams{
			FromAccountID: sql.NullInt64{Int64: args.FromAccountID, Valid: true}, //we set valid to true since we proved that both of the accounts exist
			ToAccountID:   sql.NullInt64{Int64: args.ToAccountID, Valid: true},
			Amount:        args.Amount})
		if err != nil {
			return err
		}

		// 3) update the account balance
		result.FromAccount, err = q.UpdateAccount(ctx, sqlc.UpdateAccountParams{
			ID:      args.FromAccountID,
			Balance: from_account.Balance - args.Amount})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, sqlc.UpdateAccountParams{
			ID:      args.ToAccountID,
			Balance: to_account.Balance + args.Amount})
		if err != nil {
			return err
		}

		return nil

	})

	return result, err
}
