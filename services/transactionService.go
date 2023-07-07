package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"webbc/DB"
	webcommon "webbc/common"
	"webbc/configuration"
	"webbc/eth"
	"webbc/models/transactionModel"
	"webbc/utils"

	ethereumAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/uptrace/bun"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type TransactionServiceImplementation struct {
	database      *bun.DB
	ctx           context.Context
	generalConfig *configuration.GeneralConfiguration
}

func NewTransactionService(database *bun.DB, ctx context.Context, generalConfig *configuration.GeneralConfiguration) ITransactionService {
	return &TransactionServiceImplementation{database: database, ctx: ctx, generalConfig: generalConfig}
}

func (tsi *TransactionServiceImplementation) GetLastTransactions(numberOfTransactions int) (*[]transactionModel.Transaction, error) {
	var transactions []DB.Transaction
	err := tsi.database.NewSelect().Table("transactions").Order("block_number DESC").Limit(20).Scan(tsi.ctx, &transactions)

	if err != nil {

	}

	var result []transactionModel.Transaction

	for _, v := range transactions {
		var oneResultTransaction = transactionModel.Transaction{
			Hash: v.Hash,
			From: v.From,
			To:   v.To,
			Age:  utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds()))),
		}

		result = append(result, oneResultTransaction)
	}

	return &result, nil
}

func (tsi *TransactionServiceImplementation) GetTransactionByHash(transactionHash string) (*transactionModel.Transaction, error) {
	var transaction DB.Transaction
	var logs []DB.Log
	error1 := tsi.database.NewSelect().Table("transactions").Where("hash = ?", transactionHash).Scan(tsi.ctx, &transaction)

	if error1 != nil {

	}

	error2 := tsi.database.NewSelect().Table("logs").Order("index ASC").Where("transaction_hash = ?", transactionHash).Scan(tsi.ctx, &logs)

	if error2 != nil {

	}

	p := message.NewPrinter(language.English)
	var oneResultTransaction = transactionModel.Transaction{
		Hash:              transaction.Hash,
		BlockHash:         transaction.BlockHash,
		BlockNumber:       transaction.BlockNumber,
		From:              transaction.From,
		To:                transaction.To,
		Gas:               p.Sprintf("%v", transaction.Gas),
		GasUsed:           p.Sprintf("%v", transaction.GasUsed),
		GasUsedPercentage: math.Round(((float64(transaction.GasUsed)/float64(transaction.Gas))*100)*100) / 100,
		GasPriceInGwei:    strconv.FormatFloat(float64(transaction.GasPrice)/float64(params.GWei), 'f', -1, 64),
		Nonce:             transaction.Nonce,
		TransactionIndex:  transaction.TransactionIndex,
		Value:             utils.WeiToEther(transaction.Value),
		ContractAddress:   transaction.ContractAddress,
		Status:            transaction.Status,
		Age:               utils.Convert(int(math.Round(time.Now().Sub(time.Unix(int64(transaction.Timestamp), 0)).Seconds()))),
		DateTime:          time.Unix(int64(transaction.Timestamp), 0).UTC().Format("Jan-02-2006 15:04:05"),
		InputData:         transaction.InputData,
	}

	txnFee := (float64(transaction.GasPrice) / float64(params.Ether)) * float64(transaction.GasUsed)
	oneResultTransaction.TxnFee = utils.ShowDecimalPrecision(txnFee, 18)

	gasPriceInEth := float64(transaction.GasPrice) / float64(params.Ether)
	oneResultTransaction.GasPriceInEth = utils.ShowDecimalPrecision(gasPriceInEth, 18)

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

	contractName := make(map[string]string) //Optimization map
	//logs
	for _, log := range logs {

		tsi.transferType(log, &oneResultTransaction, &contractName)

		var abi DB.Abi
		error1 = tsi.database.NewSelect().Table("abis").Where("hash = ? AND address = ?", log.Topic0, log.Address).Scan(tsi.ctx, &abi)
		if error1 != nil {

		}

		var logModel = createLogModel(&log, &abi)
		oneResultTransaction.Logs = append(oneResultTransaction.Logs, logModel)
	}

	oneResultTransaction.IsToContract = isToContract

	return &oneResultTransaction, nil
}

