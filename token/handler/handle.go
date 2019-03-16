package handler

import (
	"github.com/gin-gonic/gin"
	"token/token"
	"net/http"
	"log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"strings"
	"math/big"
	"token/conf"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
)

type Response struct {
	Status int
	Code int
	Msg string
	Data interface{}
}

type TransferData struct {
	Address string 	`json:"address"`
	Amount  int64 	`json:"amount"`
}

type BalanceData struct {
	Address string `json:"address"`
}

var (
	connecter *token.Connecter
)

func init() {
	connecter = token.NewConnecter()
}


func GetBlockNumber(c *gin.Context)  {
	var r Response
	result := connecter.BlockNumber()
	r.Status, r.Code, r.Msg, r.Data = 1, 1, "获取区块高度成功", result
	c.JSON(http.StatusOK, r)
}

func GetTokenName(c *gin.Context)  {
	var r Response
	result := connecter.ContractName()
	r.Status, r.Code, r.Msg, r.Data = 1, 1, "获取代币名称成功", result
	c.JSON(http.StatusOK, r)
}

func GetTokenTotal(c *gin.Context)  {
	var r Response
	result := connecter.TotalSupply()
	r.Status, r.Code, r.Msg, r.Data = 1, 1, "获取代币总量成功", result
	c.JSON(http.StatusOK, r)
}

func Transfer(c *gin.Context)  {

	var r Response
	var data TransferData
	err := c.BindJSON(&data)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewTransactor(strings.NewReader(conf.KEY), conf.PWD)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	result := connecter.TransferToken(auth, common.HexToAddress(data.Address), big.NewInt(data.Amount))
	if !result {
		r.Status, r.Code, r.Msg = 0, 0, "转账失败"
		c.JSON(http.StatusOK, r)
	}

	r.Status, r.Code, r.Msg, r.Data = 1, 1, "转账成功", "成功转账"+strconv.FormatInt(data.Amount, 10)
	c.JSON(http.StatusOK, r)
}

func Balance(c *gin.Context)  {
	var r Response
	var data BalanceData
	err := c.BindJSON(&data)
	if err != nil {
		log.Fatal(err)
	}

	result := connecter.BalanceOfToken(common.HexToAddress(data.Address))
	r.Status, r.Code, r.Msg, r.Data = 1, 1, "查询余额成功", result
	c.JSON(http.StatusOK, r)
}