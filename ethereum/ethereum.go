package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mindsgn-studio/pocket-money-go/database"
	"github.com/mindsgn-studio/pocket-money-go/logs"
)

type networkDetails struct {
	Name       string   `json:"name"`
	ChainID    int      `json:"chainID"`
	ChainIDHex string   `json:"ChainIDHex"`
	Currency   string   `json:"currency"`
	Mainnet    bool     `json:"mainnet"`
	RPC        []string `json:"rpc"`
}

type Wallets struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Balance     int    `json:"balance"`
	Network     int    `json:"network"`
	FiatBalance int    `json:fiatBalance`
}

type TotalBalance struct {
	Wallets          []Wallets `json:"json"`
	TotalFiatBalance int       `json:"json"`
}

var NetworkMainnetList []string = []string{
	"polygon-mainnet",
}

var NetworkTestnetList []string = []string{
	"polygon-mumbai",
}

func GetTotalBalance(password string, network string) int {
	total := 0

	if network == "mainnet" {
		for _, networkDetails := range NetworkMainnetList {
			fmt.Println(networkDetails)
		}
	} else {
		for _, networkDetails := range NetworkTestnetList {
			details := GetNetwork(networkDetails)
			list := details.RPC[0]

			client, err := ethclient.Dial(list)
			if err != nil {
				log.Fatal(err)
			}

			wallets := database.GetWallets(password)

			for _, wallet := range wallets {
				account := common.HexToAddress(wallet.Address)
				balance, err := client.BalanceAt(context.Background(), account, nil)
				if err != nil {
					log.Fatal(err)
				}

				fbalance := new(big.Float)
				fbalance.SetString(balance.String())
				ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
				fmt.Println(ethValue, details.Currency)
			}
		}
	}

	return total
}

func CreateNewEthereumWallet(password string) bool {
	newPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(newPrivateKey)
	publicKey := newPrivateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	privateKey := hexutil.Encode(privateKeyBytes)

	directory, err := database.GetDataDirectory()
	if err != nil {
		logs.LogError(err.Error())
		return false
	}

	db, err := database.OpenDatabase(directory, password)
	if err != nil {
		logs.LogError(err.Error())
		return false
	}

	database.InsertWallet(db, "ethereum", privateKey, address)

	return false
}

func GetNetwork(network string) networkDetails {
	switch network {
	case "polygon-mainnet":
		rpcList := []string{
			"wss://polygon-bor-rpc.publicnode.com",
			"https://polygon.llamarpc.com",
			"wss://polygon.drpc.org",
		}

		return networkDetails{
			Name:       "polygon",
			ChainID:    137,
			ChainIDHex: "0x89",
			Currency:   "matic",
			Mainnet:    true,
			RPC:        rpcList,
		}
	case "polygon-mumbai":
		rpcList := []string{
			"https://polygon-mumbai.gateway.tenderly.co",
			"https://polygon-mumbai.api.onfinality.io/public",
			"https://gateway.tenderly.co/public/polygon-mumbai",
		}

		return networkDetails{
			Name:       "polygon",
			ChainID:    137,
			ChainIDHex: "0x89",
			Currency:   "matic",
			Mainnet:    true,
			RPC:        rpcList,
		}

	default:
		return networkDetails{
			Name:     "",
			ChainID:  0,
			Currency: "",
			Mainnet:  false,
		}
	}
}
