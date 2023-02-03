package services

import (
	"webbc/models"
	"webbc/models/addressModel"
	"webbc/models/blockModel"
	"webbc/models/transactionModel"
)

type BlockService interface {
	GetLastBlocks(int) (*[]models.Block, error)
	GetBlockByNumber(uint64) (*models.Block, error)
	GetBlockByHash(string) (*models.Block, error)
	GetAllBlocks(int, int) (*blockModel.Blocks, error)
}

type TransactionService interface {
	GetLastTransactions(int) (*[]models.Transaction, error)
	GetTransactionByHash(string) (*models.Transaction, error)
	GetTransactionsInBlock(string, int, int) (*transactionModel.Transactions, error)
	GetAllTransactions(int, int) (*transactionModel.Transactions, error)
	GetTransactionsByAddress(string, int, int) (*transactionModel.Transactions, error)
}

type AddressService interface {
	GetAddress(string) (*addressModel.Address, error)
}
