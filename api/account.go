package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "simple_bank.sqlc.dev/app/db/sqlc"
	"simple_bank.sqlc.dev/app/token"
)

type CreateAccountRequest struct {
	// Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// insert new account to database
	args := db.CreateAccountsParams{
		// since payload as set in the request header
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccounts(ctx, args)

	if err != nil {
		if pgError, ok := err.(*pq.Error); ok {
			switch pgError.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// account := db.Account{}

	// return response code to the server
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get authenticated account compare with on-demand account query
	authenticated := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authenticated.Username != account.Owner {
		// create a custom err
		err := errors.New("Account doesn't belong to the authenticated token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccounts

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// extract payload from the request header

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.GetAccountsParams{
		Owner:  payload.Username,
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
