package services

import (
	"context"
	"math"
	"strings"
	"time"
	"webbc/DB"
	"webbc/configuration"
	"webbc/models"
	"webbc/models/transactionModel"
	"webbc/utils"

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

		// if strings.ReplaceAll(oneResultTransaction.To, " ", "") == "" {
		// 	oneResultTransaction.To = ""
		// }

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

	error2 := tsi.database.NewSelect().Table("logs").Where("transaction_hash = ?", transactionHash).Scan(tsi.ctx, &logs)

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

	// if strings.ReplaceAll(oneResultTransaction.To, " ", "") == "" {
	// 	oneResultTransaction.To = ""
	// }

	for _, element := range logs {
		var log = models.Log{
			BlockHash:       element.BlockHash,
			Index:           element.Index,
			TransactionHash: element.TransactionHash,
			Address:         element.Address,
			BlockNumber:     element.BlockNumber,
			// Topic0:          strings.ReplaceAll(element.Topic0, " ", ""),
			// Topic1:          strings.ReplaceAll(element.Topic1, " ", ""),
			// Topic2:          strings.ReplaceAll(element.Topic2, " ", ""),
			// Topic3:          strings.ReplaceAll(element.Topic3, " ", ""),
			Topic0: element.Topic0,
			Topic1: element.Topic1,
			Topic2: element.Topic2,
			Topic3: element.Topic3,
			Data:   element.Data,
		}

		oneResultTransaction.Logs = append(oneResultTransaction.Logs, log)
	}

	isToContract, _ := tsi.database.NewSelect().Table("contracts").Where("address = ?", oneResultTransaction.To).Exists(tsi.ctx)
	oneResultTransaction.IsToContract = isToContract

	return &oneResultTransaction, nil
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
		// if strings.ReplaceAll(transaction.To, " ", "") == "" {
		// 	transaction.To = ""
		// }
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
