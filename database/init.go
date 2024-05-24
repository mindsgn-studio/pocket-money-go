package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/lucsky/cuid"
	"github.com/mindsgn-studio/pocket-wallet-ethereum/logs"
	_ "github.com/mutecomm/go-sqlcipher/v4"
)

type Wallet struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	WalletType string `json:"type"`
	Address    string `json:"address"`
}

type Environment struct {
	Secret      string `json:"secret"`
	OSType      string `json:"osType"`
	PackageName string `json:"packageName"`
}

// Database Queries
func createWalletTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE IF NOT EXISTS wallet (
		"uuid" TEXT NOT NULL PRIMARY KEY,			
		"name" TEXT NOT NULL,
		"type" TEXT NOT NULL,
		"address" TEXT NOT NULL UNIQUE,
		"private_key" TEXT NOT NULL UNIQUE,
		"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		"updated_at" NOT NULL DEFAULT CURRENT_TIMESTAMP
	  );`

	statement, err := db.Prepare(createStudentTableSQL)
	if err != nil {
		logs.LogError(err.Error())
		return
	}

	statement.Exec()
}

func InsertWallet(db *sql.DB, walletType string, private string, address string) bool {
	insertStudentSQL := `INSERT INTO wallet(type, uuid, name, private_key, address) VALUES (?, ?, ?, ?, ?)`

	statement, err := db.Prepare(insertStudentSQL)
	if err != nil {
		logs.LogError(err.Error())
		return false
	}

	uuid := cuid.New()

	_, err = statement.Exec(walletType, uuid, uuid, private, address)
	if err != nil {
		logs.LogError(err.Error())
		return false
	}

	return true
}

func GetTotalWallet(db *sql.DB) bool {
	row, err := db.Query("SELECT uuid, SUM(uuid) FROM wallet GROUP BY uuid")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	var totalWallet int = 0
	for row.Next() {
		totalWallet++
	}

	if totalWallet >= 1 {
		return true
	} else {
		return false
	}
}

func GetWallets(password string) []Wallet {
	var wallets []Wallet
	directory, err := GetDataDirectory()
	if err != nil {
		logs.LogError(err.Error())
		return wallets
	}

	db, err := OpenDatabase(directory, password)
	if err != nil {
		logs.LogError(err.Error())
		return wallets
	}

	rows, err := db.Query("SELECT uuid, name, type, address type FROM wallet")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var w Wallet
		err := rows.Scan(&w.UUID, &w.Name, &w.WalletType, &w.Address)
		if err != nil {
			return wallets
		}
		wallets = append(wallets, w)
	}
	return wallets
}

func OpenDatabase(directory string, password string) (*sql.DB, error) {
	env := getEnvironment()
	key := url.QueryEscape(password + env.Secret)
	dbname := fmt.Sprintf(directory+"/wallet.db?_pragma_key='%s'", key)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		logs.LogError(err.Error())
		return db, fmt.Errorf(err.Error())
	}

	return db, nil
}

// functions
func GetDataDirectory() (string, error) {
	env := getEnvironment()
	fileType := env.OSType

	switch fileType {
	case "macos":
		return filepath.Join("./.database"), nil
	case "windows":
		return filepath.Join("./.database"), nil
	case "android":
		dir, err := os.Stat("/data/data/com.wallet/files")
		if err != nil {
			return "", err
		}
		if !dir.IsDir() {
			return "", fmt.Errorf("Path /data/data/%s/files is not a directory", "com.wallet")
		}
		return filepath.Join("data", "data", "com.wallet", "files"), nil
	}

	return "./", nil
}

func directoryExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logs.LogError(err.Error())
			return false
		}
		return false
	}
	return true
}

func WalletExists(password string) bool {
	directory, err := GetDataDirectory()
	if err != nil {
		logs.LogError(err.Error())
		return false
	}

	db, err := OpenDatabase(directory, password)
	if err != nil {
		logs.LogError(err.Error())
		return false
	}

	exists := GetTotalWallet(db)
	if exists {
		return true
	} else {
		return false
	}
}

func InitialiseWallet(password string) bool {
	directory, err := GetDataDirectory()
	if err != nil {
		logs.LogError(err.Error())
		return false
	}

	exist := directoryExist(directory)
	if exist {
		return true
	}

	created := os.MkdirAll(directory, os.ModePerm)
	if created != nil {
		logs.LogError(err.Error())
		return false
	}

	db, err := OpenDatabase(directory, password)
	if err != nil {
		logs.LogError(err.Error())
		return false
	}
	defer db.Close()

	createWalletTable(db)

	return true
}

func getEnvironment() Environment {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Environment{
		PackageName: os.Getenv("PACKAGE_NAME"),
		OSType:      os.Getenv("OS_TYPE"),
		Secret:      os.Getenv("SECRET"),
	}
}
