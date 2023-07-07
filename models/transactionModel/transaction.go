package transactionModel

type Transactions struct {
	TotalPages   int
	TotalRows    int64
	Transactions []Transaction
	MaxCount     int
}

type Transaction struct {
	Hash                 string
	Method               string
	BlockHash            string
	BlockNumber          uint64
	From                 string
	To                   string
	Gas                  string
	GasUsed              string
	GasUsedPercentage    float64
	GasPriceInGwei       string
	GasPriceInEth        string
	TxnFee               string
	Nonce                uint64
	TransactionIndex     uint64
	Value                string
	ContractAddress      string
	Status               uint64
	Age                  string
	DateTime             string
	Logs                 []Log
	InputData            string
	InputDataSig         string
	InputDataMethodId    string
	InputDataParamValues []interface{}
	DecodedInputData     DecodedInputData
	Direction            string
	IsToContract         bool
	IsUploadedABI        bool
	ERC20Transfers       []TransferModel
	ERC721Transfers      []TransferModel
	ERC1155Transfers     []TransferModel
}

type DecodedInputData struct {
	FunctionSignature string
	Parameters        []ParameterInfo
}

type ParameterInfo struct {
	Name  string
	Type  string
	Value string
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

type TransferModel struct {
	From            string
	To              string
	Value           string   //erc20 & erc1155
	TokenId         string   //erc721
	Id              string   //erc1155
	Ids             []string //erc1155
	Values          []string //erc1155
	TokenName       string
	TokenAddress    string
	TransactionHash string
	Age             string
	DateTime        string
}
