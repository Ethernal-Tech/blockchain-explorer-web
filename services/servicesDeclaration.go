package services

import (
	"webbc/models"
	"webbc/models/addressModel"
	"webbc/models/transactionModel"
)

type BlockService interface {
	GetLastBlocks(int) (*[]models.Block, error)
	GetBlockByNumber(uint64) (*models.Block, error)
	GetBlockByHash(string) (*models.Block, error)
}

type TransactionService interface {
	GetLastTransactions(int) (*[]models.Transaction, error)
	GetTransactionByHash(string) (*models.Transaction, error)
	GetAllTransactionsInBlock(uint64) (*[]models.Transaction, error)
	GetAllTransactions(int, int) (*transactionModel.Transactions, error)
}

type AddressService interface {
	GetAddress(string) (*addressModel.Address, error)
}
