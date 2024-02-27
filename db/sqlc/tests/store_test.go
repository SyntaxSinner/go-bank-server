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

	}

}
