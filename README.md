# Receipt Processor Service

A Go-based web service that processes receipts and calculates reward points based on specific rules.

## Features

- Process receipts and store them in SQLite database
- Calculate reward points based on multiple rules
- RESTful API endpoints
- Input validation
- Persistent storage

## Prerequisites

- Go 1.21 or later
- SQLite3

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd receipt-processor
```

2. Install dependencies:
```bash
go mod tidy
```

## Running the Service

Start the server:
```bash
go run main.go
```

The server will start on port 8080.

## API Endpoints

### 1. Process Receipt
Process a receipt and get a unique ID.

**Endpoint:** `POST /receipts/process`

**Request Body:**
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    }
  ],
  "total": "6.49"
}
```

**Response:**
```json
{
  "id": "adb6b560-0eef-42bc-9d16-df48f30e89b2"
}
```

### 2. Get Points
Get the points awarded for a receipt.

**Endpoint:** `GET /receipts/{id}/points`

**Response:**
```json
{
  "points": 28
}
```

## Points Calculation Rules

Points are awarded based on the following rules:

1. One point for every alphanumeric character in the retailer name
2. 50 points if the total is a round dollar amount with no cents
3. 25 points if the total is a multiple of 0.25
4. 5 points for every two items on the receipt
5. If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer
6. 6 points if the day in the purchase date is odd
7. 10 points if the time of purchase is after 2:00pm and before 4:00pm

## Validation Rules

The service validates input according to these rules:

- Retailer name: Only letters, numbers, spaces, hyphens, and ampersands allowed
- Purchase date: Must be in YYYY-MM-DD format
- Purchase time: Must be in 24-hour format (HH:MM)
- Total: Must be a decimal number with exactly two decimal places
- Items: At least one item required
- Item description: Only letters, numbers, spaces, and hyphens allowed
- Item price: Must be a decimal number with exactly two decimal places

## Example Usage

1. Process a receipt:
```bash
curl -X POST http://localhost:8080/receipts/process \
  -H "Content-Type: application/json" \
  -d '{
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
      {
        "shortDescription": "Mountain Dew 12PK",
        "price": "6.49"
      }
    ],
    "total": "6.49"
  }'
```

2. Get points for the receipt:
```bash
curl -X GET http://localhost:8080/receipts/adb6b560-0eef-42bc-9d16-df48f30e89b2/points
```

## Error Handling

The service returns appropriate HTTP status codes and error messages:

- 400 Bad Request: Invalid input data
- 404 Not Found: Receipt ID not found
- 500 Internal Server Error: Server-side errors

## Database

The service uses SQLite for data persistence. The database file (`receipts.db`) is created automatically in the project directory.

## Testing

The service includes validation for all input fields. You can test different scenarios using the example curl commands provided above.
