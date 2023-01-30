package addressModel

type Address struct {
	AddressHex       string
	Balance          string
	TransactionCount int
	Transactions     []Transaction
}

type Transaction struct {
	Hash        string
	Method      string
	BlockNumber uint64
	Timestamp   string
	From        string
	To          string
	Direction   string
	Value       string
	Gas         uint64
	GasUsed     uint64
	GasPrice    uint64
}
