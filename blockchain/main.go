package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Block struct {
	Index     int    //数据记录在区块中的位置
	Timestamp string //数据写入时间
	BPM       int    //插入区块的频率
	Hash      string //此条数据的SHA256标识符
	PrevHash  string //上一个区块的标识符
}

type Message struct {
	BPM int
}

var Blockchain []Block

func main() {
	//读取.env中的变量
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		//给定一个初始区块
		genesisBlock := Block{0, t.String(), 0, "", ""}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()
	log.Fatal(run())
}

//用于生成区块的SHA256哈希值
func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//新区块的生成
func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

//验证区块是否被篡改
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

//两个区块竞争，长的链优先
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

//设置路由，便于浏览器访问。
func makMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

//查看区块数据
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	//以json格式写会区块的数据
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

//给区块插入数据
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJson(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	if err != nil {
		respondWithJson(w, r, http.StatusInternalServerError, m)
		return
	}
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		//将结果打印在控制台
		spew.Dump(Blockchain)
	}
	respondWithJson(w, r, http.StatusCreated, newBlock)
}

//判断POST请求是否失败
func respondWithJson(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500:Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

//运行区块。
func run() error {
	mux := makMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", os.Getenv("PORT"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
