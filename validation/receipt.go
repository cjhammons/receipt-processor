package validation

import (
	"regexp"
	"time"
)

var (
	retailerRegex  = regexp.MustCompile(`^[\w\s\-&]+$`)
	totalRegex     = regexp.MustCompile(`^\d+\.\d{2}$`)
	itemDescRegex  = regexp.MustCompile(`^[\w\s\-]+$`)
	itemPriceRegex = regexp.MustCompile(`^\d+\.\d{2}$`)
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func ValidateReceipt(retailer, purchaseDate, purchaseTime, total string, items []Item) []*ValidationError {
	var errors []*ValidationError

	// Validate retailer
	if !retailerRegex.MatchString(retailer) {
		errors = append(errors, &ValidationError{
			Field:   "retailer",
			Message: "Retailer name must contain only letters, numbers, spaces, hyphens, and ampersands",
		})
	}

	// Validate purchase date
	if _, err := time.Parse("2006-01-02", purchaseDate); err != nil {
		errors = append(errors, &ValidationError{
			Field:   "purchaseDate",
			Message: "Purchase date must be in YYYY-MM-DD format",
		})
	}

	// Validate purchase time
	if _, err := time.Parse("15:04", purchaseTime); err != nil {
		errors = append(errors, &ValidationError{
			Field:   "purchaseTime",
			Message: "Purchase time must be in 24-hour format (HH:MM)",
		})
	}

	// Validate total
	if !totalRegex.MatchString(total) {
		errors = append(errors, &ValidationError{
			Field:   "total",
			Message: "Total must be a decimal number with exactly two decimal places",
		})
	}

	// Validate items
	if len(items) == 0 {
		errors = append(errors, &ValidationError{
			Field:   "items",
			Message: "At least one item is required",
		})
	}

	for _, item := range items {
		if !itemDescRegex.MatchString(item.ShortDescription) {
			errors = append(errors, &ValidationError{
				Field:   "items",
				Message: "Item description must contain only letters, numbers, spaces, and hyphens",
			})
		}

		if !itemPriceRegex.MatchString(item.Price) {
			errors = append(errors, &ValidationError{
				Field:   "items",
				Message: "Item price must be a decimal number with exactly two decimal places",
			})
		}
	}

	return errors
}
