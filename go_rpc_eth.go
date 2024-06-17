package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL = "wss://sepolia.infura.io/ws/v3/{project_id}" // add your project id
	fromKey   = ""                                           // add your private key
	toAddress = ""                                           // add the address to send the ETH to
)

func main() {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	fromPrivateKey, err := crypto.HexToECDSA(fromKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKey := fromPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	toAddress := common.HexToAddress(toAddress)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	log.Println("Listening for new blocks...")

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new head: %v", err)
	}

	blockCount := 0

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case header := <-headers:
			blockCount++
			log.Printf("New block: %d", header.Number.Int64())

			if blockCount%10 == 0 {
				log.Printf("Attempting to send 0.001 ETH from %s to %s", fromAddress.Hex(), toAddress.Hex())

				nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
				if err != nil {
					log.Printf("Failed to get nonce: %v", err)
					continue
				}

				gasPrice, err := client.SuggestGasPrice(context.Background())
				if err != nil {
					log.Printf("Failed to get gas price: %v", err)
					continue
				}

				value := big.NewInt(1000000000000000) // 0.001 ETH in wei
				gasLimit := uint64(21000)

				tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

				signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromPrivateKey)
				if err != nil {
					log.Printf("Failed to sign transaction: %v", err)
					continue
				}

				err = client.SendTransaction(context.Background(), signedTx)
				if err != nil {
					log.Printf("Failed to send transaction: %v", err)
					continue
				}

				log.Printf("Transaction sent: %s", signedTx.Hash().Hex())

				log.Println("Transaction completed after 10 blocks, exiting program.")
				return
			}
		}
	}
}
