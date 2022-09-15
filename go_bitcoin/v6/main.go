package main

func main () {

	bc := NewBlockChain("1F1CRxgLQ2tvuCPeHx17SQXYutjLJ5XNkj")
	cli := CLI{bc}
	cli.Run()


}