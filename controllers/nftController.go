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
		"transfers":   result.Transfers,
		"txsMaxCount": result.MaxCount,
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

func (nc *NftController) GetNftMetadata(context *gin.Context) {
	nftMetadata, error := nc.NftService.GetNftMetadata(context.Param("address"), context.Param("tokenid"))

	if error != nil {

	}

	data := gin.H{
		"nftMetadata": nftMetadata,
	}

	context.HTML(200, "nftMetadata.html", data)
}

func (nc *NftController) GetNftTransfers(context *gin.Context) {
	page, perPage := paginationNftTransaction(context)

	nftTransfers, error := nc.NftService.GetNftTransfers(context.Param("address"), context.Param("tokenid"), page, perPage)
	nftTransfers.PaginationData = paginationModel.PaginationData{
		NextPage:     page + 1,
		PreviousPage: page - 1,
		CurrentPage:  page,
		TotalPages:   nftTransfers.TotalPages,
		TotalRows:    nftTransfers.TotalRows,
		PerPage:      perPage,
	}

	if error != nil {

	}

	data := gin.H{
		"nftTransfers": nftTransfers,
	}

	context.HTML(200, "nftMetadataTransfers.html", data)
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
