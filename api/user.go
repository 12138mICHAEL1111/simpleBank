package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/12138mICHAEL1111/simplebank/db/sqlc"
	"github.com/12138mICHAEL1111/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"` 
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UserResponse struct{
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	FullName          string    `json:"full_name"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

//增加一个
func (server *Server) createUser(ctx *gin.Context){
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	hashedPassword , err := util.HashPassword(req.Password)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username: req.Username,
		Email: req.Email,
		FullName: req.FullName,
		HashedPassword: hashedPassword,
	}

	user,err := server.store.CreateUser(ctx,arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error);ok {
			switch pqErr.Code.Name(){
			case "unique_violation":
				ctx.JSON(http.StatusForbidden,errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	res := UserResponse{
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		CreatedAt: user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
	ctx.JSON(http.StatusOK,res)
}

type loginUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"` 
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct{
	AccessToken string `json:"access_token"`
	User UserResponse `json:"user"`
}

func (server *Server) loginUser (ctx *gin.Context){
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err!= nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	
	user,err := server.store.GetUser(ctx,req.Username)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password,user.HashedPassword)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	
	accessToken,_,err := server.tokenMaker.CreateToken(user.Username,server.config.AccessTokenDuration)
	if err !=nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	userRes := UserResponse{
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		CreatedAt: user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
	res := loginUserResponse{
		AccessToken: accessToken,
		User: userRes,
	}
	ctx.JSON(http.StatusOK,res)
}