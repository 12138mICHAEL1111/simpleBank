package api

import(
	"github.com/gin-gonic/gin"
	db "github.com/12138mICHAEL1111/simplebank/db/sqlc"
)

type Server struct {
	store *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server{
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/createAccount",server.createAccount)
	router.GET("/getAccount/:id",server.getAccount)
	router.GET("/getAccounts",server.ListAccounts)
	server.router = router
	return server
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func errorResponse(err error) gin.H{
	return gin.H{"error":err.Error()}
}