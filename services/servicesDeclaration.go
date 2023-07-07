package services

import (
	"webbc/configuration"
	"webbc/models/addressModel"
	"webbc/models/blockModel"
	"webbc/models/nftModel"
	"webbc/models/transactionModel"

	"github.com/gin-gonic/gin"
)

type IBlockService interface {
	GetLastBlocks(int) (*[]blockModel.Block, error)
	GetBlockByNumber(uint64) (*blockModel.Block, error)
	GetBlockByHash(string) (*blockModel.Block, error)
	GetAllBlocks(int, int) (*blockModel.Blocks, error)
}

type ITransactionService interface {
	GetLastTransactions(int) (*[]transactionModel.Transaction, error)
	GetTransactionByHash(string) (*transactionModel.Transaction, error)
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
	GetLatestTransfers(int, int) (*nftModel.NftTransfers, error)
	GetNftMetadata(string, string) (*nftModel.NftMetadata, error)
	GetNftTransfers(string, string, int, int) (*nftModel.NftTransfers, error)
}
