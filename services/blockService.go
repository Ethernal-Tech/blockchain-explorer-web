package services

import (
	"context"
	"math"
	"time"
	"webbc/DB"
	"webbc/models"
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

	var transactionNumberByBlock []models.TransactionNumberByBlock
	error2 := bsi.database.NewRaw("select block_number, count(*) as count from transactions where block_number in (select blocks.number from blocks order by blocks.number DESC limit 20) group by block_number").Scan(bsi.ctx, &transactionNumberByBlock)

	if error2 != nil {

	}

	var mapTransactionNumberByBlock map[uint64]int = map[uint64]int{}

	for _, v := range transactionNumberByBlock {
		mapTransactionNumberByBlock[v.Block_number] = v.Count
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
			Timestamp:          utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds()))),
			TransactionsNumber: mapTransactionNumberByBlock[v.Number],
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

	var transactionNumberInBlock models.TransactionNumberByBlock
	error2 := bsi.database.NewRaw("select count(*) as count from transactions where block_number = ?", block.Number).Scan(bsi.ctx, &transactionNumberInBlock)

	if error2 != nil {

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
		Timestamp:          utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(block.Timestamp), 0)).Seconds()))),
		TransactionsNumber: transactionNumberInBlock.Count,
	}

	return &oneResultBlock, nil
}

func (bsi *BlockServiceImplementation) GetBlockByHash(blockHash string) (*models.Block, error) {
	var block DB.Block
	error1 := bsi.database.NewSelect().Table("blocks").Where("blocks.hash = ?", blockHash).Scan(bsi.ctx, &block)

	if error1 != nil {

	}

	var transactionNumberInBlock models.TransactionNumberByBlock
	error2 := bsi.database.NewRaw("select count(*) as count from transactions where block_number = ?", block.Number).Scan(bsi.ctx, &transactionNumberInBlock)

	if error2 != nil {

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
		Timestamp:          utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(block.Timestamp), 0)).Seconds()))),
		TransactionsNumber: transactionNumberInBlock.Count,
	}

	return &oneResultBlock, nil
}
