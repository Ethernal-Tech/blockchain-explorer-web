package transactionModel

type Transactions struct {
	TotalPages   int
	TotalRows    int64
	Transactions []Transaction
	MaxCount     int
}

type Transaction struct {
	Hash            string
	Method          string
	BlockNumber     uint64
	Age             string
	DateTime        string
	From            string
	To              string
	Value           string
	TxnFee          string
	ContractAddress string
	Direction       string
	IsToContract    bool
}
