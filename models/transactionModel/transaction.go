package transactionModel

type Transactions struct {
	TotalPages   int
	TotalRows    int64
	Transactions []Transaction
}

type Transaction struct {
	Hash            string
	Method          string
	BlockNumber     uint64
	Timestamp       string
	From            string
	To              string
	Value           string
	TxnFee          uint64
	ContractAddress string
}
