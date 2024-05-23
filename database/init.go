package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

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

func CreateNewWallet() (*os.File, error) {
	dir, _ := os.Getwd()
	filepath := filepath.Join(dir, "files", "pocket")

	err = os.MkdirAll(filepath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
	}

	file, err := os.Create(filepath + "wallet.db")
	if err != nil {
		return file, err
	}

	file.Close()
	return file, nil
}
