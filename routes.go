package main

import (
	"text/template"
	"webbc/configuration"
	"webbc/controllers"

	"github.com/gin-gonic/gin"
)

// Function for defining web server routes (endpoints). The route format (pattern) is as follows:
//
// /controller/closer_description
//
// Or, if there is a need for a parameter, use:
//
// /controller/closer_description/:value
//
// For example:
//
// /block/all
//
// /transaction/txinblock/:blocknumber
//
// If there is a need for more than one parameter, use the GET query parameters or the POST http method,
// but still register the route in one of the two ways defined previously.
// Parameters and potential conflicts in that case should be resolved at the controller level.
//
// For example, if there is a need for 'actionX' in 'controllerA' with three parameters, register one
// of the following routes:
//
// /controllerA/actionX/:value1 will refer to /controllerA/actionX/:value1?parameter2=value2&parameter3=value3
//
// /controllerA/actionX will refer to /controllerA/actionX?parameter1=value1&parameter2=value2&parameter3=value3
func routes(server *gin.Engine, cont ...any) {
	server.SetFuncMap(template.FuncMap{
		"getConfig": getConfig,
	})
	server.LoadHTMLGlob("staticfiles/*.html")
	server.Static("/images", "./staticfiles/images")

	server.GET("/css/style.css", getCssFile)

	gc := cont[0].(controllers.GlobalController)
	server.GET("/", gc.GetIndex)
	server.GET("/:searchValue", gc.GetBySearchValue)

	bc := cont[1].(controllers.BlockController)
	server.GET("/block/all", bc.GetBlocks)
	server.GET("/block/number/:blocknumber", bc.GetBlockByNumber)
	server.GET("/block/hash/:blockhash", bc.GetBlockByHash)

	tc := cont[2].(controllers.TransactionController)
	server.GET("/transaction/all", tc.GetTransactions)
	server.GET("/transaction/hash/:txhash", tc.GetTransactionByHash)
	server.GET("/transaction/address/:address", tc.GetTransactionsByAddress)
	server.GET("/transaction/txinblock/:blocknumber", tc.GetTransactionsInBlock)

	ac := cont[3].(controllers.AddressController)
	server.GET("/address/hash/:addresshash", ac.GetAddress)
}

func getConfig() *configuration.TemplateConfiguration {
	return templateConfig
}

func getCssFile(c *gin.Context) {
	t, err := template.ParseFiles("staticfiles/css/style.css")

	if err != nil {
		// error handle
	}

	c.Header("Content-Type", "text/css")

	t.Execute(c.Writer, templateConfig)
}
