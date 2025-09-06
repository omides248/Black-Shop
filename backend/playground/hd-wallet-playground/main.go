package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// https://cloud.google.com/application/web3/faucet/ethereum/sepolia
const infuraUrl = "https://sepolia.infura.io/v3/fea54eb8d2df44fa824cbdfe6f10cb56"           // RPC Request (API for connect to ethereum node)
const infuraWebsocketUrl = "wss://sepolia.infura.io/ws/v3/fea54eb8d2df44fa824cbdfe6f10cb56" // RPC Request (API for connect to ethereum node)
const mnemonic = "future guard belt volume list slim final where call topple vote brush"

func main() {
	//TestBaseLogic()
	//TestPayment()
	TestCheckPaymentSub()
}

func TestCheckPaymentSub() {
	wallet, _ := CreateWallet(mnemonic)
	derivationPath := "m/44'/60'/0'/0/100"
	paymentAddress, _, _ := GenerateDerivation(derivationPath, wallet)

	fmt.Printf("Payment address created: %s\n", paymentAddress)
	fmt.Println("------------------------------------------------------------------")

	fmt.Println("Listener new block")

	SubscribeToNewBlocks(paymentAddress)
}

func TestPayment() {

	// Mnemonic and Wallet
	//mnemonic, _ := GenerateMnemonic()
	fmt.Println("mnemonic:", mnemonic)
	wallet, _ := CreateWallet(mnemonic)
	fmt.Println("Successfully created HD wallet on mnemonic:", mnemonic)

	// Create payment address
	derivationPath := "m/44'/60'/0'/0/100"
	paymentAddress, _, _ := GenerateDerivation(derivationPath, wallet)
	fmt.Println("Successfully created HD payment address:", paymentAddress)
	fmt.Println("0.000001 ETH")

	CheckPayment(paymentAddress)
}

func TestBaseLogic() {
	fmt.Printf("=========================================================%s\n", strings.Repeat("=", 65))
	mnemonic, _ := GenerateMnemonic() // Save secure place
	wordCount := len(strings.Split(strings.TrimSpace(mnemonic), " "))
	fmt.Printf("mnemonic (%d): %s\n", wordCount, mnemonic)
	fmt.Printf("=========================================================%s\n", strings.Repeat("=", 65))

	//masterKey, _ := GenerateMasterKey(mnemonic)
	//fmt.Printf("masterKey: %s\n", masterKey)
	//fmt.Printf("=========================================================%s\n", strings.Repeat("=", 65))

	wallet, _ := CreateWallet(mnemonic)
	fmt.Printf("wallet: %v\n", wallet)
	fmt.Printf("=========================================================%s\n", strings.Repeat("=", 65))

	pathString1 := "m/44'/60'/0'/0/0"
	address1, pathStr1, _ := GenerateDerivation(pathString1, wallet)
	fmt.Printf("Etherume address1: %s\n", address1)
	fmt.Printf("Etherume pathString1: %s\n", pathStr1)
	fmt.Printf("=========================================================%s\n", strings.Repeat("=", 65))

	pathString2 := "m/44'/60'/0'/0/1"
	address2, pathStr2, _ := GenerateDerivation(pathString2, wallet)
	fmt.Printf("Etherume address2: %s\n", address2)
	fmt.Printf("Etherume pathString2: %s\n", pathStr2)
	fmt.Printf("=========================================================%s\n", strings.Repeat("=", 65))
}

func GenerateMnemonic() (string, error) {
	// bitSize: 128 -> 12 words
	// bitSize: 256 -> 24 words
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

//func GenerateMasterKey(mnemonic string) (string, error) {
//	seed := bip39.NewSeed(mnemonic, "")
//
//	// Create masterKey
//	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
//	if err != nil {
//		return "", err
//	}
//
//	return masterKey.String(), nil
//}

func CreateWallet(mnemonic string) (*hdwallet.Wallet, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func GenerateDerivation(pathString string, wallet *hdwallet.Wallet) (string, string, error) {
	path := hdwallet.MustParseDerivationPath(pathString)

	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", "", err
	}
	address := account.Address.Hex()

	return address, path.String(), nil
}

func CheckPayment(paymentAddress string) {
	for {
		balance, err := GetAddressBalance(paymentAddress)
		if err != nil {
			log.Printf("Error getting balance: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// 1 Ether = 1,000,000,000,000,000,000 Wei (10^18)
		balanceInEth := new(big.Float).SetInt(balance)
		balanceInEth.Quo(balanceInEth, big.NewFloat(1e18))

		fmt.Printf("Current balance: %f ETH\n", balanceInEth)

		if balance.Cmp(big.NewInt(0)) > 0 {
			fmt.Println("Payment successfully identity")
			fmt.Printf("Received price: %f ETH\n", balanceInEth)
			// Change payment status in database
			break
		}

		time.Sleep(10 * time.Second)
	}
}

func GetAddressBalance(address string) (*big.Int, error) {
	client, err := ethclient.Dial(infuraUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Infura: %w", err)
	}
	defer client.Close()

	// Convert string address to Ethereum address
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %w", err)
	}

	return balance, nil
}

func SubscribeToNewBlocks(paymentAddress string) {
	client, err := ethclient.DialContext(context.Background(), infuraWebsocketUrl)
	if err != nil {
		log.Fatalf("failed to connecto to infuraUrl websocket %v", err)
	}
	defer client.Close()

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("failed to subscribe to new blocks %v", err)
	}

	targetAddress := common.HexToAddress(paymentAddress)

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("failed to subscribe to new blocks %v", err)
		case header := <-headers:
			fmt.Printf("Received block header: %v\n", header)
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Printf("failed to get detail block %v", err)
				continue
			}

			for _, tx := range block.Transactions() {
				if tx.To() != nil && tx.To().Hex() == targetAddress.Hex() {
					fmt.Println("\nPaid transaction find")

					chainID, err := client.NetworkID(context.Background())
					if err != nil {
						log.Printf("Error recieved Chain ID: %v", err)
						continue
					}

					from, err := types.Sender(types.LatestSignerForChainID(chainID), tx)
					if err != nil {
						log.Printf("Error to get sender address: %v", err)
						continue
					}

					fmt.Printf("   - Transaction Hash: %s\n", tx.Hash().Hex())
					fmt.Printf("   - From: %s\n", from.Hex())
					fmt.Printf("   - To: %s\n", tx.To().Hex())
					fmt.Printf("   - Amount: %s Wei\n", tx.Value().String())

					return
				}
			}
		}
	}

}
