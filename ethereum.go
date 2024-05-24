package main

import (
	"fmt"

	"github.com/mindsgn-studio/pocket-wallet-ethereum/database"
	"github.com/mindsgn-studio/pocket-wallet-ethereum/ethereum"
	"github.com/mindsgn-studio/pocket-wallet-ethereum/logs"
)

func WalletExists() bool {
	return false
}

func InitialiseWallet(password string) bool {
	initialised := database.InitialiseWallet(password)
	if initialised {
		logs.LogError("Wallet Initialised")
		return true
	} else {
		logs.LogError("Wallet Failed Initialised")
		return false
	}
}

func CreateWallet(name string, password string) bool {
	switch name {
	case "ethereum":
		ethereum.CreateNewEthereumWallet(password)
		return true
	default:
		return false
	}
}

func GetWallets(password string) []database.Wallet {
	responded := database.GetWallets(password)
	fmt.Println(responded)
	return responded
}

func main() {
	// initialize wallet
	InitialiseWallet("123456789")

	// create wallet
	// CreateWallet("ethereum", "123456789")

	// create wallet
	// CreateWallet("brown", "123456789")

	// create wallet
	GetWallets("123456789")
}
