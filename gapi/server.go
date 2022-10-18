package gapi

import (
	"fmt"

	db "github.com/12138mICHAEL1111/simplebank/db/sqlc"
	"github.com/12138mICHAEL1111/simplebank/pb"
	"github.com/12138mICHAEL1111/simplebank/token"
	"github.com/12138mICHAEL1111/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	config util.Config
	store *db.Store
	tokenMaker token.Maker
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
	
	return server,nil
}