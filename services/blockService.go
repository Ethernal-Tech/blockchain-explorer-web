package services

import (
	"context"
	"math"
	"time"
	"webbc/DB"
	"webbc/models"
	"webbc/models/blockModel"
	"webbc/utils"

	"github.com/uptrace/bun"
)

type BlockServiceImplementation struct {
	database *bun.DB
	ctx      context.Context
}

func NewBlockService(database *bun.DB, ctx context.Context) BlockService {
	return &BlockServiceImplementation{database: database, ctx: ctx}
}

func (bsi *BlockServiceImplementation) GetLastBlocks(numberOfBlocks int) (*[]models.Block, error) {
	var blocks []DB.Block
	error1 := bsi.database.NewSelect().Table("blocks").Order("number DESC").Limit(20).Scan(bsi.ctx, &blocks)

	if error1 != nil {

	}

	var result []models.Block

	for _, v := range blocks {
		var oneResultBlock models.Block = models.Block{
			Hash:               v.Hash,
			Number:             v.Number,
			ParentHash:         v.ParentHash,
			Nonce:              v.Nonce,
			Validator:          v.Miner,
			Difficulty:         v.Difficulty,
			TotalDifficulty:    v.TotalDifficulty,
			ExtraData:          v.ExtraData,
			Size:               v.Size,
			GasLimit:           v.GasLimit,
			GasUsed:            v.GasUsed,
			Age:                utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds()))),
			DateTime:           time.Unix(int64(v.Timestamp), 0).UTC().Format("2006-01-02 15:04:05"),
			TransactionsNumber: v.TransactionsCount,
		}

		result = append(result, oneResultBlock)
	}

	return &result, nil
}

func (bsi *BlockServiceImplementation) GetBlockByNumber(blockNumber uint64) (*models.Block, error) {
	var block DB.Block
	error1 := bsi.database.NewSelect().Table("blocks").Where("blocks.number = ?", blockNumber).Scan(bsi.ctx, &block)

	if error1 != nil {

	}

	var oneResultBlock models.Block = models.Block{
		Hash:               block.Hash,
		Number:             block.Number,
		ParentHash:         block.ParentHash,
		Nonce:              block.Nonce,
		Validator:          block.Miner,
		Difficulty:         block.Difficulty,
		TotalDifficulty:    utils.ToBigInt(block.TotalDifficulty).String(),
		ExtraData:          block.ExtraData,
		Size:               block.Size,
		GasLimit:           block.GasLimit,
		GasUsed:            block.GasUsed,
		Age:                utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(block.Timestamp), 0)).Seconds()))),
		DateTime:           time.Unix(int64(block.Timestamp), 0).UTC().Format("Jan-02-2006 15:04:05"),
		TransactionsNumber: block.TransactionsCount,
	}

	return &oneResultBlock, nil
}

func (bsi *BlockServiceImplementation) GetBlockByHash(blockHash string) (*models.Block, error) {
	var block DB.Block
	error1 := bsi.database.NewSelect().Table("blocks").Where("blocks.hash = ?", blockHash).Scan(bsi.ctx, &block)

	if error1 != nil {

	}

	var oneResultBlock models.Block = models.Block{
		Hash:               block.Hash,
		Number:             block.Number,
		ParentHash:         block.ParentHash,
		Nonce:              block.Nonce,
		Validator:          block.Miner,
		Difficulty:         block.Difficulty,
		TotalDifficulty:    utils.ToBigInt(block.TotalDifficulty).String(),
		ExtraData:          block.ExtraData,
		Size:               block.Size,
		GasLimit:           block.GasLimit,
		GasUsed:            block.GasUsed,
		Age:                utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(block.Timestamp), 0)).Seconds()))),
		DateTime:           time.Unix(int64(block.Timestamp), 0).UTC().Format("2006-01-02 15:04:05"),
		TransactionsNumber: block.TransactionsCount,
	}

	return &oneResultBlock, nil
}

func (bsi *BlockServiceImplementation) GetAllBlocks(page int, perPage int) (*blockModel.Blocks, error) {
	var blocks []DB.Block

	var offSet = perPage * (page - 1)
	err := bsi.database.NewSelect().Table("blocks").OrderExpr("number DESC").Limit(perPage).Offset(offSet).Scan(bsi.ctx, &blocks)
	if err != nil {
		//TODO: error handling
	}

	var result blockModel.Blocks
	for _, v := range blocks {
		var block = blockModel.Block{
			Number:             v.Number,
			Age:                utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds()))),
			DateTime:           time.Unix(int64(v.Timestamp), 0).UTC().Format("2006-01-02 15:04:05"),
			TransactionsNumber: v.TransactionsCount,
			Validator:          v.Miner,
			GasUsed:            v.GasUsed,
			GasLimit:           v.GasLimit,
		}

		result.Blocks = append(result.Blocks, block)
	}

	var totalRows int64
	bsi.database.NewRaw("SELECT count(*) FROM blocks").Scan(bsi.ctx, &totalRows)
	result.TotalRows = int64(totalRows)

	totalPages := math.Ceil(float64(totalRows) / float64(perPage))
	if totalPages == 0 {
		result.TotalPages = 1
	} else {
		result.TotalPages = int(totalPages)
	}
	return &result, nil
}
