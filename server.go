package main

import (
	"context"
	"webbc/DB"
	"webbc/configuration"
	"webbc/controllers"
	"webbc/eth"
	"webbc/services"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

var (
	ctx                   context.Context
	config                *configuration.Configuration
	templateConfig        *configuration.TemplateConfiguration
	database              *bun.DB
	server                *gin.Engine
	connection            eth.BlockchainNodeConnection
	blockService          services.BlockService
	blockController       controllers.BlockController
	transactionService    services.TransactionService
	transactionController controllers.TransactionController
	globalController      controllers.GlobalController
	addressService        services.AddressService
	addressController     controllers.AddressController
)

func main() {
	ctx = context.TODO()
	config, templateConfig = configuration.LoadConfiguration()
	database = DB.InitializationDB(config)
	server = gin.Default()
	connection := eth.BlockchainNodeConnection{
		HTTP: eth.GetClient(config.HTTPUrl),
	}

	blockService = services.NewBlockService(database, ctx)
	blockController = controllers.NewBlockController(blockService)
	transactionService = services.NewTransactionService(database, ctx)
	transactionController = controllers.NewTransactionController(transactionService)
	addressService = services.NewAddressService(database, ctx, connection.HTTP, config)
	addressController = controllers.NewAddressController(addressService)
	globalController = controllers.NewGlobalController(blockService, transactionService, addressService)

	routes(server, globalController, blockController, transactionController, addressController)
	server.Run(":10010")
}
