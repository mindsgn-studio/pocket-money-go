package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mindsgn-studio/pocket-money-go/database"
	"github.com/mindsgn-studio/pocket-money-go/logs"
)

type Wallet struct {
	Address      string  `json:"address"`
	Blockchain   string  `json:"blockchain"`
	BlockchainId string  `json:"blockchainId"`
	Decimals     uint    `json:"decimals"`
	Currency     string  `json:"currency"`
	FiatBalance  float64 `json:"fiatBalances"`
}

type Wallets struct {
	TotalFiat float64  `json:"totalFiat"`
	Currency  string   `json:"currency"`
	Wallets   []Wallet `json:"wallets"`
}

type Contract struct {
	Address      string `json:"address"`
	Blockchain   string `json:"blockchain"`
	BlockchainId string `json:"blockchainId"`
	Decimals     uint   `json:"decimals"`
}

type MarketData struct {
	Data struct {
		MarketCap         float64    `json:"market_cap"`
		MarketCapDiluted  float64    `json:"market_cap_diluted"`
		Liquidity         float64    `json:"liquidity"`
		Price             float64    `json:"price"`
		OffChainVolume    float64    `json:"off_chain_volume"`
		Volume            float64    `json:"volume"`
		VolumeChange24h   float64    `json:"volume_change_24h"`
		Volume7d          float64    `json:"volume_7d"`
		IsListed          bool       `json:"is_listed"`
		PriceChange24h    float64    `json:"price_change_24h"`
		PriceChange1h     float64    `json:"price_change_1h"`
		PriceChange7d     float64    `json:"price_change_7d"`
		PriceChange1m     float64    `json:"price_change_1m"`
		PriceChange1y     float64    `json:"price_change_1y"`
		Ath               float64    `json:"ath"`
		Atl               float64    `json:"atl"`
		Name              string     `json:"name"`
		Symbol            string     `json:"symbol"`
		Logo              string     `json:"logo"`
		Rank              int        `json:"rank"`
		Contracts         []Contract `json:"contracts"`
		TotalSupply       string     `json:"total_supply"`
		CirculatingSupply string     `json:"circulating_supply"`
	} `json:"data"`
}

type networkDetails struct {
	Name       string   `json:"name"`
	ChainID    int      `json:"chainID"`
	ChainIDHex string   `json:"ChainIDHex"`
	Currency   string   `json:"currency"`
	Mainnet    bool     `json:"mainnet"`
	RPC        []string `json:"rpc"`
}

var NetworkMainnetList []string = []string{
	"polygon-mainnet",
}

var NetworkTestnetList []string = []string{
	"polygon-mumbai",
}

func ConvertBody(body []byte) (MarketData, error) {
	var data MarketData
	err := json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetTotalBalance(password string, network string) (Wallets, error) {
	total := float64(0)
	var userWallet Wallets

	if network == "mainnet" {
		for _, networkDetails := range NetworkMainnetList {
			fmt.Println(networkDetails)
		}
	} else {
		for _, networkDetails := range NetworkTestnetList {
			details := GetNetwork(networkDetails)
			list := details.RPC[1]

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
				data, err := GetData(details.Name)
				if err != nil {
					log.Fatal(err)
				}

				price := ethValue.String()
				cryptoBalance, err := strconv.ParseFloat(price, 64)
				if err != nil {
					log.Fatal(err)
				}

				total += data.Data.Price * cryptoBalance

				walletData := Wallet{
					Address:      wallet.Address,
					Blockchain:   wallet.Name,
					BlockchainId: wallet.Name,
					Decimals:     18,
					Currency:     "USD",
					FiatBalance:  cryptoBalance * data.Data.Price,
				}

				userWallet.Wallets = append(userWallet.Wallets, walletData)
			}
		}
	}

	userWallet.TotalFiat = total
	userWallet.Currency = "USD"

	return userWallet, nil
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
			Mainnet:    false,
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

func GetData(name string) (MarketData, error) {
	url := "https://api.mobula.io/api/1/market/data?asset=" + name
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	data, err := ConvertBody(body)
	if err != nil {
		return MarketData{}, err
	}

	return data, nil
}
