package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "simple_bank.sqlc.dev/app/db/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof= USD EUR"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListAccounts struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"` // specify the range of pageSize allowed to display
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// insert new account to database
	args := db.CreateAccountsParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccounts(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// account := db.Account{}

	ctx.JSON(http.StatusCreated, account)

}
func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest

	// validation the request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get account by id from param uri
	account, err := server.store.GetAccountById(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusFound, account)
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccounts
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
