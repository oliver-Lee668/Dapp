package router

import (
	"github.com/gin-gonic/gin"
	"token/handler"
)

func Router(r *gin.Engine)  {
	r.GET("smartcontract/blocknumber", handler.GetBlockNumber)
	r.GET("smartcontract/tokenname", handler.GetTokenName)
	r.GET("smartcontract/tokentotal", handler.GetTokenTotal)
	r.POST("smartcontract/transfertoken", handler.Transfer)
	r.POST("smartcontract/balanceoftoken", handler.Balance)
}