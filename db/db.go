package db

import (
	"cjhammons/receipt-processor/models"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./receipts.db")
	if err != nil {
		return err
	}

	return models.InitDB(DB)
}
