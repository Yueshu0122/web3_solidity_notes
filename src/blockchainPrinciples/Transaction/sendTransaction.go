package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var nodeList = []string{
	"http://localhost:8001",
	"http://localhost:8002",
	"http://localhost:8003",
	"http://localhost:8004",
	"http://localhost:8005",
}

type Transaction struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Amount    int    `json:"amount"`
	Timestamp string `json:"timestamp"`
}

func main() {
	// 节点列表
	nodes := nodeList
	// 随机选择一个节点
	rand.Seed(time.Now().UnixNano())
	randomNode := nodes[rand.Intn(len(nodes))]

	// 创建一笔交易
	transaction := Transaction{
		Sender:    "Alice",
		Receiver:  "Bob",
		Amount:    100,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 将交易转换为 JSON 格式
	jsonData, err := json.Marshal(transaction)
	if err != nil {
		fmt.Println("Error marshaling transaction:", err)
		os.Exit(1)
	}

	// 发送 POST 请求到随机选择的节点
	url := fmt.Sprintf("%s/transactions/new", randomNode)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending transaction:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	// 打印响应
	fmt.Println("Response from node:", string(body))
}
