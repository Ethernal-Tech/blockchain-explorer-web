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
	transaction, error := tc.TransactionService.GetTransactionByHash(context.Param("txhash"))

	if error != nil {
	}

	data := gin.H{
		"transaction": transaction,
	}

	context.HTML(200, "transaction.html", data)
}

func (tc *TransactionController) GetTransactionsInBlock(context *gin.Context) {
	block := context.Param("blocknumber")
	page, perPage := PaginationTransaction(context)

	result, err := tc.TransactionService.GetTransactionsInBlock(block, page, perPage)

	if err != nil {
		//TODO error handling
	}

	data := gin.H{
		"block":        block,
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
	context.HTML(200, "transactionsInBlock.html", data)
}

func (tc TransactionController) GetTransactions(context *gin.Context) {
	page, perPage := PaginationTransaction(context)

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

func (tc TransactionController) GetTransactionsByAddress(context *gin.Context) {
	address := context.Param("address")

	page, perPage := PaginationTransaction(context)
	result, err := tc.TransactionService.GetTransactionsByAddress(address, page, perPage)

	if err != nil {
		//TODO error handling
	}

	data := gin.H{
		"address":      address,
		"transactions": result.Transactions,
		"txsMaxCount":  result.MaxCount,
		"pagination": paginationModel.PaginationData{
			NextPage:     page + 1,
			PreviousPage: page - 1,
			CurrentPage:  page,
			TotalPages:   result.TotalPages,
			TotalRows:    result.TotalRows,
			PerPage:      perPage,
		},
	}
	context.HTML(200, "transactionsByAddress.html", data)
}

func PaginationTransaction(context *gin.Context) (int, int) {
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
	return page, perPage
}
