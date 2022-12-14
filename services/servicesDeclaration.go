package services

import (
	"webbc/models"
)

type BlockService interface {
	GetLastBlocks(int) (*[]models.Block, error)
	GetBlockByNumber(uint64) (*models.Block, error)
}

type TransactionService interface {
	GetLastTransactions(int) (*[]models.Transaction, error)
	GetTransactionByHash(string) (*models.Transaction, error)
}
