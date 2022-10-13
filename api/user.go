package api

import (
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
	type createUserReponse struct{
		Username          string    `json:"username"`
		Email             string    `json:"email"`
		FullName          string    `json:"full_name"`
		PasswordChangedAt time.Time `json:"password_changed_at"`
		CreatedAt         time.Time `json:"created_at"`
	}
	res := createUserReponse{
		Username: user.Username,
		Email:user.Email,
		FullName: user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK,res)
}
