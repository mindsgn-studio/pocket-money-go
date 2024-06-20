package starknet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
	"github.com/mindsgn-studio/pocket-money-go/database"
	"github.com/mindsgn-studio/pocket-money-go/logs"
)

var (
	network                string = "testnet"
	predeployedClassHash   string = "0x2794ce20e5f2ff0d40e632cb53845b9f4e526ebd8471983f7dbd355b721d5a"
	accountAddress         string = "0xdeadbeef"
	ethMainnetContract     string = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	accountContractVersion int    = 0
	getBalanceMethod       string = "balanceOf"
	getDecimalsMethod      string = "decimals"
	url                    string = "https://starknet-sepolia.reddio.com/rpc/v0_7"
)

func CreateNewWallet(password string) {
	// Load environment variables
	if err := godotenv.Load(fmt.Sprintf(".env.%s", network)); err != nil {
		logs.LogError(fmt.Sprintf("Error loading .env file: %s", err))
		return
	}

	// Create new RPC client
	clientv02, err := rpc.NewProvider(url)
	if err != nil {
		logs.LogError(fmt.Sprintf("Error creating RPC provider: %s", err))
		return
	}

	// Generate random keys for the wallet
	ks, pub, _ := account.GetRandomKeys()

	// Convert account address to felt
	accountAddressFelt, err := new(felt.Felt).SetString(accountAddress)
	if err != nil {
		logs.LogError("Error converting account address to felt")
		return
	}

	// Set up the account
	acnt, err := account.NewAccount(clientv02, accountAddressFelt, pub.String(), ks, accountContractVersion)
	if err != nil {
		logs.LogError(fmt.Sprintf("Error creating new account: %s", err))
		return
	}

	// Create transaction details
	tx := rpc.DeployAccountTxn{
		Nonce:               &felt.Zero, // Contract accounts start with nonce zero.
		MaxFee:              new(felt.Felt).SetUint64(4724395326064),
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV3,
		Signature:           []*felt.Felt{},
		ClassHash:           accountAddressFelt,
		ContractAddressSalt: pub,
		ConstructorCalldata: []*felt.Felt{pub},
	}

	// Precompute address
	_, err = acnt.PrecomputeAddress(&felt.Zero, pub, accountAddressFelt, tx.ConstructorCalldata)
	if err != nil {
		logs.LogError(fmt.Sprintf("Error precomputing address: %s", err))
		return
	}

	// Get data directory for the database
	directory, err := database.GetDataDirectory()
	if err != nil {
		logs.LogError(err.Error())
		return
	}

	// Open database
	_, err = database.OpenDatabase(directory, password)
	if err != nil {
		logs.LogError(err.Error())
		return
	}

	/*
		// Print the address to send ETH to
		fmt.Println("Send ETH to:", precomputedAddress)

		_, err = GetEthBalance(precomputedAddress)
		if err != nil {
			logs.LogError(fmt.Sprintf("Error getting ETH balance: %s", err))
			return
		}

		// Sign the transaction
		err = acnt.SignDeployAccountTransaction(context.Background(), &tx, precomputedAddress)
		if err != nil {
			logs.LogError(fmt.Sprintf("Error signing transaction: %s", err))
			return
		}

		// Send the transaction to the network
		resp, err := acnt.AddDeployAccountTransaction(context.Background(), rpc.BroadcastDeployAccountTxn{DeployAccountTxn: tx})
		if err != nil {
			logs.LogError(fmt.Sprintf("Error sending transaction: %s", err))
			return
		}

		// Print response
		fmt.Println("Transaction response:", resp)

		// Insert wallet details into the database
		database.InsertStarknet(db, precomputedAddress, acnt.CairoVersion, acnt.ChainId, acnt.AccountAddress, pub, private)
	*/
}

func GetEthBalance(address string) {
	fmt.Println("Starting getTokenBalance example")

	client, err := rpc.NewProvider(url)
	if err != nil {
		logs.LogError(fmt.Sprintf("Error creating RPC provider: %s", err))
	}
	fmt.Println("Established connection with the client")

	tokenAddressInFelt, err := utils.HexToFelt(ethMainnetContract)
	if err != nil {
		fmt.Println("Failed to transform the token contract address, did you give the hex address?")
		panic(err)
	}

	accountAddressInFelt, err := utils.HexToFelt(address)
	if err != nil {
		fmt.Println("Failed to transform the account address, did you give the hex address?")
		panic(err)
	}

	// Make read contract call
	tx := rpc.FunctionCall{
		ContractAddress:    tokenAddressInFelt,
		EntryPointSelector: utils.GetSelectorFromNameFelt(getBalanceMethod),
		Calldata:           []*felt.Felt{accountAddressInFelt},
	}

	fmt.Println("Making balanceOf() request")
	callResp, rpcErr := client.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	if rpcErr != nil {
		panic(rpcErr)
	}

	// Get token's decimals
	getDecimalsTx := rpc.FunctionCall{
		ContractAddress:    tokenAddressInFelt,
		EntryPointSelector: utils.GetSelectorFromNameFelt(getDecimalsMethod),
	}
	getDecimalsResp, rpcErr := client.Call(context.Background(), getDecimalsTx, rpc.BlockID{Tag: "latest"})
	if rpcErr != nil {
		panic(rpcErr)
	}

	floatValue := new(big.Float).SetInt(utils.FeltToBigInt(callResp[0]))
	floatValue.Quo(floatValue, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), utils.FeltToBigInt(getDecimalsResp[0]), nil)))

	fmt.Printf("Token balance of %s is %f", accountAddress, floatValue)
}
