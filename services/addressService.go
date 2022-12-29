package services

import (
	"context"
	"math"
	"strings"
	"time"
	"webbc/DB"
	"webbc/configuration"
	"webbc/models/addressModel"
	"webbc/utils"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/uptrace/bun"
)

type AddressServiceImplementation struct {
	database *bun.DB
	ctx      context.Context
	client   *rpc.Client
	config   *configuration.Configuration
}

func NewAddressService(database *bun.DB, ctx context.Context, client *rpc.Client, config *configuration.Configuration) AddressService {
	return &AddressServiceImplementation{database: database, ctx: ctx, client: client, config: config}
}

func (as *AddressServiceImplementation) GetAddress(address string) (*addressModel.Address, error) {
	var result addressModel.Address
	result.AddressHex = address
	address = strings.ToLower(address)

	var transactions []DB.Transaction
	err := as.database.NewSelect().Table("transactions").Where("? = ? OR ? = ?", bun.Ident("from"), address, bun.Ident("to"), address).Order("timestamp DESC").Limit(25).Scan(as.ctx, &transactions)
	if err != nil {
		//TODO: Error handling
	}

	for _, v := range transactions {
		var transac = addressModel.Transaction{
			Hash:        v.Hash,
			BlockNumber: v.BlockNumber,
			From:        v.From,
			To:          v.To,
			Gas:         v.Gas,
			GasUsed:     v.GasUsed,
			GasPrice:    v.GasPrice,
			Timestamp:   int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds())),
		}
		value := utils.WeiToEther(utils.ToBigInt(v.Value)).String()
		transac.Value = value[:10]
		if address == v.To {
			transac.Direction = "in"
		} else {
			transac.Direction = "out"
		}
		result.Transactions = append(result.Transactions, transac)

	}
	count, err := as.database.NewSelect().Table("transactions").Where("? = ? OR ? = ?", bun.Ident("from"), address, bun.Ident("to"), address).Count(as.ctx)
	result.TransactionCount = count

	balance, err := as.getBalanceFromChainWithTimeout(address)
	if err != nil {
		//TODO: error handling
	}

	balanceBigInt := utils.ToBigInt(balance)
	result.Balance = utils.WeiToEther(balanceBigInt)
	return &result, nil
}

func (as *AddressServiceImplementation) getBalanceFromChainWithTimeout(address string) (string, error) {
	var result string
	ctxWithTimeout, cancel := context.WithTimeout(as.ctx, time.Duration(as.config.CallTimeoutInSeconds)*time.Second)
	defer cancel()
	err := as.client.CallContext(ctxWithTimeout, &result, "eth_getBalance", address, "latest")
	return result, err
}
