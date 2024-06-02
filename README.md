# Pocket

## Pocket Money Go
Pocket Money SDK is a library that allows you to seamlessly build blockchain agnostic wallets without worrying about internals.

This project is currently in early development and is not yet functional. It is a work in progress and is subject to significant changes, including the addition or removal of features and modifications to its functionality.

## Desktop 

## Supported Blockhains
- [ ] Ethereum
- [ ] Polygon
- [ ] Starknet
- [ ] Bitcoin

## Start

### Installation
```
    go install
```
### Test
```
    go test ./test
```

### Build
1. desktop
```
    go build
```

2. mobile
```
    gomobile bind -o ./../pocket-app/android/app/lib/ethereum.aar
```

## Functions

### Network
```
    network.Network(network string)
```

## types