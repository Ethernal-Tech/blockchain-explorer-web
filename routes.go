package main

import (
	"strings"
	"text/template"
	"webbc/configuration"
	"webbc/controllers"

	"github.com/gin-gonic/gin"
)

// Function for defining web server routes (endpoints). First, it is necessary to take the corresponding route group
// depending on which controller you are working with.
// Then, the route format (pattern) is as follows:
//
// /closer_description
//
// Or, if there is a need for a parameter, use:
//
// /closer_description/:value
//
// For example:
//
// /all
//
// /txinblock/:blocknumber
//
// If there is a need for more than one parameter, use the GET query parameters or the POST http method,
// but still register the route in one of the two ways defined previously.
// Parameters and potential conflicts in that case should be resolved at the controller level.
//
// For example, if there is a need for 'actionX' in 'controllerA' with three parameters, register one
// of the following routes (inside controllerA rout group):
//
// /actionX/:value1 will refer to /controllerA/actionX/:value1?parameter2=value2&parameter3=value3
//
// /actionX will refer to /controllerA/actionX?parameter1=value1&parameter2=value2&parameter3=value3
func routes(server *gin.Engine, cont ...any) {
	server.SetFuncMap(template.FuncMap{
		"getAppConfig": getAppConfig,
		"sub":          Subtraction,
		"startsWith":   StartsWith,
	})
	server.LoadHTMLGlob("staticfiles/*.html")
	server.Static("/images", "./staticfiles/images")
	server.Static("/staticfiles", "./staticfiles")

	server.GET("/css/style.css", getCssFile)

	authorized := server.Group("/", gin.BasicAuth(gin.Accounts{
		authConfig.Username: authConfig.Password,
	}))

	gc := cont[0].(controllers.GlobalController)
	server.GET("/", gc.GetIndex)
	server.GET("/:searchValue", gc.GetBySearchValue)
	authorized.GET("/configuration/:type", gc.GetConfiguration)
	authorized.POST("/configuration/:type", gc.UpdateConfiguration)

	bc := cont[1].(controllers.BlockController)
	blockRoutes := server.Group("/block")
	{
		blockRoutes.GET("/all", bc.GetBlocks)
		blockRoutes.GET("/number/:blocknumber", bc.GetBlockByNumber)
		blockRoutes.GET("/hash/:blockhash", bc.GetBlockByHash)
	}

	tc := cont[2].(controllers.TransactionController)
	transactionRoutes := server.Group("/transaction")
	{
		transactionRoutes.GET("/all", tc.GetTransactions)
		transactionRoutes.GET("/hash/:txhash", tc.GetTransactionByHash)
		transactionRoutes.GET("/address/:address", tc.GetTransactionsByAddress)
		transactionRoutes.GET("/txinblock/:blocknumber", tc.GetTransactionsInBlock)
	}

	ac := cont[3].(controllers.AddressController)
	addressRoutes := server.Group("/address")
	{
		addressRoutes.GET("/hash/:addresshash", ac.GetAddress)
		addressRoutes.POST("/uploadABI/:address", ac.UploadABI)
	}

	nc := cont[4].(controllers.NftController)
	nftRoutes := server.Group("/nft")
	{
		nftRoutes.GET("/transfers", nc.GetLatestTransfers)
		nftRoutes.GET("/:address/:tokenid", nc.GetNftMetadata)
		nftRoutes.GET("/transfers/:address/:tokenid", nc.GetNftTransfers)
	}
}

func getAppConfig() *configuration.ApplicationConfiguration {
	appConfig.Mutex.RLock()
	defer appConfig.Mutex.RUnlock()
	return appConfig
}

func getCssFile(c *gin.Context) {
	t, err := template.ParseFiles("staticfiles/css/style.css")

	if err != nil {
		// error handle
	}

	c.Header("Content-Type", "text/css")

	t.Execute(c.Writer, appConfig)
}

func Subtraction(a, b int) int {
	return a - b
}

func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}
