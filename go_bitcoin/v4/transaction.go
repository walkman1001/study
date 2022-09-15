package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 12.5

//1. 定义交易结构
type Transaction struct {
	//交易ID
	TXID []byte
	//交易输入数组
	TXInputs []TXInput
	//交易输出数组
	TXOutputs []TXOutput
}

//定义交易输入
type TXInput struct {
	//引用的交易ID
	Txid []byte
	//引用的output的索引值
	Index int64
	//解锁脚本，我们用地址模拟
	Sig string
}

//定义交易输出
type TXOutput struct {
	//转账金额
	Value float64
	//锁定脚本，我们用地址模拟
	PubKeyHash string
}

//设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

//实现一个函数，判断当前交易是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].Txid) == 0 && tx.TXInputs[0].Index == -1{
		return true
	}

	return false
}

//2. 创建挖矿交易 coinbase
func NewCoinbaseTX(address string, data string) *Transaction{
	//挖矿交易特点,
	//	1. 只有一个input
	//	2. 无需引用交易id
	//	3. 无需引用index
	//矿工由于挖矿时无需指定签名，所以这个sig字段可以由矿工自由填写，一般写矿池的名字
	input := TXInput{
		Txid:  []byte{},
		Index: -1,
		Sig:   "矿池信息",
	}
	output := TXOutput{
		Value:      reward,
		PubKeyHash: address,
	}

	//对于挖矿交易，只有一个input和output
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
	}

	//设置交易ID
	tx.SetHash()

	return &tx
}

//创建普通转账交易
//1. 找到最合理的UTXO集合，map[string]uint64
//2. 将这些UTXO逐一转成inputs
//3. 创建outputs
//4. 如果有零钱，要找零
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction{
	//1. 找到最合理的UTXO集合，map[string]uint64
	utxos, resValue := bc.FindNeedUTXOs(from, amount)

	if resValue < amount{
		fmt.Printf("余额不足，交易失败！\n")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput
	//2. 创建交易输入将这些UTXO逐一转成inputs
	for id, indexArray := range utxos{
		for _, i := range indexArray{
			input := TXInput{
				Txid:  []byte(id),
				Index: int64(i),
				Sig:   from,
			}
			inputs = append(inputs, input)
		}
	}

	//创建交易输出
	output := TXOutput{amount, to}
	outputs= append(outputs, output)

	//找零
	if resValue > amount{

		outputs = append(outputs, TXOutput{resValue - amount, from})
	}

	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  inputs,
		TXOutputs: outputs,
	}

	tx.SetHash()

	return &tx
}

