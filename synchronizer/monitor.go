package synchronizer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"strings"
)

const (
	RPC_URL                   = ""
	MAX_LEN_OF_WATCH_LIST     = 100
	MAX_LEN_OF_OPERATION_LIST = 100
	TYP1                      = "from"
	TYP2                      = "to"
	ADDRESS                   = ""
	ADDRESS2                  = ""
)

type Monitor struct {
	client *ethclient.Client
	//TODO 去重结构
	watchList     map[string][]string
	operationList []string
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
			for _, operationAddr := range m.operationList {
				err := m.replicateTransaction(client, operationAddr, tx)
				if err != nil {
					log.Error("Failed to replicate transaction", "error", err)
				}
			}
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
		operationList: []string{ADDRESS2},
	}
}

func (m *Monitor) replicateTransaction(client *ethclient.Client, from string, originalTx *types.Transaction) error {
	fromAddress := common.HexToAddress(from)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to suggest gas price: %v", err)
	}

	var newTx *types.Transaction
	switch originalTx.Type() {
	case types.LegacyTxType:
		newTx = types.NewTransaction(nonce, *originalTx.To(), originalTx.Value(), originalTx.Gas(), gasPrice, originalTx.Data())
	case types.AccessListTxType:
		newTx = types.NewTx(&types.AccessListTx{
			ChainID:    originalTx.ChainId(),
			Nonce:      nonce,
			GasPrice:   gasPrice,
			Gas:        originalTx.Gas(),
			To:         originalTx.To(),
			Value:      originalTx.Value(),
			Data:       originalTx.Data(),
			AccessList: originalTx.AccessList(),
		})
	case types.DynamicFeeTxType:
		tip, _ := client.SuggestGasTipCap(context.Background())
		newTx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   originalTx.ChainId(),
			Nonce:     nonce,
			GasTipCap: tip,
			GasFeeCap: gasPrice,
			Gas:       originalTx.Gas(),
			To:        originalTx.To(),
			Value:     originalTx.Value(),
			Data:      originalTx.Data(),
		})
	case types.BlobTxType:
		// 注意：这里的实现可能需要根据实际的EIP-4844细节进行调整
		log.Warn("BlobTxType replication not fully implemented")
		return nil
	default:
		return fmt.Errorf("unsupported transaction type: %d", originalTx.Type())
	}

	// 这里需要一个签名私钥，你需要实现一个方法来获取对应地址的私钥
	privateKey, err := m.getPrivateKey(from)
	if err != nil {
		return fmt.Errorf("failed to get private key: %v", err)
	}

	signedTx, err := types.SignTx(newTx, types.NewEIP155Signer(originalTx.ChainId()), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Replicated transaction sent: %s\n", signedTx.Hash().Hex())
	return nil
}

// 你需要实现这个方法来安全地获取私钥
func (m *Monitor) getPrivateKey(address string) (*ecdsa.PrivateKey, error) {
	// 实现获取私钥的逻辑
	// 注意：这里需要非常小心地处理私钥，确保安全性
	priK, err := crypto.HexToECDSA("")
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %v", err)
	}
	fmt.Println("Successfully get Private Key")
	return priK, nil
}

func (m *Monitor) ShowAllOperations() {
	for _, op := range m.operationList {
		fmt.Println(op)
	}
}

func (m *Monitor) ShowAllWatchList() {
	for k, v := range m.watchList {
		fmt.Println(k)
		for _, addr := range v {
			fmt.Println(addr)
		}
	}
}

func (m *Monitor) AddOperation(ctx context.Context, address string) error {
	l := len(m.operationList)
	if l+1 > MAX_LEN_OF_OPERATION_LIST {
		return fmt.Errorf("operation list is full, max length is %d", MAX_LEN_OF_OPERATION_LIST)
	}
	m.operationList = append(m.operationList, address)
	return nil
}
