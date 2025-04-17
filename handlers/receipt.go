package handlers

import (
	"cjhammons/receipt-processor/models"
	"cjhammons/receipt-processor/validation"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReceiptHandler struct {
	DB *sql.DB
}

type ProcessReceiptRequest struct {
	Retailer     string            `json:"retailer"`
	PurchaseDate string            `json:"purchaseDate"`
	PurchaseTime string            `json:"purchaseTime"`
	Items        []validation.Item `json:"items"`
	Total        string            `json:"total"`
}

type ProcessReceiptResponse struct {
	ID string `json:"id"`
}

type GetPointsResponse struct {
	Points int `json:"points"`
}

func (handler *ReceiptHandler) ProcessReceipt(c *gin.Context) {
	var request ProcessReceiptRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please verify input. Invalid JSON format"})
		return
	}

	// Validate receipt
	if errors := validation.ValidateReceipt(
		request.Retailer,
		request.PurchaseDate,
		request.PurchaseTime,
		request.Total,
		request.Items,
	); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Generate a unique ID for the receipt
	id := uuid.New().String()

	// Create receipt model
	receipt := &models.Receipt{
		ID:           id,
		Retailer:     request.Retailer,
		PurchaseDate: request.PurchaseDate,
		PurchaseTime: request.PurchaseTime,
		Total:        request.Total,
		CreatedAt:    time.Now(),
	}

	// Create receipt items
	for _, item := range request.Items {
		receipt.ReceiptItems = append(receipt.ReceiptItems, models.ReceiptItem{
			ReceiptID:        id,
			ShortDescription: item.ShortDescription,
			Price:            item.Price,
		})
	}

	// Calculate points
	receipt.Points = receipt.CalculatePoints()

	// Save to database
	if err := receipt.SaveReceipt(handler.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ProcessReceiptResponse{ID: id})
}

func (handler *ReceiptHandler) GetPoints(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	receipt, _, err := models.GetReceipt(id, handler.DB)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "No receipt found for that ID"})
		return
	}

	context.JSON(http.StatusOK, GetPointsResponse{Points: receipt.Points})
}
