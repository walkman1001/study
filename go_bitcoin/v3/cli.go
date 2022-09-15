package main

import (
	"fmt"
	"os"
)

//接收命令行参数并且控制区块操作

type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA	"add data to blockchain"
	printChain	"print all blockchain data"
`

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
	case "addBlock":
		//添加区块
		fmt.Printf("添加区块")

		//命令校验，验确保参数为4，并且第三个参数为--data
		if len(args) == 4 && args[2] == "--data"{
			//获取数据
			data := args[3]
			//添加区块
			cli.AddBlock(data)
		}else {
			fmt.Printf("添加区块参数使用不当，请检查！")
		}
	case "printChain":
		//打印区块
		fmt.Printf("打印区块\n")
		cli.PrintBlockChain()
	default:
		fmt.Printf("无效命令，请检查！")
		fmt.Printf(Usage)
		
	}
}