package main

import (
	"bytes"
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

	//block.SetHash()
	//创建一个pow对象
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	//根据挖矿结果对区块进行更新
	block.Hash= hash
	block.Nonce = nonce

	return &block
}
