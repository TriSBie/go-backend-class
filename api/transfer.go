package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "simple_bank.sqlc.dev/app/db/sqlc"
	"simple_bank.sqlc.dev/app/token"
)

type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,gt=0"`
	Amount        int64  `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccountById(ctx, accountID)
	if err != nil {
		// if no row return
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s vs %s ", accountID, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req CreateTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// validate the currency type of each account

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)

	if !valid {
		return
	}

	if fromAccount.Owner != payload.Username {
		err := errors.New("account owner doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)

	if !valid {
		return
	}

	// insert new account to database
	args := db.CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.CreateTransfer(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// result := db.Result{}

	ctx.JSON(http.StatusCreated, result)
}

type GetTransferByIdParams struct {
	ID int64 `uri:"id"`
}

func (server *Server) getTransferById(ctx *gin.Context) {
	var req GetTransferByIdParams

	// if err is only single
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetTransferById(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}
