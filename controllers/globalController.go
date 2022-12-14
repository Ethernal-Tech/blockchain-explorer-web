package controllers

import (
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type GlobalController struct {
	BlockService       services.BlockService
	TransactionService services.TransactionService
}

func NewGlobalController(blockService services.BlockService, transactionService services.TransactionService) GlobalController {
	return GlobalController{BlockService: blockService, TransactionService: transactionService}
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
