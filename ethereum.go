package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func CreateNewWallet() (string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	return hexutil.Encode(privateKeyBytes), nil
}
