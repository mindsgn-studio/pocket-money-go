package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SaveToFile(data []byte, path string) error {
	return ioutil.WriteFile(path, data, 0644)
}

func CheckDatabase() bool {
	sqliteDatabase, err := sql.Open("sqlite3", "./wallet.db")
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(sqliteDatabase)
	defer sqliteDatabase.Close()

	return false
}

func DataDir() (string, error) {
	dir, err := os.Stat("/data/data/com.wallet/files")
	if err != nil {
		return "", err // Handle error if directory doesn't exist
	}
	if !dir.IsDir() {
		return "", fmt.Errorf("Path /data/data/%s/files is not a directory", "com.wallet")
	}
	return filepath.Join("data", "data", "com.wallet", "files"), nil
}

func ConfigExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file does not exist
			return false
		}
		return false
	}
	// Config file exists
	return true
}

func ReadFileContent(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return data, nil
}

func CreateNewWallet() (string, error) {
	directory, err := DataDir()
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	filePath := filepath.Join(directory, "config.txt")

	if ConfigExists(filePath) {
		data, err := ReadFileContent(filePath)
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}

		return string(data), nil
	} else {
		created := os.MkdirAll(directory, os.ModePerm)
		if created != nil {
			return "", fmt.Errorf(created.Error())
		}

		data := []byte("This is some configuration data")

		err = SaveToFile(data, filePath)
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}

		return "saved to file", nil
	}

}
