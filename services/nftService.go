package services

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"
	"webbc/DB"
	"webbc/configuration"
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

	dbNfts := make([]DB.Nft, 0)
	err := ns.database.NewSelect().Model((*DB.Nft)(nil)).
		ColumnExpr("nft.*").
		ColumnExpr("t.timestamp AS transaction__timestamp, t.input_data AS transaction__input_data").
		ColumnExpr("tt.name AS token_type__name").
		Join("JOIN token_types AS tt ON tt.id = nft.token_type_id").
		Join("JOIN transactions AS t ON t.hash = nft.transaction_hash").
		Order("block_number DESC").
		Order("t.transaction_index DESC").
		Order("nft.index DESC").
		Order("nft.token_id DESC").
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
	ns.database.NewRaw("SELECT reltuples::bigint FROM pg_class WHERE oid = 'public.nfts' ::regclass;").Scan(ns.ctx, &totalRows)
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
