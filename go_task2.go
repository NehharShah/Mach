package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	ethInfuraURL     = "https://sepolia.infura.io/v3/{project_id}"          // add your project id
	l2InfuraURL      = "https://arbitrum-sepolia.infura.io/v3/{project_id}" // add your project id
	ethPrivateKey    = ""                                                   // add your private key
	l2PrivateKey     = ""                                                   // add your private key
	uniswapRouter    = "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	curvePool        = "0x654273fbe9445549D4B3817bb50375717894C49B"
	usdcTokenAddress = "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238"
	usdtTokenAddress = "0x30fA2FbE15c1EaDfbEF28C188b7B8dbd3c1Ff2eB"
	amountIn         = 1000000 // 1 USDC (6 decimal places)
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Perform Uniswap trade
	go func() {
		defer wg.Done()
		err := tradeOnUniswap(usdcTokenAddress, usdtTokenAddress, amountIn)
		if err != nil {
			log.Printf("Uniswap trade error: %v", err)
			return
		}
		log.Println("Uniswap trade successful")
	}()

	// Perform Curve trade
	go func() {
		defer wg.Done()
		err := tradeOnCurve(usdtTokenAddress, usdcTokenAddress, amountIn)
		if err != nil {
			log.Printf("Curve trade error: %v", err)
			return
		}
		log.Println("Curve trade successful")
	}()

	wg.Wait()
	log.Println("All trades attempted")
}

func getDynamicGasPrice(client *ethclient.Client) (*big.Int, error) {
	ctx := context.Background()
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}
	return gasPrice, nil
}

func tradeOnUniswap(fromToken, toToken string, amount int) error {
	client, err := ethclient.Dial(ethInfuraURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(ethPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	if err != nil {
		return fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = uint64(300000)

	gasPrice, err := getDynamicGasPrice(client)
	if err != nil {
		return fmt.Errorf("failed to get dynamic gas price: %v", err)
	}
	auth.GasPrice = gasPrice

	// Simulated Uniswap trade logic
	fmt.Printf("Simulated trading %d of %s to %s on Uniswap\n", amount, fromToken, toToken)

	// Simulate sending a transaction and getting a transaction hash
	tx := types.NewTransaction(0, common.HexToAddress(uniswapRouter), big.NewInt(0), auth.GasLimit, auth.GasPrice, nil)
	txHash := tx.Hash().Hex()
	fmt.Printf("Transaction hash: %s\n", txHash)

	return nil
}

func tradeOnCurve(fromToken, toToken string, amount int) error {
	client, err := ethclient.Dial(l2InfuraURL)
	if err != nil {
		return fmt.Errorf("failed to connect to L2 client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(l2PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(421613)) // Arbitrum chain ID
	if err != nil {
		return fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)

	gasPrice, err := getDynamicGasPrice(client)
	if err != nil {
		return fmt.Errorf("failed to get dynamic gas price: %v", err)
	}
	auth.GasPrice = gasPrice

	// Simulated Curve trade logic
	fmt.Printf("Simulated trading %d of %s to %s on Curve\n", amount, fromToken, toToken)

	// Simulate sending a transaction and getting a transaction hash
	tx := types.NewTransaction(0, common.HexToAddress(curvePool), big.NewInt(0), auth.GasLimit, auth.GasPrice, nil)
	txHash := tx.Hash().Hex()
	fmt.Printf("Transaction hash: %s\n", txHash)

	return nil
}
