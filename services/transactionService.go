package services

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"
	"webbc/DB"
	"webbc/models"

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
		fmt.Println("all ok")
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
			Value:            v.Value,
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
		Value:            transaction.Value,
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
	return &[]models.Transaction{}, nil
}
