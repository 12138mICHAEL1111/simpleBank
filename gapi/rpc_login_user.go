package gapi

import (
	"context"
	"database/sql"

	"github.com/12138mICHAEL1111/simplebank/pb"
	"github.com/12138mICHAEL1111/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {	
	user,err := server.store.GetUser(ctx,req.GetUsername())
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, status.Errorf(codes.NotFound,"user not found")
		}
		return nil, status.Errorf(codes.Internal,"failed to find user")
	}

	err = util.CheckPassword(req.Password,user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.NotFound,"incorrect password")
	}
	
	accessToken,_,err := server.tokenMaker.CreateToken(user.Username,server.config.AccessTokenDuration)

	if err !=nil {
		return nil, status.Errorf(codes.Internal,"failed to create toeken")
	}

	res := &pb.LoginUserResponse{
		User: convertUser(user),
		AccessToken: accessToken,
	}
	return  res,nil

}