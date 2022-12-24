package controllers

import (
	"strconv"
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService services.TransactionService
}

func NewTransactionController(transactionService services.TransactionService) TransactionController {
	return TransactionController{TransactionService: transactionService}
}

func (tc TransactionController) GetTransactionByHash(context *gin.Context) {
	transaction, error := tc.TransactionService.GetTransactionByHash(context.Param("transactionhx"))

	if error != nil {

	}

	data := gin.H{
		"transaction": transaction,
	}

	context.HTML(200, "transaction.html", data)
}

func (tc TransactionController) GetAllTransactionsInBlock(context *gin.Context) {
	blockNumber, error1 := strconv.ParseUint(context.Param("blocknumber"), 10, 64)

	if error1 != nil {

	}

	transactions, error2 := tc.TransactionService.GetAllTransactionsInBlock(blockNumber)

	if error2 != nil {

	}

	data := gin.H{
		"block":        blockNumber,
		"transactions": transactions,
	}

	context.HTML(200, "transactionsInBlock.html", data)
}
