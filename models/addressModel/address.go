package addressModel

type Address struct {
	AddressHex         string
	Balance            string
	TransactionCount   int
	Transactions       []Transaction
	IsContract         bool
	CreatorAddress     string
	CreatorTransaction string
}

type Transaction struct {
	Hash            string
	Method          string
	BlockNumber     uint64
	Age             string
	DateTime        string
	From            string
	To              string
	Direction       string
	Value           string
	Gas             uint64
	GasUsed         uint64
	GasPrice        uint64
	ContractAddress string
}
