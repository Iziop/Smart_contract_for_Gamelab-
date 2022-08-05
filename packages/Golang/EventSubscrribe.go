package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var contractAbi = `[{"inputs":[{"internalType":"string","name":"initMessage","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"oldStr","type":"string"},{"indexed":false,"internalType":"string","name":"newStr","type":"string"}],"name":"UpdatedMessages","type":"event"},{"inputs":[],"name":"message","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"newMessage","type":"string"}],"name":"update","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

func main() {
	client, err := ethclient.Dial("wss://rinkeby.infura.io/ws/v3/18311b06e26342bcb7a1faafa420d1bb")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x00000000006c3852cbef3e08e8df289169ede581")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		log.Fatal(err)
	}

	//logs := make(chan types.Log)
	//sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	//if err != nil {
	//	log.Fatal(err)
	//}

	type LogTransfer struct {
		From   common.Address
		To     common.Address
		Tokens *big.Int
	}
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	for _, vLog := range logs {
		//fmt.Println(vLog.BlockHash.Hex()) // 0x3404b8c050aa0aacd0223e91b5c32fee6400f357764771d0684fa7b3f448f1a8
		//fmt.Println(vLog.BlockNumber)     // 2394201
		//fmt.Println(vLog.TxHash.Hex())    // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6

		//event := struct {
		//	from    string
		//	to      string
		//	tokenId int
		//}{}
		err, _ := contractAbi.Unpack("Transfer", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer

			err, _ := contractAbi.Unpack("Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())
		}
	}
}
