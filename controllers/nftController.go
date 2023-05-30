package controllers

import (
	"strconv"
	"webbc/models/paginationModel"
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type NftController struct {
	NftService services.INftService
}

func NewNftController(nftService services.INftService) NftController {
	return NftController{NftService: nftService}
}

func (nc *NftController) GetLatestTransfers(context *gin.Context) {
	page, perPage := paginationNftTransaction(context)

	result, err := nc.NftService.GetLatestTransfers(page, perPage)

	if err != nil {
		//TODO error handling
	}

	data := gin.H{
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
	context.HTML(200, "nftLatestTransfers.html", data)
}

func paginationNftTransaction(context *gin.Context) (int, int) {
	page := 1
	pageStr := context.Query("p")
	if pageStr != "" {
		if pageConv, err := strconv.Atoi(pageStr); err == nil && pageConv > 0 {
			page = pageConv
		}
	}
	perPage := 25
	perPageStr := context.Query("l")
	if perPageStr != "" {
		if perPageConv, err := strconv.Atoi(perPageStr); err == nil && perPageConv > 0 {
			perPage = perPageConv
		}
	}

	return page, perPage
}