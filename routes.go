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
	// server.GET("/block/:blocknumber", bc.GetBlockByNumber)
	// server.GET("/blockhash/:blockhsh", bc.GetBlockByHash)
	// server.GET("/blocks", bc.GetBlocks)
	server.GET("/block/all", bc.GetBlocks)
	server.GET("/block/number/:blocknumber", bc.GetBlockByNumber)
	server.GET("/block/hash/:blockhash", bc.GetBlockByHash)

	tc := cont[2].(controllers.TransactionController)
	// server.GET("/transaction/:transactionhx", tc.GetTransactionByHash)
	// server.GET("/transactionsinblock/:blocknumber", tc.GetAllTransactionsInBlock)
	// server.GET("/txs", tc.GetTransactions)
	// server.GET("/txsByAddress", tc.GetTransactionsByAddress)
	server.GET("/transaction/all", tc.GetTransactions)
	server.GET("/transaction/hash/:txhash", tc.GetTransactionByHash)
	server.GET("/transaction/address/:address", tc.GetTransactionsByAddress)
	server.GET("/transaction/txinblock/:blocknumber", tc.GetAllTransactionsInBlock)

	ac := cont[3].(controllers.AddressController)
	server.GET("/address/hash/:addresshash", ac.GetAddress)
}

func getCssFile(c *gin.Context) {
	t, err := template.ParseFiles("staticfiles/css/style.css")

	if err != nil {
		// error handle
	}

	c.Header("Content-Type", "text/css")

	t.Execute(c.Writer, config)
}
