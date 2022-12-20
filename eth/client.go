package eth

import (
	"github.com/ethereum/go-ethereum/rpc"
)

type BlockchainNodeConnection struct {
	HTTP      *rpc.Client
	WebSocket *rpc.Client
}

// Connect to blockchain node, either using HTTP or Websocket connection depending on URL passed to function
func GetClient(rpcUrl string) *rpc.Client {

	rpcClient, err := rpc.Dial(rpcUrl)
	if err != nil {
		//TODO: Error handling
	}

	return rpcClient
}
