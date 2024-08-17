package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "simple_bank.sqlc.dev/app/db/sqlc"
	"simple_bank.sqlc.dev/app/token"
	"simple_bank.sqlc.dev/app/util"
)

type Server struct {
	config     util.Config
	store      db.Store    // self-contained state management
	router     *gin.Engine // router from gin
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	// create new tokenMaker
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %v", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: maker,
		config:     config,
	}

	// since gin using validator v10 under the hood -
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validator - to tag name "currency"
		// Using tag at struct validation
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouter()

	return server, nil
}

// splitting server router more smaller
func (server *Server) setUpRouter() {
	// define router in gin
	router := gin.Default()
	router.GET("/users", server.getUser)
	router.POST("/users/login", server.LoginUser)

	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouter.POST("/accounts", server.createAccount)
	authRouter.POST("/transfers", server.createTransfer)
	authRouter.GET("/transfers/:id", authMiddleware(server.tokenMaker), server.getTransferById)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.listAccounts)
	authRouter.POST("/users", server.createUser)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// define errorResponse handler return map[string]any type
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
