package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
	"webbc/DB"
	"webbc/configuration"
	"webbc/models"
	"webbc/models/transactionModel"
	"webbc/utils"

	ethereumAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/uptrace/bun"
)

type TransactionServiceImplementation struct {
	database      *bun.DB
	ctx           context.Context
	generalConfig *configuration.GeneralConfiguration
}

func NewTransactionService(database *bun.DB, ctx context.Context, generalConfig *configuration.GeneralConfiguration) TransactionService {
	return &TransactionServiceImplementation{database: database, ctx: ctx, generalConfig: generalConfig}
}

func (tsi *TransactionServiceImplementation) GetLastTransactions(numberOfTransactions int) (*[]models.Transaction, error) {
	var transactions []DB.Transaction
	err := tsi.database.NewSelect().Table("transactions").Order("block_number DESC").Limit(20).Scan(tsi.ctx, &transactions)

	if err != nil {

	}

	var result []models.Transaction

	for _, v := range transactions {
		var oneResultTransaction = models.Transaction{
			Hash:             v.Hash,
			BlockHash:        v.BlockHash,
			BlockNumber:      v.BlockNumber,
			From:             v.From,
			To:               v.To,
			Gas:              v.Gas,
			GasUsed:          v.GasUsed,
			GasPrice:         v.GasPrice,
			Nonce:            v.Nonce,
			TransactionIndex: v.TransactionIndex,
			ContractAddress:  v.ContractAddress,
			Status:           v.Status,
			Age:              utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds()))),
			DateTime:         time.Unix(int64(v.Timestamp), 0).UTC().Format("Jan-02-2006 15:04:05"),
		}

		result = append(result, oneResultTransaction)
	}

	return &result, nil
}

func (tsi *TransactionServiceImplementation) GetTransactionByHash(transactionHash string) (*models.Transaction, error) {
	var transaction DB.Transaction
	var logs []DB.Log
	error1 := tsi.database.NewSelect().Table("transactions").Where("hash = ?", transactionHash).Scan(tsi.ctx, &transaction)

	if error1 != nil {

	}

	error2 := tsi.database.NewSelect().Table("logs").Order("index ASC").Where("transaction_hash = ?", transactionHash).Scan(tsi.ctx, &logs)

	if error2 != nil {

	}

	var oneResultTransaction = models.Transaction{
		Hash:             transaction.Hash,
		BlockHash:        transaction.BlockHash,
		BlockNumber:      transaction.BlockNumber,
		From:             transaction.From,
		To:               transaction.To,
		Gas:              transaction.Gas,
		GasUsed:          transaction.GasUsed,
		GasPrice:         transaction.GasPrice,
		Nonce:            transaction.Nonce,
		TransactionIndex: transaction.TransactionIndex,
		Value:            utils.WeiToEther(transaction.Value),
		ContractAddress:  transaction.ContractAddress,
		Status:           transaction.Status,
		Age:              utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(transaction.Timestamp), 0)).Seconds()))),
		DateTime:         time.Unix(int64(transaction.Timestamp), 0).UTC().Format("Jan-02-2006 15:04:05"),
		InputData:        transaction.InputData,
	}

	isToContract, _ := tsi.database.NewSelect().Table("contracts").Where("address = ?", oneResultTransaction.To).Exists(tsi.ctx)

	//default and decoded view of input data
	if isToContract && oneResultTransaction.InputData != "0x" {
		var abi DB.Abi
		error1 = tsi.database.NewSelect().Table("abis").Where("hash = ? AND address = ?", oneResultTransaction.InputData[2:10], oneResultTransaction.To).Scan(tsi.ctx, &abi)
		if abi.Id != 0 {
			oneResultTransaction.IsUploadedABI = true
			oneResultTransaction.InputDataSig, oneResultTransaction.InputDataMethodId, oneResultTransaction.InputDataParamValues, oneResultTransaction.DecodedInputData = defaultAndDecodedViewInputData(oneResultTransaction.InputData, &abi)
		}
	}

	//logs
	for _, log := range logs {

		var abi DB.Abi
		error1 = tsi.database.NewSelect().Table("abis").Where("hash = ? AND address = ?", log.Topic0, log.Address).Scan(tsi.ctx, &abi)
		if error1 != nil {

		}
		var log = createLogModel(&log, &abi)

		oneResultTransaction.Logs = append(oneResultTransaction.Logs, log)
	}

	oneResultTransaction.IsToContract = isToContract

	return &oneResultTransaction, nil
}

