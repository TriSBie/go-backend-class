package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store interface {
	Querier
	TransferTxDeadLock(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Store provides all functions to execute db queries and transactions
// By using embedded Queries inside Store struct, all methods or functions will be inherited when using with Store
type StoreDB struct {
	*Queries // inheritance
	db       *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &StoreDB{
		db:      db,
		Queries: New(db),
	}
}

// execTX executes a function within a database transaction and return an error if occur
// define second parameters as callback function which accepts references Queries and return error
func (store *StoreDB) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return err
	}

	q := New(tx)

	// calling execution callback function
	err = fn(q)
	if err != nil {
		// rollBack a transaction if any error occur - could throw an error while rollback execution
		if rbErr := tx.Rollback(); rbErr != nil {
			// %v print the value in default format
			return fmt.Errorf("tx Error: %v, rbError : %v", err, rbErr)
		}
		log.Printf("Error %v\n", err)
		return err
	}

	// if no any exception occurs -> commit the transaction
	return tx.Commit()
}

// contains the input parameters of the transfer function
type TransferTxParams struct {
	FromAccountID int64 `json:from_account_id`
	ToAccountID   int64 `json:to_account_id`
	Amount        int64 `json:amount`
}

// contains the results of the transfers function
type TransferTxResult struct {
	Transfer    Transfer `json:transfer`
	FromAccount Account  `json:from_account`
	ToAccount   Account  `json:to_account`
	FromEntry   Entry    `json:from_entry`
	ToEntry     Entry    `json:to_entry`
}

// TransferTx performs a money transfer from one account to another
// It creates a record, add account entries, and update account's balance with a single database transaction
func (store *StoreDB) TransferTxDeadLock(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTX(ctx, func(q *Queries) error {
		var err error

		transferArg := CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		}

		// #1 - Create new transfer of amount transaction between A-B
		// after createTransfer execute successfully -> create new entry for logging amount of a specific account
		result.Transfer, err = q.CreateTransfer(ctx, transferArg)

		if err != nil {
			log.Printf("Error occurred: %v", err)
			return err
		}

		fromEntryArg := CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}

		// #2 - Stores the transition income of sender
		// Logging the amount issues by sender  sends to receivers
		result.FromEntry, err = q.CreateEntry(ctx, fromEntryArg)

		if err != nil {
			log.Printf("Error occurred: %v", err)
			return err
		}

		toEntryArg := CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}

		// #3 - Stores the transition income of receiver
		// Logging the amount income issues from sender
		result.ToEntry, err = q.CreateEntry(ctx, toEntryArg)

		if err != nil {
			log.Printf("Error occurred: %v", err)
			return err
		}

		// TODO: update the balance of both sender and receiver
		// Ensure the account query getting from the valid transaction phrases.

		// #4.1 Find information of sender account and update account balance
		// result.FromAccount, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)

		// if err != nil {
		// 	log.Printf("Error occurred: %v", err)
		// 	return err
		// }

		// // #4.1 Find information of sender account and update account balance
		// result.ToAccount, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)

		// if err != nil {
		// 	log.Printf("Error occurred: %v", err)
		// 	return err
		// }

		// ensure the sequential of update operations should be asc by id
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(context.Background(),
				q,
				arg.FromAccountID, -arg.Amount,
				arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(context.Background(),
				q,
				arg.ToAccountID, arg.Amount,
				arg.FromAccountID, -arg.Amount,
			)
		}

		return nil
	})

	return result, err
}

// Refactoring the function
func addMoney(
	ctx context.Context,
	q *Queries,
	fromAccountId int64,
	amount1 int64,
	toAccountId int64,
	amount2 int64,
) (account1 Account,
	account2 Account,
	err error,
) {
	account1, err = q.AddMoneyToAccount(ctx, AddMoneyToAccountParams{
		Balance: amount1,
		ID:      fromAccountId,
	})

	if err != nil {
		return
	}

	account2, err = q.AddMoneyToAccount(ctx, AddMoneyToAccountParams{
		Balance: amount2,
		ID:      toAccountId,
	})

	if err != nil {
		return
	}
	return account1, account2, err
}
