package services

import (
	"webbc/configuration"
	"webbc/models"
	"webbc/models/addressModel"
	"webbc/models/blockModel"
	"webbc/models/nftModel"
	"webbc/models/transactionModel"

	"github.com/gin-gonic/gin"
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

type IAddressService interface {
	GetAddress(string) (*addressModel.Address, error)
	UploadABI(*gin.Context)
}

type ConfigurationService interface {
	GetAppConfiguration() *configuration.ApplicationConfiguration
	GetGeneralConfiguration() *configuration.GeneralConfiguration
	UpdateAppConfiguration(*configuration.ApplicationConfiguration) error
	UpdateGeneralConfiguration(*configuration.GeneralConfiguration) error
}

type INftService interface {
	GetLatestTransfers(int, int) (*nftModel.NftTransactions, error)
}
