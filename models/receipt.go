package models

import (
	"database/sql"
	"math"
	"strconv"
	"strings"
	"time"
)

type Receipt struct {
	ID           string        `json:"id"`
	Retailer     string        `json:"retailer"`
	PurchaseDate string        `json:"purchaseDate"`
	PurchaseTime string        `json:"purchaseTime"`
	Total        string        `json:"total"`
	Points       int           `json:"points"`
	CreatedAt    time.Time     `json:"createdAt"`
	ReceiptItems []ReceiptItem `json:"receiptItems"`
}

type ReceiptItem struct {
	ID               int    `json:"id"`
	ReceiptID        string `json:"receiptId"`
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func InitDB(db *sql.DB) error {
	// Create receipts table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS receipts (
			id TEXT PRIMARY KEY,
			retailer TEXT NOT NULL,
			purchase_date TEXT NOT NULL,
			purchase_time TEXT NOT NULL,
			total TEXT NOT NULL,
			points INTEGER NOT NULL,
			created_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create receipt_items table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS receipt_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			receipt_id TEXT NOT NULL,
			short_description TEXT NOT NULL,
			price TEXT NOT NULL,
			FOREIGN KEY (receipt_id) REFERENCES receipts(id)
		)
	`)
	return err
}

func (receipt *Receipt) SaveReceipt(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Insert receipt
	_, err = tx.Exec(`
		INSERT INTO receipts (
			id, 
			retailer, 
			purchase_date, 
			purchase_time, 
			total, 
			points, 
			created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, receipt.ID, receipt.Retailer, receipt.PurchaseDate, receipt.PurchaseTime, receipt.Total, receipt.Points, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert items
	for _, item := range receipt.ReceiptItems {
		_, err = tx.Exec(`
			INSERT INTO receipt_items (
				receipt_id, 
				short_description, 
				price
			)
			VALUES (?, ?, ?)
		`, receipt.ID, item.ShortDescription, item.Price)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetReceipt(id string, db *sql.DB) (*Receipt, []ReceiptItem, error) {
	// Get receipt
	receipt := &Receipt{}
	err := db.QueryRow(`
		SELECT 
			id, 
			retailer, 
			purchase_date, 
			purchase_time, 
			total, 
			points
		FROM receipts
		WHERE id = ?
	`, id).Scan(&receipt.ID, &receipt.Retailer, &receipt.PurchaseDate, &receipt.PurchaseTime, &receipt.Total, &receipt.Points)
	if err != nil {
		return nil, nil, err
	}

	// Get items
	rows, err := db.Query(`
		SELECT 
			id, 
			short_description, 
			price
		FROM receipt_items
		WHERE receipt_id = ?
	`, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var items []ReceiptItem
	for rows.Next() {
		var item ReceiptItem
		err = rows.Scan(&item.ID, &item.ShortDescription, &item.Price)
		if err != nil {
			return nil, nil, err
		}
		item.ReceiptID = id
		items = append(items, item)
	}

	return receipt, items, nil
}

func (receipt *Receipt) CalculatePoints() int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			points++
		}
	}

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == math.Floor(total) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(total*100, 25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += (len(receipt.ReceiptItems) / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3
	for _, item := range receipt.ReceiptItems {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 == 1 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}
