package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Pocket Money",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

/*
package main

import (
	"fmt"

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

func GetTotalBalance(password string) {
	wallets, err := ethereum.GetTotalBalance(password, "testnet")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(wallets)
}

func main() {
	InitialiseWallet("123456789")
	GetTotalBalance("123456789")
}*/
