package main

import (
	"io/ioutil"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/elliptic"
	"os"
)

const walletFile = "wallet.dat"

//定一个 Wallets结构，它保存所有的wallet以及它的地址
type Wallets struct {
	//map[地址]钱包
	WalletsMap map[string]*Wallet
}

//创建（加载）钱包集，返回当前所有钱包的实例
func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	ws.loadFile()
	return &ws
}

//创建钱包到钱包集
func (ws *Wallets) CreateWallet() string {
	//创建一个钱包
	wallet := NewWallet()
	address := wallet.NewAddress()

	//添加到钱包集
	ws.WalletsMap[address] = wallet

	//保存包本地
	ws.saveToFile()
	//返回创建钱包的地址
	return address
}

//保存方法，把新建的wallet添加进去
func (ws *Wallets) saveToFile() {

	var buffer bytes.Buffer

	//panic: gob: type not registered for interface: elliptic.p256Curve
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	//一定要注意校验！！！
	if err != nil {
		log.Panic(err)
	}

	ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
}

//读取文件方法，把所有的wallet读出来
func (ws *Wallets) loadFile() {
	//在读取之前，要先确认文件是否在，如果不存在，直接退出
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		//ws.WalletsMap = make(map[string]*Wallet)
		return
	}

	//读取内容
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	//解码
	gob.Register(elliptic.P256())

	decoder := gob.NewDecoder(bytes.NewReader(content))

	var wsLocal Wallets

	err = decoder.Decode(&wsLocal)
	if err != nil {
		log.Panic(err)
	}

	ws.WalletsMap = wsLocal.WalletsMap
}

func (ws *Wallets) ListAllAddresses() []string {
	var addresses []string
	//遍历钱包，将所有的key取出来返回
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}

	return addresses
}