func defaultAndDecodedViewInputData(inputData string, dbAbi *DB.Abi) (string, string, []interface{}, models.DecodedInputData) {

	parsedAbi, err := ethereumAbi.JSON(strings.NewReader("[" + dbAbi.Definition + "]"))
	if err != nil {
		return "", "", []interface{}{}, models.DecodedInputData{}
	}

	var signature string
	var methodId string
	var paramValues []interface{}
	var decodedInputData models.DecodedInputData

	for _, method := range parsedAbi.Methods {

		signature = method.Name + "("
		params := inputData[10:]

		substrLen := 64
		numSubstr := len(inputData) / substrLen

		for i := 0; i < numSubstr; i++ {
			paramValues = append(paramValues, params[i*substrLen:(i+1)*substrLen])
		}

		var hasTuple bool

		for i, input := range method.Inputs {
			if len(input.Type.TupleElems) == 0 {
				if i < len(method.Inputs)-1 {
					signature = signature + input.Type.String() + " " + input.Name + ", "
				} else {
					signature = signature + input.Type.String() + " " + input.Name + ")"
				}

			} else {
				hasTuple = true
				if i < len(method.Inputs)-1 {
					signature = signature + "tuple" + " " + input.Name + ", "
				} else {
					signature = signature + "tuple" + " " + input.Name + ")"
				}
			}
		}

		decodedInputData.FunctionSignature = signature

		if !hasTuple {
			unpackValues, _ := method.Inputs.UnpackValues(common.Hex2Bytes(params))

			dataNames, dataTypes, dataValues := decodeData(unpackValues, method.Inputs)
			for i, _ := range dataValues {
				var parameter = models.ParameterInfo{
					Name:  dataNames[i],
					Type:  dataTypes[i],
					Value: dataValues[i],
				}

				decodedInputData.Parameters = append(decodedInputData.Parameters, parameter)
			}

		} else {
			//TODO

		}

		methodId = hex.EncodeToString(method.ID)
	}

	return signature, methodId, paramValues, decodedInputData
}

func createLogModel(dbLog *DB.Log, dbAbi *DB.Abi) models.Log {

	var log models.Log
	if dbAbi.Id != 0 {
		parsedAbi, err := ethereumAbi.JSON(strings.NewReader("[" + dbAbi.Definition + "]"))
		if err != nil {
			return models.Log{}
		}

		for _, event := range parsedAbi.Events {
			//var inputs []string

			paramNames := make([]string, len(event.Inputs))
			paramTypes := make([]string, len(event.Inputs))
			paramIndexed := make([]bool, len(event.Inputs))
			for index, input := range event.Inputs {
				paramNames[index] = input.Name
				paramTypes[index] = input.Type.String()
				paramIndexed[index] = input.Indexed
			}

			unpackValues, err := event.Inputs.NonIndexed().UnpackValues(common.Hex2Bytes(dbLog.Data[2:]))
			if err != nil {
				// Handle the error
				return models.Log{}
			}

			dataNames, _, dataValues := decodeData(unpackValues, event.Inputs.NonIndexed())

			log = models.Log{
				BlockHash:       dbLog.BlockHash,
				Index:           dbLog.Index,
				TransactionHash: dbLog.TransactionHash,
				Address:         dbLog.Address,
				BlockNumber:     dbLog.BlockNumber,
				EventName:       event.Name,
				Topic0:          dbLog.Topic0,
				Topic1:          dbLog.Topic1,
				Topic2:          dbLog.Topic2,
				Topic3:          dbLog.Topic3,
				ParamNames:      paramNames,
				ParamTypes:      paramTypes,
				ParamIndexed:    paramIndexed,
				DataNames:       dataNames,
				DataValues:      dataValues,
				Data:            dbLog.Data,
			}

			return log
		}
	}

	log = models.Log{
		BlockHash:       dbLog.BlockHash,
		Index:           dbLog.Index,
		TransactionHash: dbLog.TransactionHash,
		Address:         dbLog.Address,
		BlockNumber:     dbLog.BlockNumber,
		EventName:       "",
		Topic0:          dbLog.Topic0,
		Topic1:          dbLog.Topic1,
		Topic2:          dbLog.Topic2,
		Topic3:          dbLog.Topic3,
		Data:            dbLog.Data,
	}

	return log
}

