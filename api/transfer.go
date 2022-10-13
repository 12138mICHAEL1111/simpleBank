package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/12138mICHAEL1111/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID  int64 `json:"from_account_id" binding:"required,min=1"` 
	ToAccountID int64 `json:"to_account_id" binding:"required,min=1"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createTransfer(ctx *gin.Context){
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	if ! server.validAccount(ctx,req.FromAccountID,req.Currency){
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}

	account,err := server.store.TransferTx(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}

func (server *Server) validAccount(ctx *gin.Context,accountID int64,currency string) bool {
	account, err := server.store.GetAccount(ctx,accountID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return false
	} 

	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mistatch: %s vs %s", accountID,account.Currency,currency)
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return false
	}
	return true
}