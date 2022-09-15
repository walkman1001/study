package main

import (
	"log"
	"v3/bolt"
)

type BlockChainIterator struct {
	db *bolt.DB
	//游标
	currentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator{
	return &BlockChainIterator{
		db:                 bc.db,
		//最初指向区块链的最后一个区块，随着Next的调用，不断变化
		currentHashPointer: bc.tail,
	}
}

func (it *BlockChainIterator) Next() *Block{
	var block Block

	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("迭代器遍历时，bucket不应该为空！")
		}

		blockTmp := bucket.Get(it.currentHashPointer)

		//解码动作
		block = Deserialize(blockTmp)
		//游标左移
		it.currentHashPointer = block.PrevHash

		return nil
	})

	return &block
}