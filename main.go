package main

import (
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"log"
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
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "coreincp",
		Pass:         "oHGkzaGPRPcWX3xz",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, _ := rpcclient.New(connCfg, nil)

	blockHash, err := client.GetBlockHash(793980)
	if err != nil {
		log.Fatal(err)
	}

	block, err := client.GetBlock(blockHash)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历每个交易并打印详细信息
	for _, tx := range block.Transactions {
		txHash := tx.TxHash()
		fmt.Printf("TxID: %s\n", txHash.String())

		// 输出交易输入信息
		for _, input := range tx.TxIn {
			fmt.Println("Input:")
			fmt.Printf("  TxID: %s\n", input.PreviousOutPoint.Hash.String())
			fmt.Printf("  Output Index: %d\n", input.PreviousOutPoint.Index)
		}
	}
}
