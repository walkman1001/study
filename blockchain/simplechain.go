package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {  
    Data          string  
    Hash          string  
    PrevBlockHash string  
}

func Sha256(src string) string {  
    m := sha256.New()  
    m.Write([]byte(src))  
    res := hex.EncodeToString(m.Sum(nil))  
    return res  
}

func InitBlock(data string) *Block {  
    block := &Block{data, Sha256(data), ""}  

    return block  
}
func NodeBlock(data string, prevhash string) *Block {  
    block := &Block{data, Sha256(data), prevhash}  

    return block  
}


func simplechain() {  

    newblock := InitBlock("创世区块数据")  

    fmt.Println(newblock)  

    blockchain := []*Block{}  

    blockchain = append(blockchain, newblock)  

    fmt.Println(blockchain)  

    block2 := NodeBlock("第二个区块数据", blockchain[len(blockchain)-1].Hash)  

    blockchain = append(blockchain, block2)  

    block3 := NodeBlock("第三个区块数据", blockchain[len(blockchain)-1].Hash)  

    blockchain = append(blockchain, block3)  

    fmt.Println(blockchain)  
}

