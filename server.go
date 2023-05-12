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
	generalConfig         *configuration.GeneralConfiguration
	appConfig             *configuration.ApplicationConfiguration
	authConfig            *configuration.AuthConfiguration
	database              *bun.DB
	server                *gin.Engine
	blockService          services.BlockService
	blockController       controllers.BlockController
	transactionService    services.TransactionService
	transactionController controllers.TransactionController
	globalController      controllers.GlobalController
	addressService        services.IAddressService
	addressController     controllers.AddressController
	configurationService  services.ConfigurationService
)

func main() {
	ctx = context.TODO()
	generalConfig, appConfig, authConfig = configuration.LoadConfiguration()
	database = DB.InitializationDB(generalConfig)
	server = gin.Default()
	eth.NewHttpNodeClient(generalConfig.HTTPUrl)

	blockService = services.NewBlockService(database, ctx)
	blockController = controllers.NewBlockController(blockService)
	transactionService = services.NewTransactionService(database, ctx, generalConfig)
	transactionController = controllers.NewTransactionController(transactionService)
	addressService = services.NewAddressService(database, ctx, generalConfig)
	addressController = controllers.NewAddressController(addressService)
	configurationService = services.NewConfigurationService(appConfig, generalConfig, database, addressService)
	globalController = controllers.NewGlobalController(blockService, transactionService, addressService, configurationService)

	routes(server, globalController, blockController, transactionController, addressController)
	server.Run(":8080")
}
