package main

import (
	"blockchain/CLI"
	"blockchain/sample"
)

/*

git 提交的url

https://github.com/walkman1001/study.git


*/
func main() {
	//start02()
	//cryptoSample()
	//base64Sample()
	//base58aSample()
	sample.Base58bSample()
	cli := CLI.CLI{}
	cli.Run()
}
