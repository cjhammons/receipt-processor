package main

import (
	"cjhammons/receipt-processor/db"
	"cjhammons/receipt-processor/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create router
	router := gin.Default()

	// Create receipt handler
	receiptHandler := &handlers.ReceiptHandler{
		DB: db.DB,
	}

	// Set up routes
	router.POST("/receipts/process", receiptHandler.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptHandler.GetPoints)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
