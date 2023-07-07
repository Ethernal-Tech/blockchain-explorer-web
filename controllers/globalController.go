package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"webbc/configuration"
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type GlobalController struct {
	BlockService         services.IBlockService
	TransactionService   services.ITransactionService
	AddressService       services.IAddressService
	ConfigurationService services.ConfigurationService
}

func NewGlobalController(blockService services.IBlockService, transactionService services.ITransactionService, addressService services.IAddressService, configurationService services.ConfigurationService) GlobalController {
	return GlobalController{BlockService: blockService, TransactionService: transactionService, AddressService: addressService, ConfigurationService: configurationService}
}

func (gc *GlobalController) GetIndex(context *gin.Context) {
	blocks, error1 := gc.BlockService.GetLastBlocks(20)
	transactions, error2 := gc.TransactionService.GetLastTransactions(20)

	if error1 != nil || error2 != nil {
		// return 500 or something like that
	}

	cookie, err := context.Cookie("timestamp")

	if err != nil {
		cookie = "NotSet"
		currentDate := time.Now()
		yearFromNow := currentDate.AddDate(1, 0, 0)
		dateSubtraction := yearFromNow.Sub(currentDate)
		context.SetCookie("timestamp", "Age", int(dateSubtraction.Seconds()), "/", "", false, false)
	}

	fmt.Printf("Cookie value: %s \n", cookie)

	data := gin.H{
		"blocks":       blocks,
		"transactions": transactions,
	}

	context.HTML(200, "index.html", data)
}

func (gc *GlobalController) GetBySearchValue(context *gin.Context) {
	searchValue := context.Param("searchValue")
	searchValueLen := len(searchValue)

	if strings.Contains(searchValue, "0x") {
		if searchValueLen == 42 {
			address, _ := gc.AddressService.GetAddress(searchValue)
			data := gin.H{
				"address": address,
			}
			context.HTML(200, "address.html", data)
		} else if searchValueLen == 66 {
			transaction, _ := gc.TransactionService.GetTransactionByHash(searchValue)
			data := gin.H{
				"transaction": transaction,
			}

			context.HTML(200, "transaction.html", data)
		}
	} else {
		blockNumber, error1 := strconv.ParseUint(searchValue, 10, 64)
		if error1 != nil {
		}

		block, _ := gc.BlockService.GetBlockByNumber(blockNumber)
		data := gin.H{
			"block": block,
		}

		context.HTML(200, "block.html", data)
	}
}

func (gc *GlobalController) GetConfiguration(context *gin.Context) {
	configType := context.Param("type")
	switch {
	case configType == "app":
		data := gc.ConfigurationService.GetAppConfiguration()
		context.HTML(200, "appConfiguration.html", data)
	case configType == "general":
		data := gc.ConfigurationService.GetGeneralConfiguration()
		context.HTML(200, "generalConfiguration.html", data)
	}
}

func (gc *GlobalController) UpdateConfiguration(context *gin.Context) {
	configType := context.Param("type")
	switch {
	case configType == "app":
		var appConfig *configuration.ApplicationConfiguration
		if err := context.Bind(&appConfig); err != nil {
			return
		}
		if err := gc.ConfigurationService.UpdateAppConfiguration(appConfig); err != nil {

		}
	case configType == "general":
		var generalConfig *configuration.GeneralConfiguration
		if err := context.Bind(&generalConfig); err != nil {
			return
		}
		if err := gc.ConfigurationService.UpdateGeneralConfiguration(generalConfig); err != nil {

		}
	}

	context.Redirect(302, context.Request.Referer())
}
