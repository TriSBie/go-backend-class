package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "simple_bank.sqlc.dev/app/db/sqlc"
)

type Server struct {
	store  db.Store    // self-contained state management
	router *gin.Engine // router from gin
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	// define router in gin
	router := gin.Default()

	// since gin using validator v10 under the hood -
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validator - to tag name "currency"
		// Using tag at struct validation
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.POST("/transfers", server.createTransfer)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router

	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// define errorResponse handler return map[string]any type
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
