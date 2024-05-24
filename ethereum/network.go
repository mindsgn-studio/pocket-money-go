package ethereum

type networkDetails struct {
	Name       string `json:"name"`
	ChainID    int    `json:"chainID"`
	ChainIDHex string `json:"ChainIDHex"`
	Currency   string `json:"currency"`
	Mainnet    bool   `json:"mainnet"`
}

func GetNetwork(network string) networkDetails {
	switch network {
	case "ethereum-mainnet":
		return networkDetails{
			Name:       "ethereum",
			ChainID:    1,
			ChainIDHex: "0x1",
			Currency:   "eth",
			Mainnet:    true,
		}
	case "polygon-mainnet":
		return networkDetails{
			Name:       "polygon",
			ChainID:    137,
			ChainIDHex: "0x89",
			Currency:   "matic",
			Mainnet:    true,
		}
	case "gnosis":
		return networkDetails{
			Name:       "gnosis",
			ChainID:    100,
			ChainIDHex: "0x64",
			Currency:   "xdai",
			Mainnet:    true,
		}
	case "celo":
		return networkDetails{
			Name:       "celo",
			ChainID:    42220,
			ChainIDHex: "0x64",
			Currency:   "celo",
			Mainnet:    true,
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
