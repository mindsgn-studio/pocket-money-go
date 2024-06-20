package starknet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
	"github.com/mindsgn-studio/pocket-money-go/database"
	"github.com/mindsgn-studio/pocket-money-go/logs"
)

var (
	network                string = "testnet"
	predeployedClassHash   string = "0x2794ce20e5f2ff0d40e632cb53845b9f4e526ebd8471983f7dbd355b721d5a"
	accountAddress         string = "0xdeadbeef"
	accountContractVersion int    = 0
)

func CreateNewWallet(password string) {
	// Load environment variables
	if err := godotenv.Load(fmt.Sprintf(".env.%s", network)); err != nil {
		logs.LogError(fmt.Sprintf("Error loading .env file: %s", err))
		return
	}

	// Define RPC URL
	url := "https://starknet-sepolia.reddio.com/rpc/v0_7"

	// Create new RPC client
	clientv02, err := rpc.NewProvider(url)
	if err != nil {
		logs.LogError(fmt.Sprintf("Error creating RPC provider: %s", err))
		return
	}

	// Generate random keys for the wallet
	ks, pub, private := account.GetRandomKeys()

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
	precomputedAddress, err := acnt.PrecomputeAddress(&felt.Zero, pub, accountAddressFelt, tx.ConstructorCalldata)
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
	db, err := database.OpenDatabase(directory, password)
	if err != nil {
		logs.LogError(err.Error())
		return
	}

	// Print the address to send ETH to
	fmt.Println("Send ETH to:", precomputedAddress)

	balance, err := getEthBalance(clientv02, precomputedAddress)
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
}

func getEthBalance(client *rpc.Provider, address *felt.Felt) (*big.Float, error) {
	// Call the StarkNet RPC to get the balance
	balanceFelt, err := client.GetBalance(context.Background(), address)
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %w", err)
	}

	// Convert felt to big.Float (assuming balance is in wei)
	balance := new(big.Float).SetInt(balanceFelt.ToBigInt())

	// Convert wei to ETH (1 ETH = 10^18 wei)
	ethBalance := new(big.Float).Quo(balance, big.NewFloat(1e18))

	return ethBalance, nil
}
