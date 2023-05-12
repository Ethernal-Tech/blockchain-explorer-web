package eth

import (
	"github.com/ethereum/go-ethereum/rpc"
)

var httpNodeClient *rpc.Client

func NewHttpNodeClient(rpcUrl string) *rpc.Client {
	var err error
	httpNodeClient, err = rpc.Dial(rpcUrl)
	if err != nil {
		//TODO: Error handling
	}

	return httpNodeClient
}

func GetHttpNodeClient() *rpc.Client {
	return httpNodeClient
}
