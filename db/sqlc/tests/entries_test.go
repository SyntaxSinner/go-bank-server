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

func RandomEntryCreator(test *testing.T, accountID int64) sqlc.Entry {
	var arg sqlc.CreateEntryParams
	gofakeit.Struct(&arg)
	arg.AccountID = sql.NullInt64{Int64: accountID, Valid: true}
	entry, err := test_queries.CreateEntry(context.Background(), arg)
	require.NoError(test, err)
	require.NotEmpty(test, entry)

	require.Equal(test, arg.AccountID, entry.AccountID)
	require.Equal(test, arg.Balance, entry.Balance)
	require.Equal(test, arg.Currency, entry.Currency)

	require.NotZero(test, entry.ID)
	require.NotZero(test, entry.CreatedAt)
	return entry
}

func TestCreateEntry(test *testing.T) {
	test.Log("Entry Test - Starting CreateEntry test...")
	account := RandomAccCreator(test)
	RandomEntryCreator(test, account.ID)
}

func TestGetEntry(test *testing.T) {
	test.Log("Entry Test - Starting GetEntry test...")
	account := RandomAccCreator(test)
	entry := RandomEntryCreator(test, account.ID)

	retrievedEntry, err := test_queries.GetEntry(context.Background(), entry.ID)
	require.NoError(test, err)
	require.NotEmpty(test, retrievedEntry)

	require.Equal(test, entry.ID, retrievedEntry.ID)
	require.Equal(test, entry.AccountID, retrievedEntry.AccountID)
	require.Equal(test, entry.Balance, retrievedEntry.Balance)
	require.Equal(test, entry.Currency, retrievedEntry.Currency)
	require.NotZero(test, retrievedEntry.CreatedAt)
	require.WithinDuration(test, retrievedEntry.CreatedAt, time.Now(), time.Second)
}

func TestDeleteEntry(test *testing.T) {
	test.Log("Entry Test - Starting DeleteEntry test...")
	account := RandomAccCreator(test)
	entry := RandomEntryCreator(test, account.ID)

	err := test_queries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(test, err)

	retrievedEntry, err := test_queries.GetEntry(context.Background(), entry.ID)
	require.Error(test, err)
	require.Empty(test, retrievedEntry)
}
