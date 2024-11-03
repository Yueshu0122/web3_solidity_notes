package main

//
//import (
//	"bytes"
//	"crypto/sha256"
//	"encoding/hex"
//	"encoding/json"
//	"io"
//	"log"
//	"net/http"
//	"os"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//
//	"github.com/davecgh/go-spew/spew"
//	"github.com/gorilla/mux"
//	"github.com/joho/godotenv"
//)
//
//type Block struct {
//	Index     int
//	Timestamp string
//	BPM       int
//	Hash      string
//	PrevHash  string
//	Nonce     int
//}
//
//var Blockchain []Block
//
//type Message struct {
//	BPM int
//}
//
//var mutex = &sync.Mutex{}
//
//const difficulty = 4 // 难度值，表示哈希值前缀需要多少个0
//
//
//
//func main() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	go func() {
//		t := time.Now()
//		genesisBlock := Block{}
//		genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "",0}
//
//		spew.Dump(genesisBlock)
//
//		mutex.Lock()
//		Blockchain = append(Blockchain, genesisBlock)
//		mutex.Unlock()
//	}()
//
//	go func() {
//		nodes := []string{"http://node1.example.com", "http://node2.example.com"}
//		for {
//			syncBlockchain(nodes)
//			time.Sleep(10 * time.Second) // 每10秒同步一次
//		}
//	}()
//
//	log.Fatal(run())
//}
//
//// run 函数：
//// 创建路由处理器。
//// 获取环境变量中的端口号。
//// 创建HTTP服务器，设置地址、处理器、读写超时时间和最大头部字节数。
//// 启动HTTP服务器，如果启动失败则返回错误。
//func run() error {
//	mux := makeMuxRouter()
//	httpPort := os.Getenv("PORT")
//	log.Println("HTTP Server Listening on port :", httpPort)
//	s := &http.Server{
//		Addr:           ":" + httpPort,
//		Handler:        mux,
//		ReadTimeout:    10 * time.Second,
//		WriteTimeout:   10 * time.Second,
//		MaxHeaderBytes: 1 << 20,
//	}
//
//	if err := s.ListenAndServe(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// makeMuxRouter 函数：
//// 创建一个新的路由处理器。
//// 设置GET请求的处理函数为handleGetBlockchain。
//// 设置POST请求的处理函数为handleWriteBlock。
//// 返回路由处理器。
//func makeMuxRouter() http.Handler {
//	muxRouter := mux.NewRouter()
//	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
//	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
//	muxRouter.HandleFunc("/block", handleBroadcastBlock).Methods("POST")
//	return muxRouter
//}
//
//// handleGetBlockchain 函数：
//// 将区块链数据转换为JSON格式。
//// 如果转换失败，则返回500错误。
//// 将JSON数据写入响应体。
//func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
//	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	io.WriteString(w, string(bytes))
//}
//
//func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	var msg Message
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&msg); err != nil {
//		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
//		return
//	}
//
//	defer r.Body.Close()
//
//	mutex.Lock()
//	prevBlock := Blockchain[len(Blockchain)-1]
//	newBlock := generateBlock(prevBlock, msg.BPM)
//
//	if isBlockValid(newBlock, prevBlock) {
//		Blockchain = append(Blockchain, newBlock)
//		spew.Dump(Blockchain)
//
//		// 广播新区块给其他节点
//		broadcastNewBlock(newBlock)
//	}
//
//	mutex.Unlock()
//
//	respondWithJSON(w, r, http.StatusCreated, newBlock)
//}
//
//func broadcastNewBlock(block Block) {
//	// 这里假设有一个节点列表，遍历每个节点并发送区块
//	nodes := []string{"http://node1.example.com", "http://node2.example.com"}
//	for _, node := range nodes {
//		resp, err := http.Post(node+"/block", "application/json", bytes.NewBuffer(mustMarshal(block)))
//		if err != nil {
//			log.Printf("Failed to broadcast block to %s: %v", node, err)
//			continue
//		}
//		defer resp.Body.Close()
//		if resp.StatusCode != http.StatusOK {
//			log.Printf("Node %s rejected the block: %v", node, resp.Status)
//		}
//	}
//}
//
//func mustMarshal(v interface{}) []byte {
//	data, err := json.Marshal(v)
//	if err != nil {
//		log.Fatalf("Failed to marshal data: %v", err)
//	}
//	return data
//}
//
//func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
//	response, err := json.MarshalIndent(payload, "", "  ")
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("HTTP 500: Internal Server Error"))
//		return
//	}
//	w.WriteHeader(code)
//	w.Write(response)
//}
//
//func isBlockValid(newBlock, oldBlock Block) bool {
//	if oldBlock.Index+1 != newBlock.Index {
//		return false
//	}
//
//	if oldBlock.Hash != newBlock.PrevHash {
//		return false
//	}
//
//	if !isValidProofOfWork(newBlock) {
//		return false
//	}
//
//	return true
//}
//
//func isValidProofOfWork(block Block) bool {
//	hash := calculateHashWithNonce(block, block.Nonce)
//	return hash[:difficulty] == strings.Repeat("0", difficulty)
//}
//
//func handleBroadcastBlock(w http.ResponseWriter, r *http.Request) {
//	var newBlock Block
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&newBlock); err != nil {
//		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"error": "Invalid block format"})
//		return
//	}
//
//	defer r.Body.Close()
//
//	mutex.Lock()
//	defer mutex.Unlock()
//
//	if len(Blockchain) == 0 {
//		// 如果本地区块链为空，直接添加新区块
//		Blockchain = append(Blockchain, newBlock)
//		respondWithJSON(w, r, http.StatusOK, map[string]string{"message": "Block added to empty blockchain"})
//		return
//	}
//
//	prevBlock := Blockchain[len(Blockchain)-1]
//
//	if isBlockValid(newBlock, prevBlock) {
//		Blockchain = append(Blockchain, newBlock)
//		respondWithJSON(w, r, http.StatusOK, map[string]string{"message": "Block added to blockchain"})
//	} else {
//		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"error": "Invalid block"})
//	}
//}
//
//
//// calculateHash 函数：
//// 将区块的索引、时间戳、BPM和前一个哈希值拼接成字符串。
//// 使用SHA-256算法计算哈希值。
//// 将哈希值转换为十六进制字符串并返回。
//func calculateHash(block Block) string {
//	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash
//	h := sha256.New()
//	h.Write([]byte(record))
//	hashed := h.Sum(nil)
//	return hex.EncodeToString(hashed)
//}
//
//func proofOfWork(block Block) (int, string) {
//	nonce := 0
//	for {
//		nonce++
//		hash := calculateHashWithNonce(block, nonce)
//		if hash[:difficulty] == strings.Repeat("0", difficulty) {
//			return nonce, hash
//		}
//	}
//}
//
//func calculateHashWithNonce(block Block, nonce int) string {
//	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + strconv.Itoa(nonce)
//	h := sha256.New()
//	h.Write([]byte(record))
//	hashed := h.Sum(nil)
//	return hex.EncodeToString(hashed)
//}
//
//func generateBlock(oldBlock Block, BPM int) Block {
//	var newBlock Block
//
//	t := time.Now()
//
//	newBlock.Index = oldBlock.Index + 1
//	newBlock.Timestamp = t.String()
//	newBlock.BPM = BPM
//	newBlock.PrevHash = oldBlock.Hash
//	nonce, hash := proofOfWork(newBlock)
//	newBlock.Nonce = nonce
//	newBlock.Hash = hash
//
//	return newBlock
//}
//
//func syncBlockchain(nodes []string) {
//	longestChain := Blockchain
//	longestLength := len(Blockchain)
//
//	for _, node := range nodes {
//		resp, err := http.Get(node + "/")
//		if err != nil {
//			log.Printf("Failed to fetch blockchain from %s: %v", node, err)
//			continue
//		}
//		defer resp.Body.Close()
//
//		var remoteBlockchain []Block
//		if err := json.NewDecoder(resp.Body).Decode(&remoteBlockchain); err != nil {
//			log.Printf("Failed to decode blockchain from %s: %v", node, err)
//			continue
//		}
//
//		if len(remoteBlockchain) > longestLength && isChainValid(remoteBlockchain) {
//			longestChain = remoteBlockchain
//			longestLength = len(remoteBlockchain)
//		}
//	}
//
//	if  !slicesEqual(longestChain, Blockchain)  {
//		mutex.Lock()
//		Blockchain = longestChain
//		mutex.Unlock()
//		log.Println("Blockchain updated to the longest chain")
//	}
//}
//
//
//func isChainValid(chain []Block) bool {
//	for i := 1; i < len(chain); i++ {
//		if !isBlockValid(chain[i], chain[i-1]) {
//			return false
//		}
//	}
//	return true
//}
//
//func slicesEqual(a, b []Block) bool {
//	if len(a) != len(b) {
//		return false
//	}
//	for i := range a {
//		if a[i] != b[i] {
//			return false
//		}
//	}
//	return true
//}
