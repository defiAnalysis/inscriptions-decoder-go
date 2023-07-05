package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"inscription-decoder/util"
	"strings"
)

type Transaction struct {
	Hash        string
	Index       uint32
	Txinwitness []byte
}

func GetBlock(height int64) error {
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

	for _, rawTx := range rawBlock.Transactions {
		for _, value := range rawTx.TxIn {
			if len(value.Witness[1]) < 40 || !isInscribed(value.Witness[1]) {
				continue
			}

			transaction := Transaction{
				Hash:        value.PreviousOutPoint.Hash.String(),
				Index:       value.PreviousOutPoint.Index,
				Txinwitness: value.Witness[1],
			}
			fmt.Println("Hash : %s,Index:%d,Txinwitness:%s\n", transaction.Hash, transaction.Index, transaction.Txinwitness)
			datatype, data, err := ExtractOrdFile(transaction.Txinwitness)

			if err != nil {
				println(transaction.Hash, err.Error())
				continue
			}
			//println("block", sl.Height(), "has tx", tx.Hash.String(), "len", string(typ), "-", len(data), "bytes")
			if true {
				ext := datatype
				tps := strings.SplitN(string(datatype), "/", 2)
				if len(tps) == 2 {
					ext = tps[1]
				}

				fmt.Println("Hash: %s,type:%s,data:%s", transaction.Hash, ext, data)
			}
		}
	}

	return nil
}

func ExtractOrdFile(p []byte) (typ string, data []byte, e error) {
	var opcode_idx int
	var byte_idx int

	for byte_idx < len(p) {
		opcode, vchPushValue, n, er := util.GetOpcode(p[byte_idx:])
		if er != nil {
			e = errors.New("ExtractOrdinaryFile: " + er.Error())
			return
		}

		byte_idx += n

		switch opcode_idx {
		case 0:
			if len(vchPushValue) != 32 {
				e = errors.New("opcode_idx 0: No push data 32 bytes")
				return
			}
		case 1:
			if opcode != util.OP_CHECKSIG {
				e = errors.New("opcode_idx 1: OP_CHECKSIG missing")
				return
			}
		case 2:
			if opcode != util.OP_FALSE {
				e = errors.New("opcode_idx 2: OP_FALSE missing")
				return
			}
		case 3:
			if opcode != util.OP_IF {
				e = errors.New("opcode_idx 3: OP_IF missing")
				return
			}
		case 4:
			if len(vchPushValue) != 3 || string(vchPushValue) != "ord" {
				e = errors.New("opcode_idx 4: missing ord string")
				return
			}
		case 5:
			if len(vchPushValue) != 1 || vchPushValue[0] != 1 {
				//println("opcode_idx 5:", hex.EncodeToString(vchPushValue), string(vchPushValue), "-ignore")
				opcode_idx-- // ignore this one
			}
		case 6:
			typ = string(vchPushValue)
		case 7:
			if opcode != util.OP_FALSE {
				e = errors.New("opcode_idx 7: OP_FALSE missing")
				return
			}
		default:
			if opcode == util.OP_ENDIF {
				return
			}
			data = append(data, vchPushValue...)
		}

		opcode_idx++
	}
	return
}

func isInscribed(s []byte) bool {
	isncPattern, _ := hex.DecodeString("0063036f7264")
	return bytes.Contains(s, isncPattern)
}

func main() {
	if err := GetBlock(793980); err != nil {
		fmt.Println("err:", err.Error())
		return
	}

	fmt.Println("end======================")
}
