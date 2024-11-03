package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := &Block{}
		genesisBlock = NewBlock(0, t, 0, "")
		mutex.Lock()
		Blockchain = append(Blockchain, *genesisBlock)
		mutex.Unlock()
	}()

	go func() {
		nodes := []string{"http://node1.example.com", "http://node2.example.com"}
		for {
			syncBlockchain(nodes)
			time.Sleep(10 * time.Second) // 每10秒同步一次
		}
	}()

	log.Fatal(run())
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
