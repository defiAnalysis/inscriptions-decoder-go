package main

import (
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
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
	//connCfg := &rpcclient.ConnConfig{
	//	Host:         "localhost:8332",
	//	User:         "coreincp",
	//	Pass:         "oHGkzaGPRPcWX3xz",
	//	HTTPPostMode: true,
	//	DisableTLS:   true,
	//}
	//
	//client, _ := rpcclient.New(connCfg, nil)
	//
	//blockHash, err := client.GetBlockHash(793980)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//block, err := client.GetBlock(blockHash)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 遍历每个交易并打印详细信息
	//for _, tx := range block.Transactions {
	//	txHash := tx.TxHash()
	//	fmt.Printf("TxID: %s\n", txHash.String())
	//
	//	// 输出交易输入信息
	//	for _, input := range tx.TxIn {
	//		fmt.Println("Input:")
	//		fmt.Printf("  TxID: %s\n", input.PreviousOutPoint.Hash.String())
	//		fmt.Printf("  Output Index: %d\n", input.PreviousOutPoint.Index)
	//	}
	//}
	if err := GetBlock1(793980); err != nil {
		fmt.Errorf("err:", err.Error())
		return
	}
}

func GetBlock1(height int64) error {
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "coreincp",
		Pass:         "oHGkzaGPRPcWX3xz",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, _ := rpcclient.New(connCfg, nil)

	hash, err := client.GetBlockHash(height)
	if err != nil {
		fmt.Errorf("GetBlock GetBlockHash err:", err.Error())
		return err
	}

	rawBlock, err := client.GetBlock(hash)
	if err != nil {
		fmt.Errorf("GetBlock GetBlockVerboseTx:", err.Error())
		return err
	}

	for key, rawTx := range rawBlock.Transactions {
		if key > 3 {
			return nil
		}
		fmt.Println("TxHash:", rawTx.TxHash().String())
		for id, value := range rawTx.TxIn {
			fmt.Println("input id:", id)
			fmt.Println("Hash:", value.PreviousOutPoint.Hash)
			fmt.Println("Index:", value.PreviousOutPoint.Index)

			key := fmt.Sprintf("%s%s%d", value.PreviousOutPoint.Hash.String(), "i", value.PreviousOutPoint.Index)
			fmt.Println("key====:", key)

			//if len(value.Witness) <= 1 {
			//	continue
			//}
			//
			//if len(value.Witness[1]) < 40 || !isInscribed(value.Witness[1]) {
			//	continue
			//}
			//
			//transaction := Transaction1{
			//	Hash:        value.PreviousOutPoint.Hash.String(),
			//	Index:       value.PreviousOutPoint.Index,
			//	Txinwitness: value.Witness[1],
			//}
			//
			//datatype, data, err := ExtractOrdFile(transaction.Txinwitness)
			//
			//if err != nil {
			//	continue
			//}
			////println("block", sl.Height(), "has tx", tx.Hash.String(), "len", string(typ), "-", len(data), "bytes")
			//if true {
			//	ext := datatype
			//	tps := strings.SplitN(string(datatype), "/", 2)
			//	if len(tps) == 2 {
			//		ext = tps[1]
			//	}
			//
			//	fmt.Println("Hash: %s,type:%s,data:%s", transaction.Hash, transaction.Index, ext, string(data))
			//}
		}

		fmt.Println("==============================")
	}

	return nil
}
