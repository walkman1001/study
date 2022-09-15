package main


func main () {

	//创建一个区块链
	bc := NewBlockChain("1F1CRxgLQ2tvuCPeHx17SQXYutjLJ5XNkj")
	//调用命令行命令
	cli := CLI{bc}
	//处理相应请求
	cli.Run()


}