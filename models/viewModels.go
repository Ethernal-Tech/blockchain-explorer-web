package models

type Block struct {
	Hash               string
	Number             uint64
	ParentHash         string
	Nonce              string
	Validator          string
	Difficulty         string
	TotalDifficulty    string
	ExtraData          []byte
	Size               uint64
	GasLimit           uint64
	GasUsed            uint64
	Timestamp          string
	TransactionsNumber int
}

type Transaction struct {
	Hash             string
	BlockHash        string
	BlockNumber      uint64
	From             string
	To               string
	Gas              uint64
	GasUsed          uint64
	GasPrice         uint64
	Nonce            uint64
	TransactionIndex uint64
	Value            string
	ContractAddress  string
	Status           uint64
	Timestamp        string
}
