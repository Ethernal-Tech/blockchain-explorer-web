package main

import (
	"bytes"
	"text/template"
	"webbc/controllers"

	"github.com/gin-gonic/gin"
)

func routes(server *gin.Engine, cont ...any) {
	server.LoadHTMLGlob("staticfiles/*.html")
	//server.Static("/css", "./staticfiles/css")
	server.GET("/css/style.css", func(c *gin.Context) {
		t, _ := template.ParseFiles("staticfiles/css/style.css")
		var doc bytes.Buffer
		t.Execute(&doc, &config)
		c.Header("Content-Type", "text/css")
		c.Writer.Write(doc.Bytes())
	})

	server.Static("/images", "./staticfiles/images")

	gc := cont[0].(controllers.GlobalController)
	server.GET("/", gc.GetIndex)
	server.GET("/:searchValue", gc.GetBySearchValue)
	bc := cont[1].(controllers.BlockController)
	server.GET("/block/:blocknumber", bc.GetBlockByNumber)
	server.GET("/blockhash/:blockhsh", bc.GetBlockByHash)
	tc := cont[2].(controllers.TransactionController)
	server.GET("/transaction/:transactionhx", tc.GetTransactionByHash)
	server.GET("/transactionsinblock/:blocknumber", tc.GetAllTransactionsInBlock)
	server.GET("/txs", tc.GetTransactions)
	ac := cont[3].(controllers.AddressController)
	server.GET("/address/:addresshex", ac.GetAddress)
}
