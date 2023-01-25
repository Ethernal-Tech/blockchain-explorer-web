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
	Value       string
	TxnFee      uint64
}
