package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var Blockchain []Block
var Mempool []Transaction
var mutex = &sync.Mutex{}

func main() {
	fmt.Println("Starting Blockchain Simulation...")

	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 通过命令行参数指定端口号
	portPtr := flag.Int("port", 8080, "HTTP server port")
	flag.Parse()
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = strconv.Itoa(*portPtr)
	}

	go initializeGenesisBlock()
	go syncBlockchainPeriodically()
	go generateBlocksPeriodically()

	log.Fatal(run())
}

func initializeGenesisBlock() {
	t := time.Now()
	genesisBlock := NewBlock(0, t, 0, "", []Transaction{})
	mutex.Lock()
	Blockchain = append(Blockchain, *genesisBlock)
	mutex.Unlock()
}

func syncBlockchainPeriodically() {
	nodes := nodeList
	for {
		syncBlockchain(nodes)
		time.Sleep(10 * time.Second) // 每10秒同步一次
	}
}

func generateBlocksPeriodically() {
	for {
		time.Sleep(30 * time.Second) // 每30秒生成一个新区块
		mutex.Lock()
		if len(Mempool) > 0 {
			prevBlock := Blockchain[len(Blockchain)-1]
			newBlock := generateBlock(prevBlock, 0, Mempool)
			if isBlockValid(newBlock, prevBlock) {
				Blockchain = append(Blockchain, newBlock)
				Mempool = []Transaction{} // 清空内存池
				broadcastNewBlock(newBlock)
			}
		}
		mutex.Unlock()
	}
}

func run() error {
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
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
