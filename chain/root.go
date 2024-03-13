package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/big"
	"scanner/env"
	"scanner/repo"
	. "scanner/types"
	"scanner/utils"
	"sync/atomic"
	"time"
)

type Chain struct {
	env *env.Env

	chainID *big.Int
	signer  types.EIP155Signer

	repo *repo.Repo
}

func ScanBlock(env *env.Env, repo *repo.Repo, startBlock, endBlock uint64) {
	c := &Chain{
		env:  env,
		repo: repo,
	}

	var err error

	if c.chainID = c.getChainID(); err != nil {
		panic(err)
	} else {
		c.signer = types.NewEIP155Signer(c.chainID)

		c.scanBlock(startBlock, endBlock)
	}
}

func (c *Chain) scanBlock(start, end uint64) {

	startBlock := start

	// end 값이 들어 올때만 병렬처리
	if end != 0 {
		// start 부터 end 까지만 구동, end 에 도달하면, 모듈을 죽임
		c.readBlock(start, end)
	} else {
		// start 부터 최신 블럭을 계속 조회를 하면서 구동
		for {
			time.Sleep(3 * time.Second)

			latestBlock := c.getLatestBlock()

			if latestBlock == uint64(big.NewInt(0).Int64()) {
				log.Print("Failed to get LatestBlock")
			} else if latestBlock < startBlock {
				log.Print("StartBlock over LatestBlock")
			} else {
				go c.readBlock(startBlock, latestBlock)
				atomic.StoreUint64(&startBlock, latestBlock)
			}
		}
	}
}

func (c *Chain) readBlock(start, end uint64) {
	for i := start; i <= end; i++ {
		if blockToRead := c.getBlockByNumber(big.NewInt(int64(i))); blockToRead == nil {
			log.Print("Failed to get Block : ", i)
			continue
		} else if blockToRead.Transactions().Len() == 0 {
			log.Println("Debug Transactions len zero : ", i)
			continue
		} else {
			log.Println("Scan block success save Block & Tx : ", i)

			go c.saveBlock(blockToRead)
			go c.saveTx(blockToRead, blockToRead.Transactions().Len(), blockToRead.Header())
		}
	}
}

func (c *Chain) saveBlock(block *types.Block) {
	if err := c.repo.DB.SaveBlock(MakeCustomBlock(block, c.chainID.Int64())); err != nil {
		log.Println("Failed to save Block : ", block.Hash())
	}
}

func (c *Chain) saveTx(block *types.Block, length int, header *types.Header) {
	var writeModel []mongo.WriteModel

	for j := 0; j < length; j++ {
		tx := block.Transactions()[j]

		if receipt := c.getReceipt(tx.Hash()); receipt == nil {
			log.Println("Failed to get Tx Receipt : ", tx.Hash())
			continue
		} else {

			customTx := MakeCustomTx(tx, receipt, header, c.signer)

			// customTx 에 대한 내용 저장
			if json, err := utils.ToJson(customTx); err != nil {
				log.Println("Failed ToJson", tx.Hash())
				continue
			} else {
				writeModel = append(
					writeModel,
					mongo.NewUpdateOneModel().SetUpsert(true).
						SetFilter(bson.M{"tx": hexutil.Encode(customTx.Tx[:])}). // 소문자 변경
						SetUpdate(bson.M{"$set": json}),
				)
			}
		}
	}
	if len(writeModel) != 0 {
		if err := c.repo.DB.SaveTxByBulk(writeModel); err != nil {
			log.Println("Failed to save Txs : ", block.Hash())
		}
	}

	//if err := c.repo.DB.SaveTx(MakeCustomTx(block, length, header)); err != nil {
	//	log.Println("Failed to save Tx : ", block.Hash())
	//}
}

func (c *Chain) getReceipt(hash common.Hash) *types.Receipt {
	return c.repo.Node.GetReceiptByHash(hash)
}

func (c *Chain) getChainID() *big.Int {
	// Node 호출 => repository 의 Node 함수 출력
	return c.repo.Node.GetChainID()
}

func (c *Chain) getLatestBlock() uint64 {
	return c.repo.Node.GetLatestBlock()
}

func (c *Chain) getBlockByNumber(number *big.Int) *types.Block {
	return c.repo.Node.GetBlockByNumber(number)
}

func (c *Chain) getClient() *ethclient.Client {
	return c.repo.Node.GetClient()
}
