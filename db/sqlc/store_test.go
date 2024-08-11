package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestTransferTransaction(t *testing.T) {
// 	store := NewStore(testDb)

// 	existed := make(map[int]bool)

// 	account1 := createRandomAccount(t)
// 	account2 := createRandomAccount(t)

// 	fmt.Println("Before ", account1.Balance, account2.Balance)

// 	// generate channel for storing the results and errors emitted from goroutine
// 	errors := make(chan error)
// 	results := make(chan TransferTxResult)

// 	n := 2
// 	amount := int64(10) //using typecasting

// 	for i := 0; i < n; i++ {
// 		// run concurrent go-routines
// 		go func() {
// 			result, err := store.TransferTxDeadLock(context.Background(), TransferTxParams{
// 				FromAccountID: account1.ID,
// 				ToAccountID:   account2.ID,
// 				Amount:        amount,
// 			})
// 			errors <- err
// 			results <- result
// 		}()
// 	}

// 	// Iterating through the channel and handle inside
// 	for i := 0; i < n; i++ {
// 		err := <-errors
// 		require.NoError(t, err)

// 		result := <-results
// 		require.NotEmpty(t, result)

// 		// operator <- only use for channel instead of value
// 		transfer := result.Transfer

// 		require.Equal(t, account1.ID, transfer.FromAccountID)
// 		require.Equal(t, account2.ID, transfer.ToAccountID)
// 		require.Equal(t, amount, transfer.Amount)
// 		require.NotZero(t, transfer.ID)
// 		require.NotZero(t, transfer.CreatedAt)

// 		// test get transfer by id
// 		_, err = testQueries.GetTransferById(context.Background(), transfer.ID)
// 		require.NoError(t, err)

// 		// test get entry from both sender and receiver
// 		fromEntry := result.FromEntry
// 		require.NotEmpty(t, fromEntry)
// 		require.Equal(t, account1.ID, fromEntry.AccountID)
// 		require.Equal(t, -amount, fromEntry.Amount)

// 		_, err = testQueries.GetEntryById(context.Background(), fromEntry.ID)
// 		require.NoError(t, err)

// 		toEntry := result.ToEntry
// 		require.NotEmpty(t, toEntry)
// 		require.Equal(t, account2.ID, toEntry.AccountID)
// 		require.Equal(t, amount, toEntry.Amount)

// 		_, err = testQueries.GetEntryById(context.Background(), toEntry.ID)
// 		require.NoError(t, err)

// 		// get account balance from sender
// 		fromAccount := result.FromAccount

// 		require.NotEmpty(t, fromAccount)
// 		require.Equal(t, account1.ID, fromAccount.ID)

// 		_, err = testQueries.GetAccountById(context.Background(), fromAccount.ID)
// 		require.NoError(t, err)

// 		// get account balance from receiver
// 		toAccount := result.ToAccount

// 		require.NotEmpty(t, toAccount)
// 		require.Equal(t, account2.ID, toAccount.ID)

// 		_, err = testQueries.GetAccountById(context.Background(), toAccount.ID)
// 		require.NoError(t, err)

// 		// check account balance
// 		diff1 := account1.Balance - fromAccount.Balance
// 		diff2 := toAccount.Balance - account2.Balance

// 		require.Equal(t, diff1, diff2)
// 		require.True(t, diff1 > 0)
// 		require.True(t, diff1%amount == 0) // check if the balance is divisible by the amount

// 		k := int(diff1 / amount)          // k must be the number of transactions
// 		require.True(t, k >= 1 && k <= n) // k must be greater than 1 and less than n
// 		// ensure that the key is unique
// 		require.NotContains(t, existed, k)
// 		existed[k] = true
// 	}

// 	updateAccount1, err := testQueries.GetAccountById(context.Background(), account1.ID)
// 	require.NoError(t, err)

// 	updateAccount2, err := testQueries.GetAccountById(context.Background(), account2.ID)
// 	require.NoError(t, err)

// 	fmt.Println("After: ", updateAccount1.Balance, updateAccount2.Balance)

// 	// require.Equal(t, account1.Balance-amount, updateAccount1.Balance)
// 	// require.Equal(t, account2.Balance+amount, updateAccount2.Balance)
// }

func TestTransferTransactionDeadlock(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("Before ", account1.Balance, account2.Balance)

	// generate channel for storing the results and errors emitted from goroutine
	errors := make(chan error)

	n := 10
	amount := int64(10) //using typecasting

	// Generate 10 iteration loops with accountA ( 5 times receiving & sending amount ) as well as with accountB
	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID
		// run concurrent go-routines
		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}
		go func() {
			_, err := store.TransferTxDeadLock(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})
			errors <- err
		}()
	}

	// Iterating through the channel and handle inside
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

	}

	updateAccount1, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("After: ", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
