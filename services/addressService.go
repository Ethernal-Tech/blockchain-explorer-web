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
	config   *configuration.GeneralConfiguration
}

func NewAddressService(database *bun.DB, ctx context.Context, client *rpc.Client, config *configuration.GeneralConfiguration) AddressService {
	return &AddressServiceImplementation{database: database, ctx: ctx, client: client, config: config}
}

func (as *AddressServiceImplementation) GetAddress(address string) (*addressModel.Address, error) {
	var result addressModel.Address
	result.AddressHex = address
	address = strings.ToLower(address)

	var transactions []DB.Transaction
	err := as.database.NewSelect().Table("transactions").Where("? = ? OR ? = ? OR ? = ?", bun.Ident("from"), address, bun.Ident("to"), address, bun.Ident("contract_address"), address).Order("timestamp DESC").Limit(25).Scan(as.ctx, &transactions)
	if err != nil {
		//TODO: Error handling
	}

	for _, v := range transactions {
		var transac = addressModel.Transaction{
			Hash:            v.Hash,
			Method:          getMethodName(v.InputData),
			BlockNumber:     v.BlockNumber,
			From:            v.From,
			To:              v.To,
			Value:           utils.WeiToEther(v.Value),
			Gas:             v.Gas,
			GasUsed:         v.GasUsed,
			GasPrice:        v.GasPrice,
			Age:             utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds()))),
			DateTime:        time.Unix(int64(v.Timestamp), 0).UTC().Format("2006-01-02 15:04:05"),
			ContractAddress: v.ContractAddress,
		}
		if strings.ReplaceAll(transac.To, " ", "") == "" {
			transac.To = ""
		}

		// value := utils.WeiToEther(v.Value)
		// transac.Value = value

		if address == v.From {
			transac.Direction = "out"
		} else {
			transac.Direction = "in"
		}
		result.Transactions = append(result.Transactions, transac)

	}
	count, err := as.database.NewSelect().Table("transactions").Where("? = ? OR ? = ? OR ? = ?", bun.Ident("from"), address, bun.Ident("to"), address, bun.Ident("contract_address"), address).Count(as.ctx)
	result.TransactionCount = count

	balance, err := as.getBalanceFromChainWithTimeout(address)
	if err != nil {
		//TODO: error handling
	}

	//balanceBigInt := utils.ToBigInt(balance)
	result.Balance = utils.WeiToEther(balance)
	var isContract bool
	isContract, err = as.database.NewSelect().Table("transactions").Where("contract_address = ?", address).Exists(as.ctx)
	result.IsContract = isContract
	return &result, nil
}

func (as *AddressServiceImplementation) getBalanceFromChainWithTimeout(address string) (string, error) {
	var result string
	ctxWithTimeout, cancel := context.WithTimeout(as.ctx, time.Duration(as.config.CallTimeoutInSeconds)*time.Second)
	defer cancel()
	err := as.client.CallContext(ctxWithTimeout, &result, "eth_getBalance", address, "latest")
	return result, err
}

func (as *AddressServiceImplementation) ChangeClient(client *rpc.Client) {
	as.client = client
}
