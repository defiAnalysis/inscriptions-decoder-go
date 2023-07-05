package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Transaction struct {
	TxID    string   `json:"txid"`
	Value   int      `json:"value"`
	Address string   `json:"address"`
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

type Input struct {
	TxID     string   `json:"txid"`
	Output   int      `json:"output"`
	PKScript string   `json:"pkscript"`
	Sequence int      `json:"sequence"`
	Witness  []string `json:"witness"`
}

type Output struct {
	Value   int    `json:"value"`
	Address string `json:"address"`
}

type Block struct {
	Transactions []Transaction `json:"tx"`
}

func main() {
	// 替换为您的比特币节点的URL和区块哈希
	nodeURL := "http://localhost:8332"
	blockHash := "22e97ecac38499addbb1140c0a1500b9959f4ac3224791d2f6e28f7d55472d7b"

	// 发起HTTP请求获取区块数据
	resp, err := http.Get(nodeURL + "/rest/block/" + blockHash + ".json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 解析HTTP响应的JSON数据
	var block Block
	err = json.NewDecoder(resp.Body).Decode(&block)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("block:", block)

	// 遍历每个交易并打印详细信息
	for _, tx := range block.Transactions {
		fmt.Println("Transaction ID:", tx.TxID)
		fmt.Println("Value:", tx.Value)
		fmt.Println("Address:", tx.Address)

		fmt.Println("Inputs:")
		for _, input := range tx.Inputs {
			fmt.Println("  TxID:", input.TxID)
			fmt.Println("  Output:", input.Output)
			fmt.Println("  PKScript:", input.PKScript)
			fmt.Println("  Sequence:", input.Sequence)
			fmt.Println("  Witness:", input.Witness)
		}

		fmt.Println("Outputs:")
		for _, output := range tx.Outputs {
			fmt.Println("  Value:", output.Value)
			fmt.Println("  Address:", output.Address)
		}

		fmt.Println("----------------------------------")
	}
}
