package ethereum

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mindsgn-studio/pocket-money-go/database"
	"github.com/mindsgn-studio/pocket-money-go/logs"
)

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
