package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

//定义简单区块结构
type Block struct {
	//版本
	Version uint64
	//前区块哈希
	PrevHash []byte
	//Merkel根（梅克尔根，这就是一个哈希值）
	MerkelRoot []byte
	//时间戳
	TimeStamp uint64
	//难度值
	Difficulty uint64
	//随机数
	Nonce uint64

	//当前区块哈希,正常情况下比特币区块中没有当前区块哈希，我们为了方便做了简化
	Hash []byte
	//数据
	Data []byte
}

//辅助函数，将uint64转化为[]byte
func Uint64ToByte(num uint64) []byte{
	//使用二进制转换
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil{
		log.Panic(err)
	}
	return buffer.Bytes()
}

//添加区块函数
func NewBlock (data string, prevBlockHash []byte) *Block{
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}

	block.SetHash()

	return &block
}

//创建生成当前区块的方法
func (block *Block) SetHash(){
	//拼装数据, 把data打散成一个个byte
	/*
	blockInfo := []byte{}
	blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
	blockInfo = append(blockInfo, block.PrevHash...)
	blockInfo = append(blockInfo, block.MerkelRoot...)
	blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
	blockInfo = append(blockInfo, block.Data...)
	 */
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}
	blockInfo := bytes.Join(tmp, []byte{})

	//对数据进行sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]

}
