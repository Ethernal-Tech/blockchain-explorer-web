package addressModel

import "math/big"

type Address struct {
	AddressHex       string
	Balance          *big.Float
	TransactionCount int
	Transactions     []Transaction
}

type Transaction struct {
	Hash        string
	Method      string
	BlockNumber uint64
	Timestamp   int
	From        string
	To          string
	Direction   string
	Value       uint64
	Gas         uint64
	GasUsed     uint64
	GasPrice    uint64
}
