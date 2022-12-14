package main
//
//import (
//	"crypto/ecdsa"
//	"crypto/elliptic"
//	"crypto/rand"
//	"crypto/sha256"
//	"fmt"
//	"log"
//	"math/big"
//)
//
////1. 演示如何使用ecdsa生成公私钥
////2. 签名校验
//
//func main() {
//	//创建曲线
//	curve := elliptic.P256()
//
//	//生成私钥，曲线、随机数
//	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
//	if err != nil{
//		log.Panic(err)
//	}
//	//生成公钥
//	pubKey := privateKey.PublicKey
//
//	data := "hello world!"
//	hash := sha256.Sum256([]byte(data))
//
//	//签名
//	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
//	if err != nil{
//		log.Panic(err)
//	}
//
//	//把r,s进行序列化
//	signature := append(r.Bytes(), s.Bytes()...)
//	//1. 定义两个辅助的big.int
//	r1 := big.Int{}
//	s1 := big.Int{}
//	//2. 拆分我们的signature，平分前半部分给r,后半部分给s
//	r1.SetBytes(signature[:len(signature)/2])
//	s1.SetBytes(signature[len(signature)/2:])
//	//校验：数据、签名、公钥
//	res := ecdsa.Verify(&pubKey, hash[:], &r1, &s1)
//	fmt.Printf("校验结果：%v", res)
//}