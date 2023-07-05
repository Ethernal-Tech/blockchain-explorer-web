package controllers

import (
	"strconv"
	"webbc/models/paginationModel"
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type BlockController struct {
	BlockService services.IBlockService
}

func NewBlockController(blockService services.IBlockService) BlockController {
	return BlockController{BlockService: blockService}
}

func (bc *BlockController) GetBlockByNumber(context *gin.Context) {
	blockNumber, error1 := strconv.ParseUint(context.Param("blocknumber"), 10, 64)

	if error1 != nil {

	}

	block, error2 := bc.BlockService.GetBlockByNumber(blockNumber)

	if error2 != nil {

	}

	data := gin.H{
		"block": block,
	}

	context.HTML(200, "block.html", data)
}

func (bc *BlockController) GetBlockByHash(context *gin.Context) {
	block, error1 := bc.BlockService.GetBlockByHash(context.Param("blockhash"))

	if error1 != nil {
	}

	data := gin.H{
		"block": block,
	}

	context.HTML(200, "block.html", data)
}

func (bc BlockController) GetBlocks(context *gin.Context) {
	page, perPage := PaginationBlock(context)

	result, err := bc.BlockService.GetAllBlocks(page, perPage)

	if err != nil {
		//TODO error handling
	}

	pag := paginationModel.PaginationData{
		NextPage:     page + 1,
		PreviousPage: page - 1,
		CurrentPage:  page,
		TotalPages:   result.TotalPages,
		TotalRows:    result.TotalRows,
		PerPage:      perPage,
	}

	data := gin.H{
		"blocks":     result.Blocks,
		"pagination": pag,
		"last":       len(result.Blocks) - 1,
	}

	context.HTML(200, "blocks.html", data)
}

func PaginationBlock(context *gin.Context) (int, int) {
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