func decodeData(unpackValues []interface{}, arguments ethereumAbi.Arguments) ([]string, []string, []string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()

	//nonIndexed := event.Inputs.NonIndexed()
	dataNames := make([]string, len(arguments))
	dataTypes := make([]string, len(arguments))
	dataValues := make([]string, len(arguments))

	for index, value := range unpackValues {
		t := reflect.TypeOf(value)
		v := reflect.ValueOf(value)
		dataNames[index] = arguments[index].Name
		dataTypes[index] = arguments[index].Type.String()

		switch t.Kind() {
		case reflect.Slice:
			if v.Type().Elem().Kind() == reflect.Ptr {
				var intSlice []string
				for i := 0; i < v.Len(); i++ {
					intSlice = append(intSlice, fmt.Sprintf("%v", v.Index(i)))
				}
				dataValues[index] = fmt.Sprintf("%s", strings.Join(intSlice, " "))
			} else if v.Type().Elem().Kind() == reflect.Uint8 {
				dataValues[index] = hex.EncodeToString(v.Bytes())
			}
		case reflect.Array:
			if t.String() == "common.Address" {
				dataValues[index] = fmt.Sprintf("%v", v)
			} else {
				byteSlice := make([]byte, v.Len())
				for i := 0; i < v.Len(); i++ {
					byteSlice[i] = byte(v.Index(i).Uint())
				}
				dataValues[index] = hex.EncodeToString(byteSlice)
			}
		default:
			dataValues[index] = fmt.Sprintf("%v", v)
		}
	}

	return dataNames, dataTypes, dataValues
}

func (tsi *TransactionServiceImplementation) GetTransactionsInBlock(blockNumber string, page int, perPage int) (*transactionModel.Transactions, error) {
	var transactions []DB.Transaction

	var offSet = perPage * (page - 1)
	err := tsi.database.NewSelect().Table("transactions").OrderExpr("transaction_index DESC").Where("block_number = ?", blockNumber).Limit(perPage).Offset(offSet).Scan(tsi.ctx, &transactions)

	if err != nil {
		//TODO: error handling
	}

	var result transactionModel.Transactions

	for _, v := range transactions {
		var transaction = transactionModel.Transaction{
			Hash:            v.Hash,
			Method:          getMethodName(v.InputData),
			BlockNumber:     v.BlockNumber,
			Age:             utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds()))),
			DateTime:        time.Unix(int64(v.Timestamp), 0).UTC().Format("2006-01-02 15:04:05"),
			From:            v.From,
			To:              v.To,
			Value:           utils.WeiToEther(v.Value),
			TxnFee:          0000000,
			ContractAddress: v.ContractAddress,
		}

		isToContract, _ := tsi.database.NewSelect().Table("contracts").Where("address = ?", transaction.To).Exists(tsi.ctx)
		transaction.IsToContract = isToContract

		result.Transactions = append(result.Transactions, transaction)
	}

	var totalRows int
	totalRows, err = tsi.database.NewSelect().Table("transactions").Where("? = ?", bun.Ident("block_number"), blockNumber).Count(tsi.ctx)
	result.TotalRows = int64(totalRows)

	totalPages := math.Ceil(float64(totalRows) / float64(perPage))
	if totalPages == 0 {
		result.TotalPages = 1
	} else {
		result.TotalPages = int(totalPages)
	}

	return &result, nil
}

