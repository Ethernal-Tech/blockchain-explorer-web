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
	Age                string
	DateTime           string
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
	Age              string
	DateTime         string
	Logs             []Log
	InputData        string
	IsToContract     bool
}

type Log struct {
	BlockHash       string
	Index           uint32
	TransactionHash string
	Address         string
	BlockNumber     uint64
	EventName       string
	ParamNames      []string
	ParamTypes      []string
	ParamIndexed    []bool
	Topic0          string
	Topic1          string
	Topic2          string
	Topic3          string
	Data            string
	DataNames       []string
	DataValues      []string
}
