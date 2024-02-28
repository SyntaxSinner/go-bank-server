package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/SyntaxSinner/BankCRUD_API/db/sqlc"
	"github.com/SyntaxSinner/BankCRUD_API/db/sqlc/utils"
	"github.com/stretchr/testify/require"
)

func RandomAccCreator(test *testing.T) sqlc.Account {
	arg := utils.To_account(utils.GenerateRandomOwner())
	fmt.Print(arg)
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
	test.Log("Account Test - Starting CreateAccount test...")
	RandomAccCreator(test)
}

func TestGetAccount(test *testing.T) {
	test.Log("Account Test - Starting GetAccount test...")

	account := RandomAccCreator(test)
	account, err := test_queries.GetAccount(context.Background(), account.ID)
	require.NoError(test, err)
	require.NotEmpty(test, account)

	require.Equal(test, account.ID, account.ID)
	require.Equal(test, account.Owner, account.Owner)
	require.Equal(test, account.Balance, account.Balance)
	require.Equal(test, account.Currency, account.Currency)
	require.NotZero(test, account.CreatedAt)
	require.WithinDuration(test, account.CreatedAt, time.Now(), time.Second)

}

func TestDeleteAccount(test *testing.T) {
	test.Log("Account Test - Starting DeleteAccount test...")

	account := RandomAccCreator(test)
	err := test_queries.DeleteAccount(context.Background(), account.ID)
	require.NoError(test, err)

	account, err = test_queries.GetAccount(context.Background(), account.ID)
	require.Error(test, err)
	require.Empty(test, account)
}
