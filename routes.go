package main

import (
	"text/template"
	"webbc/controllers"

	"github.com/gin-gonic/gin"
)

func routes(server *gin.Engine, cont ...any) {
	server.LoadHTMLGlob("staticfiles/*.html")
	server.Static("/images", "./staticfiles/images")

	server.GET("/css/style.css", getCssFile)

	gc := cont[0].(controllers.GlobalController)
	server.GET("/", gc.GetIndex)
	server.GET("/:searchValue", gc.GetBySearchValue)

	bc := cont[1].(controllers.BlockController)
	server.GET("/block/:blocknumber", bc.GetBlockByNumber)
	server.GET("/blockhash/:blockhsh", bc.GetBlockByHash)
	server.GET("/blocks", bc.GetBlocks)

	tc := cont[2].(controllers.TransactionController)
	server.GET("/transaction/:transactionhx", tc.GetTransactionByHash)
	server.GET("/transactionsinblock/:blocknumber", tc.GetAllTransactionsInBlock)
	server.GET("/txs", tc.GetTransactions)
	server.GET("/txsByAddress", tc.GetTransactionsByAddress)

	ac := cont[3].(controllers.AddressController)
	server.GET("/address/:addresshex", ac.GetAddress)
}

func getCssFile(c *gin.Context) {
	t, err := template.ParseFiles("staticfiles/css/style.css")

	if err != nil {
		// error handle
	}

	c.Header("Content-Type", "text/css")

	t.Execute(c.Writer, config)
}
