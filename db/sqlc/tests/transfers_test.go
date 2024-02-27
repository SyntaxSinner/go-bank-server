package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/SyntaxSinner/BankCRUD_API/db/sqlc"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func RandomTransferCreator(test *testing.T, fromAccountID, toAccountID int64) sqlc.Transfer {
	var arg sqlc.CreateTransferParams
	gofakeit.Struct(&arg)
	arg.FromAccountID = sql.NullInt64{Int64: fromAccountID, Valid: true}
	arg.ToAccountID = sql.NullInt64{Int64: toAccountID, Valid: true}
	transfer, err := test_queries.CreateTransfer(context.Background(), arg)

	require.NoError(test, err)
	require.NotEmpty(test, transfer)

	require.Equal(test, arg.FromAccountID.Int64, transfer.FromAccountID.Int64)
	require.Equal(test, arg.ToAccountID.Int64, transfer.ToAccountID.Int64)
	require.Equal(test, arg.Amount, transfer.Amount)
	require.Equal(test, arg.Currency, transfer.Currency)

	require.NotZero(test, transfer.ID)
	require.NotZero(test, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(test *testing.T) {
	test.Log("Transfers Test - Starting CreateTransfer test...")
	fromAccount := RandomAccCreator(test)
	toAccount := RandomAccCreator(test)
	RandomTransferCreator(test, fromAccount.ID, toAccount.ID)
}

func TestGetTransfer(test *testing.T) {
	test.Log("Transfers Test - Starting GetTransfer test...")
	fromAccount := RandomAccCreator(test)
	toAccount := RandomAccCreator(test)
	transfer := RandomTransferCreator(test, fromAccount.ID, toAccount.ID)

	retrievedTransfer, err := test_queries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(test, err)
	require.NotEmpty(test, retrievedTransfer)

	require.Equal(test, transfer.ID, retrievedTransfer.ID)
	require.Equal(test, transfer.FromAccountID.Int64, retrievedTransfer.FromAccountID.Int64)
	require.Equal(test, transfer.ToAccountID.Int64, retrievedTransfer.ToAccountID.Int64)
	require.Equal(test, transfer.Amount, retrievedTransfer.Amount)
	require.Equal(test, transfer.Currency, retrievedTransfer.Currency)
	require.NotZero(test, retrievedTransfer.CreatedAt)
	require.WithinDuration(test, retrievedTransfer.CreatedAt, time.Now(), 2*time.Hour)
}

func TestDeleteTransfer(test *testing.T) {
	test.Log("Transfers Test - Starting DeleteTransfer test...")
	fromAccount := RandomAccCreator(test)
	toAccount := RandomAccCreator(test)
	transfer := RandomTransferCreator(test, fromAccount.ID, toAccount.ID)

	err := test_queries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(test, err)

	retrievedTransfer, err := test_queries.GetTransfer(context.Background(), transfer.ID)
	require.Error(test, err)
	require.Empty(test, retrievedTransfer)
}
