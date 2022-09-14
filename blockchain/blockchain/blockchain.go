package blockchain
 
import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"blockchain/block"
	"blockchain/pow"
	"blockchain/transaction"
	"blockchain/wallet"
	"log"
	"errors"
)
/*
	区块链实现
*/
const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
//区块链
type Blockchain struct {
	tip		[]byte
	db 		*bolt.DB
}
 
//工厂模式db
func(bc *Blockchain) Db() *bolt.DB {
	return bc.db
}
 
//把区块添加进区块链,挖矿
func (bc *Blockchain) MineBlock(transactions []*transaction.Transaction) {
	var lastHash []byte
 
	//在一笔交易被放入一个块之前进行验证
	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: 无效 transaction")
		}
	}
	//只读的方式浏览数据库，获取当前区块链顶端区块的哈希，为加入下一区块做准备
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))	//通过键"l"拿到区块链顶端区块哈希
 
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
 
	//prevBlock := bc.Blocks[len(bc.Blocks)-1]
	//求出新区块
	newBlock := pow.NewBlock(transactions,lastHash)
	// bc.Blocks = append(bc.Blocks,newBlock)
	//把新区块加入到数据库区块链中
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash,newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"),newBlock.Hash)
		bc.tip = newBlock.Hash
 
		return nil
	})
}
 
//创建创世区块  /修改/
func NewGenesisBlock(coinbase *transaction.Transaction) *block.Block {
	return pow.NewBlock([]*transaction.Transaction{coinbase},[]byte{})
}
//实例化一个区块链,默认存储了创世区块 ,接收一个地址为挖矿奖励地址 /修改/
func NewBlockchain(address string) *Blockchain {
	//return &Blockchain{[]*block.Block{NewGenesisBlock()}}
	var tip []byte
	//打开一个数据库文件，如果文件不存在则创建该名字的文件
	db,err := bolt.Open(dbFile,0600,nil)
	if err != nil {
		log.Panic(err)
	}
	//读写操作数据库
	err = db.Update(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(blocksBucket))
		//查看名字为blocksBucket的Bucket是否存在
		if b == nil {
			//不存在则从头 创建
			fmt.Println("不存在区块链，需要重新创建一个区块链...")
			
			//genesis := NewGenesisBlock()	//创建创世区块
			//此时的创世区块就要包含交易coinbaseTx
			cbtx := transaction.NewCoinbaseTX(address,genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)
 
			b,err := tx.CreateBucket([]byte(blocksBucket)) //创建名为blocksBucket的桶
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash,genesis.Serialize()) //写入键值对，区块哈希对应序列化后的区块
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"),genesis.Hash) //"l"键对应区块链顶端区块的哈希
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash //指向最后一个区块，这里也就是创世区块
		} else {
			//如果存在blocksBucket桶，也就是存在区块链
			//通过键"l"映射出顶端区块的Hash值
			tip = b.Get([]byte("l"))
		}
 
		return nil
	})
 
	bc := Blockchain{tip,db}  //此时Blockchain结构体字段已经变成这样了
	return &bc
}
 
//分割线——————迭代器——————
type BlockchainIterator struct {
	currentHash 	[]byte
	db 				*bolt.DB
}
//当需要遍历当前区块链时，创建一个此区块链的迭代器
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip,bc.db}
 
	return bci
}
 
//迭代器的任务就是返回链中的下一个区块
func (i *BlockchainIterator) Next() *block.Block {
	var Block *block.Block
 
	//只读方式打开区块链数据库
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		//获取数据库中当前区块哈希对应的被序列化后的区块
		encodeBlock := b.Get(i.currentHash)
		//反序列化，获得区块
		Block = block.DeserializeBlock(encodeBlock)
 
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
 
	//把迭代器中的当前区块哈希设置为上一区块的哈希，实现迭代的作用
	i.currentHash =Block.PrevBlockHash
 
	return Block
 
}
 
