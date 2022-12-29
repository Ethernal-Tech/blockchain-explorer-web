package controllers

import (
	"strconv"
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type BlockController struct {
	BlockService services.BlockService
}

func NewBlockController(blockService services.BlockService) BlockController {
	return BlockController{BlockService: blockService}
}

func (bc BlockController) GetBlockByNumber(context *gin.Context) {
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

func (bc BlockController) GetBlockByHash(context *gin.Context) {
	_, er := bc.BlockService.GetBlockByHash(context.Param("blockhsh"))

	if er != nil {
	}

}
