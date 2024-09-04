package synchronizer

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"strings"
)

const (
	RPC_URL               = "wss://mainnet.infura.io/ws/v3/fbc59ac1af714c0e8f92a67ead436d49"
	MAX_LEN_OF_WATCH_LIST = 100
	TYP1                  = "from"
	TYP2                  = "to"
	ADDRESS               = "0x0604cc2a4d90d0854d4551133c31d6c55232c749"
)

type Monitor struct {
	client *ethclient.Client
	//TODO 去重结构
	watchList map[string][]string
}

func (m *Monitor) Start(ctx context.Context) error {
	fmt.Print("ROW27")
	//TODO implement me
	headers := make(chan *types.Header)
	fmt.Println("ROW30")
	sub, err := m.client.SubscribeNewHead(context.Background(), headers)
	fmt.Println("ROW32")
	log.Debug("subscribed to new headers")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	//defer sub.Unsubscribe()
	fmt.Println("WTF")
	for {
		select {
		case err := <-sub.Err():
			log.Error(err.Error())
			return err
		case header := <-headers:
			fmt.Println("fuck")
			m.processBlock(m.client, header.Number.Uint64())
		}
	}
}

// TODO add db
func (m *Monitor) processBlock(client *ethclient.Client, blockNumber uint64) {
	//TODO implement me
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, tx := range block.Transactions() {
		var from common.Address
		switch tx.Type() {
		case types.LegacyTxType:
			signer := types.NewEIP155Signer(tx.ChainId())
			from, err = types.Sender(signer, tx)
		case types.AccessListTxType:
			signer := types.NewEIP2930Signer(tx.ChainId())
			from, err = types.Sender(signer, tx)
		case types.DynamicFeeTxType:
			signer := types.NewLondonSigner(tx.ChainId())
			from, err = types.Sender(signer, tx)
		case types.BlobTxType: // 这是 EIP-4844 引入的新类型
			signer := types.NewCancunSigner(tx.ChainId())
			from, err = types.Sender(signer, tx)
		default:
			log.Warn("Unsupported transaction type", "type", tx.Type())
			continue
		}

		to := tx.To()
		fromAddr := strings.ToLower(from.Hex())
		toAddr := "contract creation"
		if to != nil {
			toAddr = strings.ToLower(to.Hex())
		}

		// 移除 "0x" 前缀进行比较
		addressToCompare := strings.ToLower(strings.TrimPrefix(ADDRESS, "0x"))
		fromAddrCompare := strings.TrimPrefix(fromAddr, "0x")
		toAddrCompare := strings.TrimPrefix(toAddr, "0x")

		if fromAddrCompare == addressToCompare || toAddrCompare == addressToCompare {
			fmt.Println("Matching transaction found!")
			fmt.Println("Block number:", blockNumber)
			fmt.Println("From:", fromAddr)
			fmt.Println("To:", toAddr)
			fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())
			// 这里可以添加更多的交易信息打印
		}
	}
}

func (m *Monitor) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *Monitor) Stopped(ctx context.Context) bool {
	//TODO implement me
	panic("implement me")
}

func (m *Monitor) AddAddress(ctx context.Context, typ string, address string) error {
	//TODO implement me
	l := len(m.watchList)
	if typ != TYP1 && typ != TYP2 {
		return fmt.Errorf("type is not valid")
	}
	if l+1 > MAX_LEN_OF_WATCH_LIST {
		//TODO error
		return fmt.Errorf("watch list is full, max length is %d", MAX_LEN_OF_WATCH_LIST)
	}
	m.watchList[typ] = append(m.watchList[typ], address)
	return nil
}

func (m *Monitor) AddAddresses(ctx context.Context, list map[string][]string) error {
	l := len(m.watchList)
	if l+len(list) > MAX_LEN_OF_WATCH_LIST {
		//TODO error
		return fmt.Errorf("watch list is full, max length is %d", MAX_LEN_OF_WATCH_LIST)
	}
	for k, v := range list {
		if k != TYP1 && k != TYP2 {
			return fmt.Errorf("type is not valid")
		}
		m.watchList[k] = append(m.watchList[k], v...)
	}
	return nil
}

type IMonitor interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Stopped(ctx context.Context) bool
	AddAddress(ctx context.Context, typ string, address string) error
	//TODO 注意长度限制
	AddAddresses(ctx context.Context, list map[string][]string) error
}

func NewMonitor(ctx context.Context) *Monitor {
	cli, err := EthClient(ctx, RPC_URL)
	if err != nil {
		panic(err)
	}

	return &Monitor{
		client: cli,
		watchList: map[string][]string{
			TYP1: []string{ADDRESS},
			TYP2: []string{ADDRESS},
		},
	}
}