func defaultAndDecodedViewInputData(inputData string, dbAbi *DB.Abi) (string, string, []interface{}, transactionModel.DecodedInputData) {

	parsedAbi, err := ethereumAbi.JSON(strings.NewReader("[" + dbAbi.Definition + "]"))
	if err != nil {
		return "", "", []interface{}{}, transactionModel.DecodedInputData{}
	}

	var signature string
	var methodId string
	var paramValues []interface{}
	var decodedInputData transactionModel.DecodedInputData

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
				var parameter = transactionModel.ParameterInfo{
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

func (tsi *TransactionServiceImplementation) transferType(log DB.Log, oneResultTransaction *transactionModel.Transaction, contractName *map[string]string) {
	if log.Topic0 == webcommon.Erc20TransferEvent.Signature && log.Topic1 != "" && log.Topic2 != "" && log.Topic3 == "" {
		transferModel := transactionModel.TransferModel{
			From: "0x" + log.Topic1[len(log.Topic1)-40:],
			To:   "0x" + log.Topic2[len(log.Topic2)-40:],
		}
		name, decimals := tsi.getNameAndDecimals(log.Address)
		parsedAbi, err := ethereumAbi.JSON(strings.NewReader("[" + webcommon.Erc20TransferEvent.Abi + "]"))
		if err == nil {
			for _, event := range parsedAbi.Events {
				unpackValues, err := event.Inputs.NonIndexed().UnpackValues(common.Hex2Bytes(log.Data[2:]))
				if err != nil {
					break
				}
				_, _, dataValues := decodeData(unpackValues, event.Inputs.NonIndexed())
				amount, ok := new(big.Int).SetString(dataValues[0], 10)
				if ok && decimals != 0 {
					exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
					tokenAmount := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(exp))
					resultStr := tokenAmount.Text('f', 18)
					resultStr = strings.TrimRight(resultStr, "0")
					resultStr = strings.TrimRight(resultStr, ".")
					transferModel.Value = resultStr
				} else {
					transferModel.Value = dataValues[0]
				}
			}
		}
		transferModel.TokenAddress = log.Address
		transferModel.TokenName = name
		oneResultTransaction.ERC20Transfers = append(oneResultTransaction.ERC20Transfers, transferModel)
	} else if log.Topic0 == webcommon.Erc721TransferEvent.Signature && log.Topic1 != "" && log.Topic2 != "" && log.Topic3 != "" {
		transferModel := transactionModel.TransferModel{
			From:    "0x" + log.Topic1[len(log.Topic1)-40:],
			To:      "0x" + log.Topic2[len(log.Topic2)-40:],
			TokenId: utils.HexNumberToString(log.Topic3),
		}
		name, ok := (*contractName)[log.Address]
		if !ok {
			name = tsi.GetName(log.Address)
			(*contractName)[log.Address] = name
		}
		transferModel.TokenAddress = log.Address
		transferModel.TokenName = name
		oneResultTransaction.ERC721Transfers = append(oneResultTransaction.ERC721Transfers, transferModel)
	} else if log.Topic0 == webcommon.Erc1155TransferSingleEvent.Signature && log.Topic2 != "" && log.Topic3 != "" {
		transferModel := transactionModel.TransferModel{
			From: "0x" + log.Topic2[len(log.Topic2)-40:],
			To:   "0x" + log.Topic3[len(log.Topic3)-40:],
		}
		parsedAbi, err := ethereumAbi.JSON(strings.NewReader("[" + webcommon.Erc1155TransferSingleEvent.Abi + "]"))
		if err == nil {
			for _, event := range parsedAbi.Events {
				unpackValues, err := event.Inputs.NonIndexed().UnpackValues(common.Hex2Bytes(log.Data[2:]))
				if err != nil {
					break
				}
				_, _, dataValues := decodeData(unpackValues, event.Inputs.NonIndexed())
				transferModel.TokenId = dataValues[0]
				transferModel.Value = dataValues[1]
			}
		}

		name, ok := (*contractName)[log.Address]
		if !ok {
			name = tsi.GetName(log.Address)
			(*contractName)[log.Address] = name
		}
		transferModel.TokenAddress = log.Address
		transferModel.TokenName = name
		oneResultTransaction.ERC1155Transfers = append(oneResultTransaction.ERC1155Transfers, transferModel)
	} else if log.Topic0 == webcommon.Erc1155TransferBatchEvent.Signature && log.Topic2 != "" && log.Topic3 != "" {
		name, ok := (*contractName)[log.Address]
		if !ok {
			name = tsi.GetName(log.Address)
			(*contractName)[log.Address] = name
		}
		parsedAbi, err := ethereumAbi.JSON(strings.NewReader("[" + webcommon.Erc1155TransferBatchEvent.Abi + "]"))
		var dataValues []string
		if err == nil {
			for _, event := range parsedAbi.Events {
				unpackValues, err := event.Inputs.NonIndexed().UnpackValues(common.Hex2Bytes(log.Data[2:]))
				if err != nil {
					break
				}
				_, _, dataValues = decodeData(unpackValues, event.Inputs.NonIndexed())
			}
		}
		tokenIds := strings.Split(dataValues[0], " ")
		values := strings.Split(dataValues[1], " ")
		for index, id := range tokenIds {
			transferModel := transactionModel.TransferModel{
				From:         "0x" + log.Topic2[len(log.Topic2)-40:],
				To:           "0x" + log.Topic3[len(log.Topic3)-40:],
				TokenId:      id,
				Value:        values[index],
				TokenAddress: log.Address,
				TokenName:    name,
			}
			oneResultTransaction.ERC1155Transfers = append(oneResultTransaction.ERC1155Transfers, transferModel)
		}
	}
}

func (tsi *TransactionServiceImplementation) getNameAndDecimals(address string) (string, uint64) {
	var nameStr string
	var decimalsStr string
	ctxWithTimeout, cancel := context.WithTimeout(tsi.ctx, time.Duration(tsi.generalConfig.CallTimeoutInSeconds)*time.Second)
	defer cancel()
	type request struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}

	var elems []rpc.BatchElem

	elems = append(elems, rpc.BatchElem{
		Method: "eth_call",
		Args:   []interface{}{request{address, webcommon.NameMethodSignature}, "latest"},
		Result: &nameStr,
	})

	elems = append(elems, rpc.BatchElem{
		Method: "eth_call",
		Args:   []interface{}{request{address, webcommon.DecimalsMethodSignature}, "latest"},
		Result: &decimalsStr,
	})
	err := eth.GetHttpNodeClient().BatchCallContext(ctxWithTimeout, elems)
	if err != nil {
		nameStr = ""
		//TODO: Error handling
	}
	if len(nameStr) > 0 && nameStr[0:2] == "0x" {
		nameStr = nameStr[2:]
	}
	bs, err := hex.DecodeString(nameStr)
	re := regexp.MustCompile("[^a-zA-Z0-9// -.]+")
	str := re.ReplaceAllString(string(bs), "")

	var decimal uint64 = 0
	if len(decimalsStr) > 0 {
		decimal = utils.ToUint64(decimalsStr)
	}
	return str, decimal
}

