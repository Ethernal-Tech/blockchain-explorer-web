package DB

type Block struct {
	Hash              string `bun:",pk,type:char(66)"`
	Number            uint64 `bun:"type:bigint,notnull,unique"`
	ParentHash        string `bun:"type:char(66),notnull"`
	Nonce             string `bun:"type:varchar,notnull"`
	Miner             string `bun:"type:char(42),notnull"`
	Difficulty        string `bun:"type:varchar,notnull"`
	TotalDifficulty   string `bun:"type:varchar,notnull"`
	ExtraData         []byte `bun:"type:bytea"`
	Size              uint64 `bun:"type:bigint,notnull"`
	GasLimit          uint64 `bun:"type:bigint,notnull"`
	GasUsed           uint64 `bun:"type:bigint,notnull"`
	Timestamp         uint64 `bun:"type:bigint,notnull"`
	TransactionsCount int    `bun:"type:integer,notnull"`
}

type Transaction struct {
	Hash             string `bun:",pk,type:char(66)"`
	BlockHash        string `bun:"type:char(66),notnull"`
	BlockNumber      uint64 `bun:"type:bigint,notnull"`
	From             string `bun:"type:char(42),notnull"`
	To               string `bun:"type:char(42)"`
	Gas              uint64 `bun:"type:bigint,notnull"`
	GasUsed          uint64 `bun:"type:bigint,notnull"`
	GasPrice         uint64 `bun:"type:bigint,notnull"`
	Nonce            uint64 `bun:"type:bigint,notnull"`
	TransactionIndex uint64 `bun:"type:integer,notnull"`
	Value            string `bun:"type:varchar"`
	ContractAddress  string `bun:"type:varchar(42)"`
	Status           uint64 `bun:"type:smallint,notnull"`
	Timestamp        uint64 `bun:"type:bigint,notnull"`
	InputData        string `bun:"type:varchar"`
}
