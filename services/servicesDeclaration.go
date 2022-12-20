package services

import (
	"webbc/models"
	"webbc/models/addressModel"
)

type BlockService interface {
	GetLastBlocks(int) (*[]models.Block, error)
	GetBlockByNumber(uint64) (*models.Block, error)
}

type TransactionService interface {
	GetLastTransactions(int) (*[]models.Transaction, error)
	GetTransactionByHash(string) (*models.Transaction, error)
}

type AddressService interface {
	GetAddress(string) (*addressModel.Address, error)
}
