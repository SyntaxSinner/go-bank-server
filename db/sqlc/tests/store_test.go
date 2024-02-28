package tests

import (
	"context"
	"testing"

	"github.com/SyntaxSinner/BankCRUD_API/db/sqlc/transactions"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(test *testing.T) {

	store := transactions.NewStore(test_db)

	account1 := RandomAccCreator(test)
	account2 := RandomAccCreator(test)

	n := 3
	amount := int64(10)

	errs := make(chan error)
	results := make(chan transactions.TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), transactions.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {

		//test errors
		err := <-errs
		require.NoError(test, err)

		//result not empty
		result := <-results
		require.NotEmpty(test, result)

		//test the transfer object
		transfer := result.Transfer
		require.NotEmpty(test, transfer)
		require.Equal(test, account1.ID, transfer.FromAccountID.Int64)
		require.Equal(test, account2.ID, transfer.ToAccountID.Int64)
		require.Equal(test, amount, transfer.Amount)
		require.NotZero(test, transfer.ID)
		require.NotZero(test, transfer.CreatedAt)

		//test the entries

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(test, err) //check wether the transfer exists

		from_entry := result.FromEntry
		require.NotEmpty(test, from_entry)
		require.Equal(test, account1.ID, from_entry.AccountID.Int64)
		require.Equal(test, -amount, from_entry.Balance)
		require.NotZero(test, from_entry.ID)
		require.NotZero(test, from_entry.CreatedAt)

		_, err = store.GetEntry(context.Background(), from_entry.ID)
		require.NoError(test, err) //check wether the entry exists

		to_entry := result.ToEntry
		require.NotEmpty(test, to_entry)
		require.Equal(test, account2.ID, to_entry.AccountID.Int64)
		require.Equal(test, amount, to_entry.Balance)
		require.NotZero(test, to_entry.ID)
		require.NotZero(test, to_entry.CreatedAt)

		_, err = store.GetEntry(context.Background(), to_entry.ID)
		require.NoError(test, err) //check wether the entry exists

		//test the accounts
		from_account := result.FromAccount
		require.NotEmpty(test, from_account)
		require.Equal(test, account1.ID, from_account.ID)
		require.Equal(test, account1.Balance-amount, from_account.Balance)

		to_account := result.ToAccount
		require.NotEmpty(test, to_account)
		require.Equal(test, account2.ID, to_account.ID)
		require.Equal(test, account2.Balance+amount, to_account.Balance)

	}

}
