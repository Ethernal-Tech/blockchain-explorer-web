package nftModel

type NftTransactions struct {
	TotalPages   int
	TotalRows    int64
	Transactions []NftTransaction
	MaxCount     int
}

type NftTransaction struct {
	Hash            string
	Method          string
	Age             string
	DateTime        string
	From            string
	To              string
	Type            string
	Item            string
	IsFromContract  bool
	IsToContract    bool
	NftId           string
	Value           string
	ContractAddress string
}
