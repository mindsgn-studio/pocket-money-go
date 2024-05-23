package network

func Network(name string) (string, error) {
	switch name {
	case "ethereum":
		return "ethereum", nil
	case "celo":
		return "celo", nil
	case "polygon":
		return "polygon", nil
	case "sol":
		return "sol", nil
	case "safe":
		return "safe", nil
	default:
		return "", nil
	}
}
