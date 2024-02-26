package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func RandomAccCreator(test *testing.T) Account {
	var arg CreateAccountParams
	gofakeit.Struct(&arg)
	account, err := test_queries.CreateAccount(context.Background(), arg)

	require.NoError(test, err)
	require.NotEmpty(test, account)

	require.Equal(test, arg.Owner, account.Owner)
	require.Equal(test, arg.Balance, account.Balance)
	require.Equal(test, arg.Currency, account.Currency)

	require.NotZero(test, account.ID)
	require.NotZero(test, account.CreatedAt)
	return account
}

func TestCreateAccount(test *testing.T) {
	RandomAccCreator(test)
}

func TestGetAccount(test *testing.T) {
	account := RandomAccCreator(test)
	account, err := test_queries.GetAccount(context.Background(), account.ID)
	require.NoError(test, err)
	require.NotEmpty(test, account)

	require.Equal(test, account.ID, account.ID)
	require.Equal(test, account.Owner, account.Owner)
	require.Equal(test, account.Balance, account.Balance)
	require.Equal(test, account.Currency, account.Currency)
	require.NotZero(test, account.CreatedAt)

}
