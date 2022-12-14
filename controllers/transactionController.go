package controllers

import (
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService services.TransactionService
}

func NewTransactionController(transactionService services.TransactionService) TransactionController {
	return TransactionController{TransactionService: transactionService}
}

func (gc TransactionController) GetTransactionByHash(context *gin.Context) {
	transaction, error := gc.TransactionService.GetTransactionByHash(context.Param("transactionhx"))

	if error != nil {

	}

	// var transactions *[]models.Transaction = &[]models.Transaction{}
	// *transactions = append(*transactions, *transaction)

	data := gin.H{
		"transaction": transaction,
	}

	context.HTML(200, "transaction.html", data)
}
