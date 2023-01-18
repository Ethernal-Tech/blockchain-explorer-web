package transactionModel

type Transactions struct {
	TotalPages   int
	TotalRows    uint64
	Transactions []Transaction
}

type Transaction struct {
	Hash        string
	Method      string
	BlockNumber uint64
	Timestamp   int
	From        string
	To          string
	Value       uint64
	TxnFee      uint64
}