func (tsi *TransactionServiceImplementation) GetName(address string) string {
	var nameStr string
	ctxWithTimeout, cancel := context.WithTimeout(tsi.ctx, time.Duration(tsi.generalConfig.CallTimeoutInSeconds)*time.Second)
	defer cancel()
	type request struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}
	var elems []rpc.BatchElem

	elems = append(elems, rpc.BatchElem{
		Method: "eth_call",
		Args:   []interface{}{request{address, webcommon.NameMethodSignature}, "latest"},
		Result: &nameStr,
	})
	err := eth.GetHttpNodeClient().BatchCallContext(ctxWithTimeout, elems)
	if err != nil {
		nameStr = ""
		//TODO: Error handling
	}
	if len(nameStr) > 0 && nameStr[0:2] == "0x" {
		nameStr = nameStr[2:]
	}
	bs, err := hex.DecodeString(nameStr)
	re := regexp.MustCompile("[^a-zA-Z0-9// -.]+")
	str := re.ReplaceAllString(string(bs), "")
	return str
}

func createLogModel(dbLog *DB.Log, dbAbi *DB.Abi) transactionModel.Log {

	var log transactionModel.Log
	if dbAbi.Id != 0 {
		parsedAbi, err := ethereumAbi.JSON(strings.NewReader("[" + dbAbi.Definition + "]"))
		if err != nil {
			return transactionModel.Log{}
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
				return transactionModel.Log{}
			}

			dataNames, _, dataValues := decodeData(unpackValues, event.Inputs.NonIndexed())

			log = transactionModel.Log{
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

	log = transactionModel.Log{
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
			ContractAddress: v.ContractAddress,
		}

		txnFee := (float64(v.GasPrice) / float64(params.Ether)) * float64(v.GasUsed)
		transaction.TxnFee = utils.ShowDecimalPrecision(txnFee, 8)

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
			ContractAddress: v.ContractAddress,
		}

		txnFee := (float64(v.GasPrice) / float64(params.Ether)) * float64(v.GasUsed)
		transaction.TxnFee = utils.ShowDecimalPrecision(txnFee, 8)

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
			ContractAddress: v.ContractAddress,
		}

		txnFee := (float64(v.GasPrice) / float64(params.Ether)) * float64(v.GasUsed)
		transaction.TxnFee = utils.ShowDecimalPrecision(txnFee, 8)

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
