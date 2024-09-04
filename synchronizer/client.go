package synchronizer

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
)

func EthClient(ctx context.Context, url string) (*ethclient.Client, error) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	//TODO add retry
	client, err := ethclient.DialContext(ctx2, url)
	if err != nil {
		return nil, err
	}
	return client, nil
}
