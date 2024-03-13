package node

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"scanner/env"
	"scanner/utils"
)

type Node struct {
	env *env.Env

	client *ethclient.Client
}

type NodeImpl interface {
	GetChainID() *big.Int
	GetLatestBlock() uint64
	GetBlockByNumber(number *big.Int) *types.Block
	GetClient() *ethclient.Client
	GetReceiptsByHash(hash common.Hash) *types.Receipt
}

func NewNode(env *env.Env) (NodeImpl, error) {
	n := &Node{
		env: env,
	}

	var err error

	if n.client, err = ethclient.Dial(env.Node.Dial); err != nil {
		panic(err)
	} else {
		return n, nil
	}
}

func (n *Node) GetReceiptsByHash(hash common.Hash) *types.Receipt {
	if res, err := n.client.TransactionReceipt(utils.Context(), hash); err != nil {
		log.Print(err)
		return nil
	} else {
		return res
	}
}

func (n *Node) GetLatestBlock() uint64 {
	if res, err := n.client.BlockNumber(utils.Context()); err != nil {
		log.Print(err)
		return 0
	} else {
		return res
	}
}

func (n *Node) GetBlockByNumber(number *big.Int) *types.Block {
	if res, err := n.client.BlockByNumber(utils.Context(), number); err != nil {
		log.Print(err)
		return nil
	} else {
		return res
	}
}

func (n *Node) GetChainID() *big.Int {
	if res, err := n.client.ChainID(utils.Context()); err != nil {
		// 에러 Log 처리
		log.Print(err)
		return nil
	} else {
		return res
	}
}

func (n *Node) GetClient() *ethclient.Client {
	return n.client
}
