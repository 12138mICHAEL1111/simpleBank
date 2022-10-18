package api

import (
	"fmt"

	db "github.com/12138mICHAEL1111/simplebank/db/sqlc"
	"github.com/12138mICHAEL1111/simplebank/token"
	"github.com/12138mICHAEL1111/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config util.Config
	store *db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(config util.Config,store *db.Store) (*Server,error){
	tokenMaker,err := token.NewPasetoMaker(config.TokenSymmetricKey) //或者NewJWTMaker
	if err != nil {
		return nil,fmt.Errorf("cannot create token maker %w", err)
	}
	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}
	
	//对currency做校验
	if v,ok:= binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterValidation("currency",validCurrency)
	}
	server.setupRouter()
	return server,nil
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func (server *Server) setupRouter (){
	router := gin.Default()
	router.POST("/createUser",server.createUser)
	router.POST("/users/login",server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/createAccount",server.createAccount)
	
	authRoutes.GET("/getAccount/:id",server.getAccount)
	authRoutes.GET("/getAccounts",server.ListAccounts)

	authRoutes.POST("/createTransfer",server.createTransfer)
	server.router = router
}

func errorResponse(err error) gin.H{
	return gin.H{"error":err.Error()}
}