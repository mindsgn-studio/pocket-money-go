package ethereum

import (
	"github.com/mindsgn-studio/pocket-wallet-ethereum/database"
)

type Wallet struct {
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
	Network    string `json:"network"`
	NetworkID  string `json:"networkID"`
	Provider   string `json:"provider"`
}

func CreatePrivateKey() string {
	// privateKey, err := crypto.GenerateKey()
	// if err != nil {
	//	return "", fmt.Errorf(err.Error())
	//}

	//privateKeyBytes := crypto.FromECDSA(privateKey)
	// address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	// hash := sha3.NewLegacyKeccak256()
	// hash.Write(publicKeyBytes[1:])

	// return hexutil.Encode(privateKeyBytes), nil
	return ""
}

func InitializeWallet() string {
	// exist := database.CheckDatabase()
	// if exist {
	//	return "false"
	// }

	file, err := database.CreateNewWallet()
	if err != nil {
		return err.Error()
	}

	return file
}
