package main

import (
	"log"
	"v3/bolt"
)

const blockChianDb = "blockChain.db"
const blockBucket = "blockBucket"

//使用数值引入区块链
type BlockChain struct {
	//定义一个区块类型的数组
	//blocks []*Block

	//使用数据库代替数组
	//key是区块的hash值，value为区块的字节流
	db *bolt.DB
	//存储最后一个区块的哈希
	tail []byte
}

//创建区块链
func NewBlockChain () *BlockChain{

	//最后一个区块的哈希,从数据库中读出来的
	var lastHash []byte

	//打开数据库
	db, err := bolt.Open(blockChianDb, 0600, nil)
	//defer db.Close()
	if err != nil{
		log.Panic(err)
	}

	//写数据
	db.Update(func(tx *bolt.Tx) error {
		//找到bucket,(如果没有就创建，没有要找的bucket就代表要对一个新链进行操作，否则就是已有的链，进行追加即可)
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有bucket，创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket(b1)失败")
			}

			//定义创世块
			genesisBlock := GenesisBlock()
			//block的哈希作为key，block的字节流作为value
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			//修改最后一个区块的哈希
			bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash

			//测试
			//blockBytes := bucket.Get(genesisBlock.Hash)
			//block := Deserialize(blockBytes)
			//fmt.Printf("block info: %v\n", block)
		}else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		//return nil代表整个事务操作完成，不需要回滚
		return nil
	})

	//返回刚刚操作的区块链
	return &BlockChain{
		db:   db,
		tail: lastHash,
	}
}
//定义创世块
func  GenesisBlock() *Block {
	return NewBlock("创世块！\n", []byte{})

}

//添加区块到区块链
func (bc *BlockChain) AddBlock (data string) {

	//获取区块链
	db := bc.db
	//获取最后一个区块哈希
	lastHash := bc.tail

	db.Update(func(tx *bolt.Tx) error {

		//完成区块添加
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket 不应该为空，请检查！")
		}

		//1. 创建新区块
		block := NewBlock(data, lastHash)

		//2. 添加区块到数据库中
		//hash作为key, block的字节流作为value
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("LastHashKey"), block.Hash)

		//3. 更新内存中的区块链
		bc.tail = block.Hash

		return nil
	})
}

