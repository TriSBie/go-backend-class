package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"simple_bank.sqlc.dev/app/util"
)

func createRandomAccount(t *testing.T) Account {
	ctx := context.Background()

	user := createRandomUser(t)
	arg := CreateAccountsParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// using context with testQueries (*Queries) initialized from main test
	account, err := testQueries.CreateAccounts(ctx, arg)

	// Ensure the execution run without nil or error
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Owner, account1.Owner)
	require.Equal(t, account2.Currency, account1.Currency)
	require.Equal(t, account2.Balance, account1.Balance)

	// ensure two time of assert is within a delta duration
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	// initialize the params of arg
	arg := UpdateAccountBalanceParams{
		Balance: util.RandomMoney(),
		ID:      account1.ID,
	}

	updateRes, err := testQueries.UpdateAccountBalance(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updateRes)

	require.Equal(t, account1.ID, updateRes.ID)
	require.Equal(t, arg.Balance, updateRes.Balance)
	require.Equal(t, account1.Currency, updateRes.Currency)
	require.Equal(t, account1.Owner, updateRes.Owner)

	// ensure two time of assert is within a delta duration
	require.WithinDuration(t, account1.CreatedAt, updateRes.CreatedAt, time.Second*1)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	// initialize the params of arg
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccountById(context.Background(), account.ID)

	// Ensure error occurs
	require.Error(t, err)
	// Ensure error No Rows found are caught
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		// create 10 random accounts & assign account to the last reference to test
		lastAccount = createRandomAccount(t)
	}

	args := GetAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.GetAccounts(context.Background(), args)

	require.NoError(t, err)
	// require.Len(t, accounts, 5)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
