package api

import (
	db "github.com/12138mICHAEL1111/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server{
	server := &Server{store: store}
	router := gin.Default()
	//对currency做校验
	if v,ok:= binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterValidation("currency",validCurrency)
	}
	router.POST("/createUser",server.createUser)
	router.POST("/createAccount",server.createAccount)
	router.GET("/getAccount/:id",server.getAccount)
	router.GET("/getAccounts",server.ListAccounts)
	router.POST("/createTransfer",server.createTransfer)
	server.router = router
	return server
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func errorResponse(err error) gin.H{
	return gin.H{"error":err.Error()}
}