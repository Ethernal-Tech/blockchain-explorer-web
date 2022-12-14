package main

import (
	"context"
	"webbc/DB"
	"webbc/configuration"
	"webbc/controllers"
	"webbc/services"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

var (
	ctx                   context.Context
	config                *configuration.Configuration
	database              *bun.DB
	server                *gin.Engine
	blockService          services.BlockService
	blockController       controllers.BlockController
	transactionService    services.TransactionService
	transactionController controllers.TransactionController
	globalController      controllers.GlobalController
)

func main() {
	ctx = context.TODO()
	config = configuration.LoadConfiguration()
	database = DB.InitializationDB(config)
	server = gin.Default()
	blockService = services.NewBlockService(database, ctx)
	blockController = controllers.NewBlockController(blockService)
	transactionService = services.NewTransactionService(database, ctx)
	transactionController = controllers.NewTransactionController(transactionService)
	globalController = controllers.NewGlobalController(blockService, transactionService)

	routes(server, globalController, blockController, transactionController)
	server.Run(":10000")
}
