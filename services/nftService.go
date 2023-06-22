package services

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"
	"webbc/DB"
	"webbc/configuration"
	"webbc/models"
	"webbc/models/nftModel"
	"webbc/utils"

	"github.com/uptrace/bun"
)

type NftService struct {
	database      *bun.DB
	ctx           context.Context
	generalConfig *configuration.GeneralConfiguration
}

func NewNftService(database *bun.DB, ctx context.Context, generalConfig *configuration.GeneralConfiguration) INftService {
	return &NftService{database: database, ctx: ctx, generalConfig: generalConfig}
}

func (ns *NftService) GetLatestTransfers(page int, perPage int) (*nftModel.NftTransactions, error) {
	var offSet = perPage * (page - 1)

	dbNfts := make([]DB.NftTransfer, 0)
	err := ns.database.NewSelect().Model((*DB.NftTransfer)(nil)).
		ColumnExpr("nft_transfer.*").
		ColumnExpr("t.timestamp AS transaction__timestamp, t.input_data AS transaction__input_data").
		ColumnExpr("tt.name AS token_type__name").
		Join("JOIN token_types AS tt ON tt.id = nft_transfer.token_type_id").
		Join("JOIN transactions AS t ON t.hash = nft_transfer.transaction_hash").
		Order("block_number DESC").
		Order("t.transaction_index DESC").
		Order("nft_transfer.index DESC").
		Order("nft_transfer.token_id DESC").
		Limit(perPage).Offset(offSet).Scan(ns.ctx, &dbNfts)

	if err != nil {
		//TODO: error handling
	}

	var result nftModel.NftTransactions

	isAddressContract := make(map[string]bool)
	for _, dbNft := range dbNfts {

		transaction := nftModel.NftTransaction{
			Type:            dbNft.TokenType.Name,
			From:            dbNft.From,
			To:              dbNft.To,
			NftId:           dbNft.TokenId,
			Hash:            dbNft.TransactionHash,
			Method:          getMethodName(dbNft.Transaction.InputData),
			Age:             utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(dbNft.Transaction.Timestamp), 0)).Seconds()))),
			DateTime:        time.Unix(int64(dbNft.Transaction.Timestamp), 0).UTC().Format("2006-01-02 15:04:05"),
			ContractAddress: dbNft.Address,
		}
		number := new(big.Int)
		if value, ok := number.SetString(dbNft.Value, 10); ok == true {
			if value.Cmp(big.NewInt(99)) == 1 {
				transaction.Value = ">99"
			} else if value.Cmp(big.NewInt(1)) == 1 && value.Cmp(big.NewInt(100)) == -1 {
				transaction.Value = fmt.Sprintf("x%v", dbNft.Value)
			}
		}

		isFromContract, exists := isAddressContract[transaction.From]
		if !exists {
			isFromContract, _ = ns.database.NewSelect().Table("contracts").Where("address = ?", transaction.From).Exists(ns.ctx)
			isAddressContract[transaction.From] = isFromContract
		}

		isToContract, exists := isAddressContract[transaction.To]
		if !exists {
			isToContract, _ = ns.database.NewSelect().Table("contracts").Where("address = ?", transaction.To).Exists(ns.ctx)
			isAddressContract[transaction.To] = isToContract
		}

		transaction.IsFromContract = isFromContract
		transaction.IsToContract = isToContract

		result.Transactions = append(result.Transactions, transaction)
	}

	var totalRows int
	ns.database.NewRaw("SELECT reltuples::bigint FROM pg_class WHERE oid = 'public.nft_transfers' ::regclass;").Scan(ns.ctx, &totalRows)
	result.TotalRows = int64(totalRows)

	maxRows := int(ns.generalConfig.NftLatestTransfersMaxCount)
	if totalRows > maxRows {
		totalRows = maxRows
	}

	totalPages := math.Ceil(float64(totalRows) / float64(perPage))
	if totalPages == 0 {
		result.TotalPages = 1
	} else {
		result.TotalPages = int(totalPages)
	}
	result.MaxCount = int(maxRows)
	return &result, nil
}

func (ns *NftService) GetNftMetadata(address string, tokenId string) (*models.NftMetadataModel, error) {
	var nftMetadata DB.NftMetadata
	error1 := ns.database.NewSelect().Table("nft_metadata").Where("token_id = ?", tokenId).Scan(ns.ctx, &nftMetadata)

	if error1 != nil {

	}

	if strings.HasPrefix(nftMetadata.Image, "ipfs") {
		nftMetadata.Image = "http://nswd.ddns.net:8088/ipfs/" + strings.Split(nftMetadata.Image, "//")[1]
	}

	var nftMetadataModel = models.NftMetadataModel{
		Id:      nftMetadata.Id,
		TokenId: nftMetadata.TokenId,
		Address: nftMetadata.Address,
		Name:    nftMetadata.Name,
		Image:   nftMetadata.Image,
	}

	var transfersDB []DB.NftTransfer
	err := ns.database.NewSelect().Model((*DB.NftTransfer)(nil)).
		ColumnExpr("nft_transfer.*").
		ColumnExpr("t.timestamp AS transaction__timestamp, t.input_data AS transaction__input_data").
		ColumnExpr("tt.name AS token_type__name").
		Join("JOIN token_types AS tt ON tt.id = nft_transfer.token_type_id").
		Join("JOIN transactions AS t ON t.hash = nft_transfer.transaction_hash").
		Order("block_number DESC").
		Order("t.transaction_index DESC").
		Order("nft_transfer.index DESC").
		Order("nft_transfer.token_id DESC").
		Where("? = ? AND ? = ?", bun.Ident("address"), strings.ToLower(address), bun.Ident("token_id"), tokenId).
		Limit(25).Offset(0).Scan(ns.ctx, &transfersDB)

	if err != nil {
		//TODO: Error handling
	}

	var transfers []models.TransferModel

	for _, v := range transfersDB {
		var oneTransfer = models.TransferModel{
			From:            v.From,
			To:              v.To,
			Age:             utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Transaction.Timestamp), 0)).Seconds()))),
			DateTime:        time.Unix(int64(v.Transaction.Timestamp), 0).UTC().Format("Jan-02-2006 15:04:05"),
			TransactionHash: v.TransactionHash,
		}

		transfers = append(transfers, oneTransfer)
	}

	nftMetadataModel.Transfers = transfers

	var totalRows int
	err = ns.database.QueryRow("SELECT COUNT(*) FROM nft_transfers WHERE address = ? AND token_id = ?", strings.ToLower(address), tokenId).Scan(&totalRows)
	nftMetadataModel.TotalRows = int64(totalRows)

	return &nftMetadataModel, nil
}
