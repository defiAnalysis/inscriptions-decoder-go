package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"log"
	"os"
	"strings"
)

var PROTOCOL_ID = []byte{0x72, 0x6f, 0x6c, 0x6c}

var (
	pointer   int
	inputData string = "2024531bd7c0de19af5b3009033c130e757702f7d4b933aa848be3d825e090a041ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d38003a7b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2273617473222c22616d74223a22313030303030303030227d68"
)

func readBytes(raw []byte, n int) []byte {
	value := raw[pointer : pointer+n]
	pointer += n
	return value
}

func getInitialPosition(raw []byte) (int, error) {
	inscriptionMark := []byte{0x00, 0x63, 0x03, 0x6f, 0x72, 0x64}
	pos := strings.Index(string(raw), string(inscriptionMark))
	if pos == -1 {
		return 0, errors.New("No ordinal inscription found in transaction")
	}
	return pos + len(inscriptionMark), nil
}

func readContentType(raw []byte) (string, error) {
	OP_1 := byte(0x51)

	b := readBytes(raw, 1)[0]
	if b != OP_1 {
		if b != 0x01 || readBytes(raw, 1)[0] != 0x01 {
			return "", errors.New("Invalid byte sequence")
		}
	}

	size := int(readBytes(raw, 1)[0])
	contentType := readBytes(raw, size)
	return string(contentType), nil
}

func readPushdata(raw []byte, opcode byte) ([]byte, error) {
	intOpcode := int(opcode)

	if 0x01 <= intOpcode && intOpcode <= 0x4b {
		return readBytes(raw, intOpcode), nil
	}

	numBytes := 0
	switch intOpcode {
	case 0x4c:
		numBytes = 1
	case 0x4d:
		numBytes = 2
	case 0x4e:
		numBytes = 4
	default:
		return nil, fmt.Errorf("Invalid push opcode %x at position %d", intOpcode, pointer)
	}

	if pointer+numBytes > len(raw) {
		return nil, fmt.Errorf("Invalid data length at position %d", pointer)
	}

	sizeBytes := readBytes(raw, numBytes)
	var size int
	switch numBytes {
	case 1:
		size = int(sizeBytes[0])
	case 2:
		size = int(binary.LittleEndian.Uint16(sizeBytes))
	case 4:
		size = int(binary.LittleEndian.Uint32(sizeBytes))
	}

	if pointer+size > len(raw) {
		return nil, fmt.Errorf("Invalid data length at position %d", pointer)
	}

	return readBytes(raw, size), nil
}

//func writeDataUri(data []byte, contentType string) {
//	dataBase64 := base64.StdEncoding.EncodeToString(data)
//	fmt.Printf("data:%s;base64,%s", contentType, dataBase64)
//}

func writeFile(data []byte, filename string) {
	if filename == "" {
		filename = "out.txt"
	}

	filename = "out/" + filename

	i := 1
	baseFilename := filename
	for _, err := os.Stat(filename); !os.IsNotExist(err); _, err = os.Stat(filename) {
		i++
		filename = fmt.Sprintf("%s%d", baseFilename, i)
	}

	fmt.Printf("Writing contents to file \"%s\"\n", filename)
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main11() {
	//flag := IsTaprootAddress("512023b7432f5010b5fd9178639462e658757b2dcc4cb2dffcb4ab6c8976350d09ba")
	//"0014b5a33d86d07fba16e6d0ef3ced5e1d81c74a21e9"
	//golangAddress := "512023b7432f5010b5fd9178639462e658757b2dcc4cb2dffcb4ab6c8976350d09ba"

	//IsTaprootAddress(golangAddress)
}

func main2() {
	raw, err := hex.DecodeString(inputData)
	if err != nil {
		log.Fatal(err)
	}

	pointer, _ = getInitialPosition(raw)

	contentType, _ := readContentType(raw)
	fmt.Printf("Content type: %s\n", contentType)
	if readBytes(raw, 1)[0] != byte(0x00) {
		fmt.Println("Error: Invalid byte sequence")
		os.Exit(1)
	}

	data := []byte{}

	OP_ENDIF := byte(0x68)
	opcode := readBytes(raw, 1)[0]
	for opcode != OP_ENDIF {
		chunk, _ := readPushdata(raw, opcode)
		data = append(data, chunk...)
		opcode = readBytes(raw, 1)[0]
	}

	fmt.Println("data:", string(data))

	fmt.Printf("Total size: %d bytes\n", len(data))
	//writeFile(data, "output")
	fmt.Println("\nDone")
}

func main() {
	hash := "7c62902993ba36491a1c5bfe6547bbad9c21958717ed4aa2878c79c20cd87968"
	strHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		fmt.Println("NewHashFromStr err:", err.Error())
		return
	}
	datas, err := ReadTransaction(strHash)
	if err != nil {
		fmt.Println("ReadTransaction err:", err.Error())
		return
	}

	fmt.Println("datas:", string(datas))
}

