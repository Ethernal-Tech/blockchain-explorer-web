package blockModel

type Blocks struct {
	TotalPages int
	TotalRows  int64
	Blocks     []Block
}

type Block struct {
	Number             uint64
	Age                string
	DateTime           string
	TransactionsNumber int
	Validator          string
	GasUsed            uint64
	GasLimit           uint64
}
