package main

import (
	"fmt"
	"os"
	"strconv"
)

//接收命令行参数并且控制区块操作

type CLI struct {
	bc *BlockChain
}

const Usage = `
	printChain				正向打印区块链
	
	printChainR				反向打印区块链
	
	getBalance --address ADDRESS		获取指定地址的余额"
	
	send FROM TO AMONUNT MINER DATA		由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA
	
	newWallet				创建一个新的钱包（私钥公钥对）
	
	listAddresses				列举所有的钱包地址
	
`

//根据命令行参数处理对应请求
func (cli *CLI) Run() {

	//1. 获取命令
	args := os.Args
	//	校验参数是否准确
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}

	//2. 分析命令
	cmd := args[1]
	switch cmd {
	case "printChain":
		//打印区块
		fmt.Printf("正向打印区块\n")
		cli.PrintBlockChain()
	case "printChainR":
		//反向打印区块
		fmt.Printf("反向打印区块\n")
		cli.PrintBlockChainReverse()
	case "getBalance":
		//获取余额
		fmt.Printf("获取余额\n")
		if len(args) == 4 && args[2] == "--address" {
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Printf("转账开始...\n")
		if len(args) != 7{
			fmt.Printf("参数错误，请检查！\n")
			fmt.Printf(Usage)
			return
		}
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(os.Args[4], 64)
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	case "newWallet":
		fmt.Printf("创建新的钱包...\n")
		cli.NewWallet()
	case "listAddresses":
		fmt.Printf("列举所有地址...\n")
		cli.ListAddresses()
	default:
		fmt.Printf("无效命令，请检查！")
		fmt.Printf(Usage)
		
	}
}