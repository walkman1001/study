package block

import (
	"blockchain/transaction"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//区块的结构体
type Block struct {
	Timestamp     int64
	Transactions  []*transaction.Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

//区块交易字段的哈希,因为每个交易ID是序列化并哈希后的交易数据结构，所以我们只需把所以交易的ID进行哈希就可以了
func (b *Block) HashTransactions() []byte {
	var txHash [32]byte
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Hash())
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

//0.3 实现Block的序列化
func (b *Block) Serialize() []byte {
	//首先定义一个buffer存储序列化后的数据
	var result bytes.Buffer
	//实例化一个序列化实例,结果保存到result中
	encoder := gob.NewEncoder(&result)
	//对区块进行实例化
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//0.3 实现反序列化函数
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
