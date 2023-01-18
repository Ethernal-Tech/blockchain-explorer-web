package controllers

import (
	"strconv"
	"webbc/models/paginationModel"
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService services.TransactionService
}

func NewTransactionController(transactionService services.TransactionService) TransactionController {
	return TransactionController{TransactionService: transactionService}
}

func (tc *TransactionController) GetTransactionByHash(context *gin.Context) {
	transaction, error := tc.TransactionService.GetTransactionByHash(context.Param("transactionhx"))

	if error != nil {

	}

	data := gin.H{
		"transaction": transaction,
	}

	context.HTML(200, "transaction.html", data)
}

func (tc *TransactionController) GetAllTransactionsInBlock(context *gin.Context) {
	blockNumber, error1 := strconv.ParseUint(context.Param("blocknumber"), 10, 64)

	if error1 != nil {

	}

	transactions, error2 := tc.TransactionService.GetAllTransactionsInBlock(blockNumber)

	if error2 != nil {

	}

	data := gin.H{
		"blockNumber":          blockNumber,
		"numberOfTransactions": len(*transactions),
		"transactions":         transactions,
	}

	context.HTML(200, "transactionsInBlock.html", data)
}

func (tc TransactionController) GetTransactions(context *gin.Context) {
	address := context.Query("a")
	block := context.Query("block")

	if address != "" {
		//get transactions by address
	} else if block != "" {
		//get transactions by block
	} else {
		page := 1
		pageStr := context.Query("p")
		if pageStr != "" {
			page, _ = strconv.Atoi(context.Query("p"))
		}
		perPage := 50
		perPageStr := context.Query("l")
		if perPageStr != "" {
			perPage, _ = strconv.Atoi(context.Query("l"))
		}

		result, err := tc.TransactionService.GetAllTransactions(page, perPage)

		if err != nil {
			//TODO error handling
		}

		data := gin.H{
			"transactions": result.Transactions,
			"pagination": paginationModel.PaginationData{
				NextPage:     page + 1,
				PreviousPage: page - 1,
				CurrentPage:  page,
				TotalPages:   result.TotalPages,
				TotalRows:    result.TotalRows,
				PerPage:      perPage,
			},
		}
		context.HTML(200, "transactions.html", data)
	}
}
