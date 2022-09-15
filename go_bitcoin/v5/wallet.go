package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	//"golang.org/x/crypto/ripemd160"

	"v5/lib/ripemd160"
	//"github.com/btcsuite/btcutil/base58"

	"v5/lib/base58"
)

//这里的钱包时一结构，每一个钱包保存了公钥,私钥对

type Wallet struct {
	//私钥
	Private *ecdsa.PrivateKey
	//约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）
	PubKey []byte
}

//创建钱包
func NewWallet() *Wallet {
	//创建曲线
	curve := elliptic.P256()
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	//生成公钥
	pubKeyOrig := privateKey.PublicKey

	//拼接X, Y
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)

	return &Wallet{Private: privateKey, PubKey: pubKey}
}

//生成地址
func (w *Wallet) NewAddress() string {
	//钱包公钥
	pubKey := w.PubKey

	//计算公钥哈希和ripe60
	rip160HashValue := HashPubKey(pubKey)
	version := byte(00)
	//拼接version
	payload := append([]byte{version}, rip160HashValue...)

	//校验码checksum
	checkCode := CheckSum(payload)

	//拼接版本、哈希值、校验码、25字节数据
	payload = append(payload, checkCode...)

	//base58编码
	address := base58.Encode(payload)

	return address
}

//Ripe160
func HashPubKey(data []byte) []byte {
	hash := sha256.Sum256(data)

	//理解为编码器
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])

	if err != nil {
		log.Panic(err)
	}

	//返回rip160的哈希结果
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
}

//校验码
func CheckSum(data []byte) []byte {
	//两次sha256
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])

	//前4字节校验码
	checkCode := hash2[:4]
	return checkCode
}
