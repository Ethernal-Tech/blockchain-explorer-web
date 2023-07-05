package blockModel

type Blocks struct {
	TotalPages int
	TotalRows  int64
	Blocks     []Block
}

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
	Age                string
	DateTime           string
	TransactionsNumber int
}