func ReadTransaction(hash *chainhash.Hash) ([]byte, error) {
	fmt.Println("ReadTransaction:")

	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "coreincp",
		Pass:         "oHGkzaGPRPcWX3xz",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, _ := rpcclient.New(connCfg, nil)
	fmt.Println("hash:", hash.String())
	tx, err := client.GetRawTransaction(hash)
	if err != nil {
		fmt.Println("err:", err.Error())
		return nil, err
	}

	if len(tx.MsgTx().TxIn[0].Witness) > 1 {
		raw := tx.MsgTx().TxIn[0].Witness[1]

		pointer, _ = getInitialPosition(raw)

		contentType, _ := readContentType(raw)
		fmt.Printf("Content type: %s\n", contentType)
		if readBytes(raw, 1)[0] != byte(0x00) {
			fmt.Println("Error: Invalid byte sequence")
			os.Exit(1)
		}

		data := []byte{}

		OP_ENDIF := byte(0x68)
		opcode := readBytes(raw, 1)[0]
		for opcode != OP_ENDIF {
			chunk, _ := readPushdata(raw, opcode)
			data = append(data, chunk...)
			opcode = readBytes(raw, 1)[0]
		}

		fmt.Println("data:", string(data))

		return data, nil
	}

	//// 解析交易
	//txBytes, err := btcutil.DecodeHex(tx)
	//if err != nil {
	//	fmt.Println("无法解析交易数据:", err)
	//	return nil, nil
	//}
	//
	//msgTx := btcutil.NewTx(txBytes)

	// 遍历输出脚本以找到Taproot地址
	//for _, txOut := range tx.MsgTx().TxOut {
	//	pkScript := txOut.PkScript
	//	_, address, _, err := txscript.ExtractPkScriptAddrs(
	//		pkScript, &chaincfg.MainNetParams,
	//	)
	//
	//	fmt.Println("address:", address)
	//
	//	if err != nil {
	//		fmt.Println("err:", err.Error())
	//		return nil, nil
	//	}
	//
	//}

	//address := hex.EncodeToString(tx.MsgTx().TxOut[0].PkScript)
	//fmt.Println("to address:", address)
	//
	//value := tx.MsgTx().TxOut[0].Value
	//fmt.Println("value:", value)

	//Witness := tx.MsgTx().TxIn[0].Witness[1]
	//
	//fmt.Println("tx:", hex.EncodeToString(Witness))
	//
	//contentType, _ := readContentType(Witness)
	//fmt.Printf("Content type: %s\n", contentType)
	//
	//if len(tx.MsgTx().TxIn[0].Witness) > 1 {
	//	witness := tx.MsgTx().TxIn[0].Witness[1]
	//	pushData, err := ExtractPushData(0, witness)
	//	if err != nil {
	//		return nil, err
	//	}
	//	// skip PROTOCOL_ID
	//	if pushData != nil && bytes.HasPrefix(pushData, PROTOCOL_ID) {
	//		return pushData[4:], nil
	//	}
	//}

	return nil, nil
}