//在区块链上找到每一个区块中属于address用户的未花费交易输出,返回未花费输出的交易切片
func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []transaction.Transaction {
	var unspentTXs []transaction.Transaction
	//创建一个map，存储已经花费了的交易输出
	spentTXOs := make(map[string][]int)
	//因为要在链上遍历区块，所以要使用到迭代器
	bci := bc.Iterator()
 
	for {
		block := bci.Next()  //迭代
 
		//遍历当前区块上的交易
		for _,tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) //把交易ID转换成string类型，方便存入map中
		
		//标签
		Outputs:
		//遍历当前交易中的输出切片，取出交易输出
			for outIdx,out := range tx.Vout {
				//在已经花费了的交易输出map中，如果没有找到对应的交易输出，则表示当前交易的输出未花费
				//反之如下
				if spentTXOs[txID] != nil {
					//存在当前交易的输出中有已经花费的交易输出，
					//则我们遍历map中保存的该交易ID对应的输出的index 
					//提示：(这里的已经花费的交易输出index其实就是输入TXInput结构体中的Vout字段)
					for _,spentOutIdx := range spentTXOs[txID] {
						//首先要清楚当前交易输出是一个切片，里面有很多输出，
						//如果map里存储的引用的输出和我们当前遍历到的输出index重合,则表示该输出被引用了
						if spentOutIdx == outIdx {
							continue Outputs  //我们就继续遍历下一轮，找到未被引用的输出
						}
					}
				}
				//到这里是得到此交易输出切片中未被引用的输出
 
				//这里就要从这些未被引用的输出中筛选出属于该用户address地址的输出
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs,*tx)
				}
			}
			//判断是否为coinbase交易
			if tx.IsCoinbase() == false { 		
				//如果不是,则遍历当前交易的输入
				for _,in := range tx.Vin {
					//如果当前交易的输入是被该地址address所花费的，就会有对应的该地址的引用输出
					//则在map上记录该输入引用的该地址对应的交易输出
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID],in.Vout)
					}
				}
			}
		}
		//退出for循环的条件就是遍历到的创世区块后
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}
 
//通过找到未花费输出交易的集合，我们返回集合中的所有未花费的交易输出
func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []transaction.TXOutput {
	var UTXOs []transaction.TXOutput
	//找到address地址下的未花费交易输出的交易的集合
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)
	//遍历交易集合得到交易，从交易中提取出输出字段Vout,从输出字段中提取出属于address的输出
	for _,tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs,out)
			}
		}
	}
	//返回未花费交易输出
	return UTXOs 
}
 
//发送币操作,相当于创建一笔未花费输出交易
func NewUTXOTransaction(from,to string,amount int,bc *Blockchain) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput
	//validOutputs是一个存放要用到的未花费输出的交易/输出的map 
	//acc,validOutputs := bc.FindSpendableOutputs(from,amount)
	wallets,err := wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	_wallet := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(_wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)
	if acc < amount {
		log.Panic("ERROR:Not enough tokens...")
	}
	//通过validOutputs里面的数据来放入建立一个输入列表
	for txid,outs := range validOutputs {
		//反序列化得到txID
		txID,err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		//遍历输出outs切片,得到TXInput里的Vout字段值
		for _,out := range outs {
			//input := transaction.TXInput{txID,out,from}
			input := transaction.TXInput{txID,out,nil,_wallet.PublicKey}
			inputs = append(inputs,input)
		}
	}
	//建立一个输出列表
	//outputs = append(outputs,transaction.TXOutput{amount,to})
	outputs = append(outputs,*transaction.NewTXOutput(amount,to))
	if acc > amount {
		//outputs = append(outputs,transaction.TXOutput{acc - amount,from}) //相当于找零
		outputs = append(outputs,*transaction.NewTXOutput(acc - amount,from)) //相当于找零
	}
	tx := transaction.Transaction{nil,inputs,outputs}
	//tx.SetID()
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, _wallet.PrivateKey)
 
	return &tx
}
 
//找到可以花费的交易输出,这是基于上面的FindUnspentTransactions 方法
func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte,amount int) (int,map[string][]int) {
	//未花费交易输出map集合
	unspentOutputs := make(map[string][]int)
	//未花费交易
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0	//累加未花费交易输出中的Value值
 
	Work:
		for _,tx := range unspentTXs {
			txID := hex.EncodeToString(tx.ID)
 
			for outIdx,out := range tx.Vout {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID],outIdx)
 
					if accumulated >= amount {
						break Work
					}
				}
			}
		}
		return accumulated,unspentOutputs
}
 
//通过交易ID找到一个交易
func (bc *Blockchain) FindTransaction(ID []byte) (transaction.Transaction,error) {
	bci := bc.Iterator()
	for {
		block := bci.Next()
 
		for _,tx := range block.Transactions {
			if bytes.Compare(tx.ID,ID) == 0 {
				return *tx,nil
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return transaction.Transaction{},errors.New("Transaction is not found")
}
//对交易输入进行签名
func (bc *Blockchain) SignTransaction(tx *transaction.Transaction,privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]transaction.Transaction)
	for _,vin :=range tx.Vin {
		prevTX,err := bc.FindTransaction(vin.Txid) //找到输入引用的输出所在的交易
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey,prevTXs)
}
 
//验证交易
func (bc *Blockchain) VerifyTransaction(tx *transaction.Transaction) bool {
	prevTXs := make(map[string]transaction.Transaction)
 
	for _, vin := range tx.Vin {
		prevTX,err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs) //验证签名
}