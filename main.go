package main

import (
	"github.com/mindsgn-studio/pocket-money-go/database"
	"github.com/mindsgn-studio/pocket-money-go/ethereum"
	"github.com/mindsgn-studio/pocket-money-go/logs"
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
	wallets := database.GetWallets(password)
	return wallets
}

func main() {
	// initialize wallet
	InitialiseWallet("123456789")

	// create wallet
	GetWallets("123456789")
}
