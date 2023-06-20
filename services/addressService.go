package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
	"webbc/DB"
	"webbc/configuration"
	"webbc/eth"
	"webbc/models/abiModel"
	"webbc/models/addressModel"
	"webbc/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type AddressService struct {
	database *bun.DB
	ctx      context.Context
	config   *configuration.GeneralConfiguration
}

func NewAddressService(database *bun.DB, ctx context.Context, config *configuration.GeneralConfiguration) IAddressService {
	return &AddressService{database: database, ctx: ctx, config: config}
}

func (as *AddressService) GetAddress(address string) (*addressModel.Address, error) {
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

		txnFee := (float64(v.GasPrice) / float64(params.Ether)) * float64(v.GasUsed)
		transac.TxnFee = utils.ShowDecimalPrecision(txnFee, 8)

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
	isContract, err = as.database.NewSelect().Table("contracts").Where("address = ?", address).Exists(as.ctx)

	if isContract {
		var creatorTransaction string
		err = as.database.NewSelect().Table("contracts").ColumnExpr("transaction_hash").Where("address = ?", address).Scan(as.ctx, &creatorTransaction)
		if err != nil {

		}

		var creatorAddress string
		err = as.database.NewSelect().Table("transactions").ColumnExpr(`"from"`).Where("hash = ?", creatorTransaction).Scan(as.ctx, &creatorAddress)

		result.CreatorTransaction = creatorTransaction
		result.CreatorAddress = creatorAddress
	}

	result.IsContract = isContract
	return &result, nil
}

func (as *AddressService) getBalanceFromChainWithTimeout(address string) (string, error) {
	var result string
	ctxWithTimeout, cancel := context.WithTimeout(as.ctx, time.Duration(as.config.CallTimeoutInSeconds)*time.Second)
	defer cancel()
	err := eth.GetHttpNodeClient().CallContext(ctxWithTimeout, &result, "eth_getBalance", address, "latest")
	return result, err
}

func (as *AddressService) UploadABI(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))

	var abiItems []abiModel.AbiItem
	if err := c.BindJSON(&abiItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbAbis []DB.Abi = []DB.Abi{}
	for _, el := range abiItems {
		if el.Type == "event" {
			var event abiModel.EventItem = abiModel.EventItem{
				Anonymous: el.Anonymous,
				Inputs:    el.Inputs,
				Name:      el.Name,
				Type:      el.Type,
			}
			jsonBytes, err := json.Marshal(event)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			parsedEvent, err := abi.JSON(strings.NewReader("[" + string(jsonBytes) + "]"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			dbAbi := DB.Abi{
				Hash:       parsedEvent.Events[el.Name].ID.String(),
				Address:    address,
				AbiTypeId:  2,
				Definition: string(jsonBytes),
			}
			dbAbis = append(dbAbis, dbAbi)
		} else if el.Type == "function" {
			var function abiModel.FunctionItem = abiModel.FunctionItem{
				Inputs:          el.Inputs,
				Outputs:         el.Outputs,
				Name:            el.Name,
				StateMutability: el.StateMutability,
				Type:            el.Type,
			}
			jsonBytes, err := json.Marshal(function)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			parsedFunction, err := abi.JSON(strings.NewReader("[" + string(jsonBytes) + "]"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			dbAbi := DB.Abi{
				Hash:       hex.EncodeToString(parsedFunction.Methods[el.Name].ID),
				Address:    address,
				AbiTypeId:  3,
				Definition: string(jsonBytes),
			}
			dbAbis = append(dbAbis, dbAbi)
		} else if el.Type == "constructor" {
			fmt.Println("constructor")
		}
	}

	_, abiError := as.database.NewInsert().Model(&dbAbis).Exec(as.ctx)
	if abiError != nil {
		//TODO: error handling
		c.JSON(http.StatusBadRequest, gin.H{"error": abiError.Error()})
		return
	}

	c.Status(http.StatusOK)
}