func (tsi TransactionServiceImplementation) GetAllTransactions(page int, perPage int) (*transactionModel.Transactions, error) {
	var transactions []DB.Transaction

	var offSet = perPage * (page - 1)
	err := tsi.database.NewSelect().Table("transactions").OrderExpr("block_number DESC").Limit(perPage).Offset(offSet).Scan(tsi.ctx, &transactions)
	if err != nil {
		//TODO: error handling
	}

	var result transactionModel.Transactions
	for _, v := range transactions {
		var transaction = transactionModel.Transaction{
			Hash:            v.Hash,
			Method:          getMethodName(v.InputData),
			BlockNumber:     v.BlockNumber,
			Age:             utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds()))),
			DateTime:        time.Unix(int64(v.Timestamp), 0).UTC().Format("2006-01-02 15:04:05"),
			From:            v.From,
			To:              v.To,
			Value:           utils.WeiToEther(v.Value),
			TxnFee:          0000000,
			ContractAddress: v.ContractAddress,
		}

		isToContract, _ := tsi.database.NewSelect().Table("contracts").Where("address = ?", transaction.To).Exists(tsi.ctx)
		transaction.IsToContract = isToContract
		// if strings.ReplaceAll(transaction.To, " ", "") == "" {
		// 	transaction.To = ""
		// }
		result.Transactions = append(result.Transactions, transaction)
	}

	var totalRows int64
	tsi.database.NewRaw("SELECT reltuples::bigint FROM pg_class WHERE oid = 'public.transactions' ::regclass;").Scan(tsi.ctx, &totalRows)
	result.TotalRows = int64(totalRows)

	maxRows := int64(tsi.generalConfig.TransactionsMaxCount)
	if totalRows > maxRows {
		totalRows = maxRows
	}

	totalPages := math.Ceil(float64(totalRows) / float64(perPage))
	if totalPages == 0 {
		result.TotalPages = 1
	} else {
		result.TotalPages = int(totalPages)
	}
	return &result, nil
}

func (tsi TransactionServiceImplementation) GetTransactionsByAddress(address string, page int, perPage int) (*transactionModel.Transactions, error) {
	var transactions []DB.Transaction

	var offest = perPage * (page - 1)
	address = strings.ToLower(address)
	err := tsi.database.NewSelect().Table("transactions").OrderExpr("block_number DESC").Where("? = ? OR ? = ?", bun.Ident("from"), address, bun.Ident("to"), address).Offset(offest).Limit(perPage).Scan(tsi.ctx, &transactions)
	if err != nil {
		//TODO: error handling
	}

	var result transactionModel.Transactions
	for _, v := range transactions {
		var transaction = transactionModel.Transaction{
			Hash:            v.Hash,
			Method:          getMethodName(v.InputData),
			BlockNumber:     v.BlockNumber,
			Age:             utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds()))),
			DateTime:        time.Unix(int64(v.Timestamp), 0).UTC().Format("2006-01-02 15:04:05"),
			From:            v.From,
			To:              v.To,
			Value:           utils.WeiToEther(v.Value),
			TxnFee:          0000000,
			ContractAddress: v.ContractAddress,
		}
		// if strings.ReplaceAll(transaction.To, " ", "") == "" {
		// 	transaction.To = ""
		// }
		if address == v.To {
			transaction.Direction = "in"
		} else {
			transaction.Direction = "out"
		}
		result.Transactions = append(result.Transactions, transaction)
	}

	var totalRows int
	totalRows, err = tsi.database.NewSelect().Table("transactions").Where("? = ? OR ? = ?", bun.Ident("from"), address, bun.Ident("to"), address).Count(tsi.ctx)
	result.TotalRows = int64(totalRows)

	maxRows := tsi.generalConfig.TransactionsByAddressMaxCount
	if totalRows > int(maxRows) {
		totalRows = int(maxRows)
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

func getMethodName(str string) string {
	if len(str) >= 10 {
		return str[:10]
	} else {
		return str
	}
}
