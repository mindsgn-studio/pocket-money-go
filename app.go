package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mindsgn-studio/pocket-money-go/database"
	"github.com/mindsgn-studio/pocket-money-go/ethereum"
	"github.com/mindsgn-studio/pocket-money-go/logs"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) StartUp(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OnDomReady(ctx context.Context) {
	a.ctx = ctx
}

// start up
func (a *App) InitialiseWallet(password string) bool {
	initialised := database.InitialiseWallet(password)
	if initialised {
		logs.LogError("Wallet Initialised")
		return true
	} else {
		logs.LogError("Wallet Failed Initialised")
		return false
	}
}

// user create new wallet
func (a *App) CreateWallet(name string, password string) bool {
	switch name {
	case "ethereum":
		ethereum.CreateNewEthereumWallet(password)
		return true
	default:
		return false
	}
}

// get list of wallet and balance
func (a *App) GetWallets(password string) []database.Wallet {
	wallets := database.GetWallets(password)
	return wallets
}

func (a *App) WalletExists() (bool, error) {
	return false, errors.New("Wallet not found")
}

func (a *App) GetTotalBalance(password string) {
	wallets, err := ethereum.GetTotalBalance(password, "testnet")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(wallets)
}
