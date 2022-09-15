package main

func main () {

	bc := NewBlockChain("我的地址")
	cli := CLI{bc}
	cli.Run()


}