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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mindsgn-studio/pocket-money-go/database"
	"github.com/mindsgn-studio/pocket-money-go/logs"
)

type Wallet struct {
	Address       string  `json:"address"`
	Blockchain    string  `json:"blockchain"`
	BlockchainId  string  `json:"blockchainId"`
	Decimals      uint    `json:"decimals"`
	Currency      string  `json:"currency"`
	FiatBalance   float64 `json:"fiatBalances"`
	CryptoBalance float64 `json:"cryptoBalances"`
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
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	ChainID  int    `json:"chainID"`
	Currency string `json:"currency"`
	Mainnet  bool   `json:"mainnet"`
	RPC      string `json:"rpc"`
	Explorer string `json:"explorer"`
}

var NetworkMainnetList []string = []string{
	"polygon-mainnet",
	// "polygon-zkevm",
	"Starknet-mainnet",
	"ethereum",
	"optimism",
	"arbitrum-one",
	"astar",
	"binance-smart-chain",
}

var NetworkTestnetList []string = []string{
	"polygon-mumbai",
	"ethereum-sepolia",
	"ethereum-goerli",
	// "polygon-zkevm-testnet",
	"polygon-mainnet",
	"optimism-sepolia",
	"arbitrum-goerli",
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
	var isMainnet bool = false

	if network == "mainnet" {
		isMainnet = true
	}

	for _, networkDetails := range NetworkTestnetList {
		details, err := GetNetwork(networkDetails)

		if details.Mainnet != isMainnet {
			continue
		}

		if err != nil {
			return Wallets{}, err
		}

		wallets := database.GetWallets(password)

		for _, wallet := range wallets {
			list := details.RPC
			client, err := ethclient.Dial(list)
			if err != nil {
				log.Fatal(err)
			}

			account := common.HexToAddress(wallet.Address)
			balance, err := client.BalanceAt(context.Background(), account, nil)
			if err != nil {
				fmt.Println(list, wallet.Address)
				fmt.Println(err)
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

			if cryptoBalance != 0 {
				walletData := Wallet{
					Blockchain:    details.FullName,
					CryptoBalance: cryptoBalance,
					FiatBalance:   cryptoBalance * data.Data.Price,
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

func GetNetwork(network string) (networkDetails, error) {
	switch network {
	case "ethereum":
		return networkDetails{
			Name:     "ethereum",
			FullName: "ethereum mainnet",
			ChainID:  1,
			Currency: "eth",
			Mainnet:  true,
			RPC:      "https://etherscan.io",
			Explorer: "https://1rpc.io/eth",
		}, nil
	case "optimism":
		return networkDetails{
			Name:     "optimism",
			FullName: "optimism mainnet",
			ChainID:  10,
			Currency: "eth",
			Mainnet:  true,
			RPC:      "https://1rpc.io/op",
			Explorer: "https://optimistic.etherscan.io",
		}, nil
	case "polygon-mainnet":
		return networkDetails{
			Name:     "polygon",
			FullName: "Polygon mainnet",
			ChainID:  137,
			Currency: "matic",
			Mainnet:  true,
			RPC:      "https://1rpc.io/matic",
			Explorer: "https://polygonscan.com",
		}, nil
	case "arbitrum-one":
		return networkDetails{
			Name:     "arbitrum",
			FullName: "arbitrum one mainnet",
			ChainID:  42161,
			Currency: "matic",
			Mainnet:  true,
			RPC:      "https://1rpc.io/arb",
			Explorer: "https://arbiscan.io",
		}, nil
	case "polygon-zkevm":
		return networkDetails{
			Name:     "Polygon",
			FullName: "polygon zk evm mainnet",
			ChainID:  1101,
			Currency: "matic",
			Mainnet:  true,
			RPC:      "https://zkevm-rpc.com",
			Explorer: "https://zkevm.polygonscan.com/",
		}, nil
	case "astar":
		return networkDetails{
			Name:     "astar",
			FullName: "astar mainnet",
			ChainID:  592,
			Currency: "astr",
			Mainnet:  true,
			RPC:      "https://1rpc.io/astr",
			Explorer: "https://astar.subscan.io",
		}, nil
	case "ethereum-sepolia":
		return networkDetails{
			Name:     "ethereum",
			FullName: "ethereum sepolia",
			ChainID:  11155111,
			Currency: "eth",
			Mainnet:  false,
			RPC:      "https://eth-sepolia.public.blastapi.io",
			Explorer: "https://sepolia.etherscan.io",
		}, nil
	case "ethereum-goerli":
		return networkDetails{
			Name:     "ethereum",
			FullName: "ethereum goerli",
			ChainID:  5,
			Currency: "eth",
			Mainnet:  false,
			RPC:      "https://eth-goerli.public.blastapi.io",
			Explorer: "https://goerli.etherscan.io",
		}, nil
	case "optimism-sepolia":
		return networkDetails{
			Name:     "optimism",
			FullName: "optimism sepolia",
			ChainID:  11155420,
			Currency: "eth",
			Mainnet:  false,
			RPC:      "https://optimism-goerli.public.blastapi.io",
			Explorer: "https://sepolia-optimistic.etherscan.io/",
		}, nil
	case "arbitrum-goerli":
		return networkDetails{
			Name:     "arbitrum",
			FullName: "arbitrum goerli",
			ChainID:  5,
			Currency: "eth",
			Mainnet:  false,
			RPC:      "https://1rpc.io/astr",
			Explorer: "https://astar.subscan.io",
		}, nil
	case "polygon-zkevm-testnet":
		return networkDetails{
			Name:     "polygon",
			FullName: "polygon zk evm testnet",
			ChainID:  1422,
			Currency: "eth",
			Mainnet:  false,
			RPC:      "https://rpc.public.zkevm-test.net",
			Explorer: "https://zkevm.polygonscan.com/",
		}, nil
	case "polygon-mumbai":
		return networkDetails{
			Name:     "polygon",
			FullName: "polygon mumbai",
			ChainID:  80001,
			Currency: "matic",
			Mainnet:  false,
			RPC:      "https://polygon-testnet.public.blastapi.io",
			Explorer: "https://mumbai.polygonscan.com",
		}, nil
	case "binance-smart-chain":
		return networkDetails{
			Name:     "binance",
			FullName: "binance smart chain",
			ChainID:  56,
			Currency: "bnb",
			Mainnet:  false,
			RPC:      "https://bsc-dataseed1.defibit.io",
			Explorer: "https://bscscan.com",
		}, nil
	default:
		return networkDetails{
			Name:     "",
			ChainID:  0,
			Currency: "",
			Mainnet:  false,
		}, fmt.Errorf(network, "not found")
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

func SendNativePayment() {
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

func SendERC20() {
	return
}
