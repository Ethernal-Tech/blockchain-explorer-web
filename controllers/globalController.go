package controllers

import (
	"strconv"
	"strings"
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type GlobalController struct {
	BlockService       services.BlockService
	TransactionService services.TransactionService
	AddressService     services.AddressService
}

func NewGlobalController(blockService services.BlockService, transactionService services.TransactionService, addressService services.AddressService) GlobalController {
	return GlobalController{BlockService: blockService, TransactionService: transactionService, AddressService: addressService}
}

func (gc GlobalController) GetIndex(context *gin.Context) {
	blocks, error1 := gc.BlockService.GetLastBlocks(20)
	transactions, error2 := gc.TransactionService.GetLastTransactions(20)

	if error1 != nil || error2 != nil {
		// return 500 or something like that
	}

	data := gin.H{
		"blocks":       blocks,
		"transactions": transactions,
	}

	context.HTML(200, "index.html", data)
}

func (gc GlobalController) GetBySearchValue(context *gin.Context) {
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
