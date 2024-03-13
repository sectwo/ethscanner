package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type CustomBlock struct {
	Hash       common.Hash    `json:"hash"`
	ParentHash common.Hash    `json:"parentHash"`
	Miner      common.Address `json:"miner"`
	Root       common.Hash    `json:"root"`

	Time        uint64 `json:"time"`
	BlockNumber uint64 `json:"blockNumber"`
	Size        uint64 `json:"size"`

	ChainID int64 `json:"chainID"`
}

type CustomTx struct {
	Tx   common.Hash     `json:"tx"`
	From common.Address  `json:"from"`
	To   *common.Address `json:"to"`

	BlockNumber *big.Int `json:"blockNumber"`
	Fee         *big.Int `json:"fee"`

	Size   uint64 `json:"size"`
	Amount string `json:"amount"`
	Nonce  uint64 `json:"nonce"`
	Time   int64  `json:"time"`
}

func MakeCustomBlock(block *types.Block, chainID int64) *CustomBlock {
	customBlock := &CustomBlock{
		Hash:        block.Hash(),
		ParentHash:  block.ParentHash(),
		Miner:       block.Coinbase(),
		Root:        block.Root(),
		Time:        block.Time(),
		BlockNumber: block.NumberU64(),
		Size:        block.Size(),

		ChainID: chainID,
	}

	return customBlock
}

func MakeCustomTx(transaction *types.Transaction, receipt *types.Receipt, header *types.Header, signer types.EIP155Signer) *CustomTx {
	tx := &CustomTx{
		Tx:          receipt.TxHash,
		To:          transaction.To(),
		BlockNumber: header.Number,
		Fee:         new(big.Int).Mul(transaction.GasPrice(), big.NewInt(int64(receipt.GasUsed))), // gas price 에서 사용된 gas 의 양을 곱해야함
		Size:        transaction.Size(),
		Amount:      transaction.Value().String(),
		Nonce:       transaction.Nonce(),
		Time:        transaction.Time().Unix(),
	}

	// From
	tx.From, _ = types.Sender(signer, transaction)

	return tx
}
