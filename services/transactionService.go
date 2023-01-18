package services

import (
	"context"
	"math"
	"strings"
	"time"
	"webbc/DB"
	"webbc/models"
	"webbc/models/transactionModel"
	"webbc/utils"

	"github.com/uptrace/bun"
)

type TransactionServiceImplementation struct {
	database *bun.DB
	ctx      context.Context
}

func NewTransactionService(database *bun.DB, ctx context.Context) TransactionService {
	return &TransactionServiceImplementation{database: database, ctx: ctx}
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
			Value:            utils.ToUint64(v.Value),
			ContractAddress:  v.ContractAddress,
			Status:           v.Status,
			Timestamp:        int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds())),
		}

		if strings.ReplaceAll(oneResultTransaction.To, " ", "") == "" {
			oneResultTransaction.To = ""
		}

		result = append(result, oneResultTransaction)
	}

	return &result, nil
}

func (tsi *TransactionServiceImplementation) GetTransactionByHash(transactionHash string) (*models.Transaction, error) {
	var transaction DB.Transaction
	error1 := tsi.database.NewSelect().Table("transactions").Where("hash = ?", transactionHash).Scan(tsi.ctx, &transaction)

	if error1 != nil {

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
		Value:            utils.ToUint64(transaction.Value),
		ContractAddress:  transaction.ContractAddress,
		Status:           transaction.Status,
		Timestamp:        int(math.Round(time.Now().Sub(time.Unix(int64(transaction.Timestamp), 0)).Seconds())),
	}

	if strings.ReplaceAll(oneResultTransaction.To, " ", "") == "" {
		oneResultTransaction.To = ""
	}

	return &oneResultTransaction, nil
}

func (tsi *TransactionServiceImplementation) GetAllTransactionsInBlock(blockNumber uint64) (*[]models.Transaction, error) {
	var transactions []DB.Transaction
	err := tsi.database.NewSelect().Table("transactions").Where("block_number = ?", blockNumber).Scan(tsi.ctx, &transactions)

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
			Value:            utils.ToUint64(v.Value),
			ContractAddress:  v.ContractAddress,
			Status:           v.Status,
			Timestamp:        int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds())),
		}

		if strings.ReplaceAll(oneResultTransaction.To, " ", "") == "" {
			oneResultTransaction.To = ""
		}

		result = append(result, oneResultTransaction)
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
			Hash:        v.Hash,
			Method:      "",
			BlockNumber: v.BlockNumber,
			Timestamp:   int(math.Round(time.Now().Sub(time.Unix(int64(v.Timestamp), 0)).Seconds())),
			From:        v.From,
			To:          v.To,
			Value:       utils.ToUint64(v.Value),
			TxnFee:      0000000,
		}

		result.Transactions = append(result.Transactions, transaction)
	}

	var totalRows int64
	tsi.database.NewRaw("SELECT reltuples::bigint FROM pg_class WHERE oid = 'public.transactions' ::regclass;").Scan(tsi.ctx, &totalRows)

	result.TotalRows = uint64(totalRows)
	if totalRows > 500000 {
		totalRows = 500000
	}

	totalPages := math.Ceil(float64(totalRows / int64(perPage)))
	result.TotalPages = int(totalPages)
	return &result, nil
}
