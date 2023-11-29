package nftModel

import (
	"webbc/models/paginationModel"
)

type NftMetadata struct {
	Id      int64
	TokenId string
	Address string
	Name    string
	Image   string
	Owner   string
	Creator string
}

type NftTransfers struct {
	TotalPages     int
	TotalRows      int64
	Transfers      []NftTransfer
	MaxCount       int
	PaginationData paginationModel.PaginationData
}

type NftTransfer struct {
	TransactionHash string
	Method          string
	Age             string
	DateTime        string
	From            string
	To              string
	Type            string
	Item            string
	IsFromContract  bool
	IsToContract    bool
	TokenId         string
	Value           string
	Address         string
	NftName         string
	NftImage        string
}
