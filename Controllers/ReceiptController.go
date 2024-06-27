package Controllers

import (
	"Fetch-Rewards-API/Interfaces"
	"Fetch-Rewards-API/Structs"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type receiptController struct {
	DataService   Interfaces.DatabaseService
	NetworkClient *echo.Echo
	Logger        *zerolog.Logger
	eventBus      Interfaces.EventBus
}

type Receipt struct {
	Id           string          `json:"id"`
	Retailer     string          `json:"retailer"`
	PurchaseDate string          `json:"purchaseDate"`
	PurchaseTime string          `json:"purchaseTime"`
	Items        []PurchasedItem `json:"items"`
	Total        string          `json:"total"`
}

type PurchasedItem struct {
	Id               string `json:"id"`
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// RegisterReceiptController registers the receipt controller with the provided Echo instance.
func RegisterReceiptController(logger *zerolog.Logger, echoClient *echo.Echo, dataService Interfaces.DatabaseService) {
	controller := &receiptController{
		DataService:   dataService,
		NetworkClient: echoClient,
		Logger:        logger,
	}
	controller.NetworkClient.POST("/receipts/process", controller.ProcessReceipt)
}

// ProcessReceipt handles the receipt processing request.
func (r *receiptController) ProcessReceipt(c echo.Context) error {
	var req Receipt
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Error decoding request body")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	// Generate UUID for Receipt ID
	receiptID := uuid.New().String()

	// Create ReceiptEntity
	receiptEntity := Structs.Receipt{
		Id:           receiptID,
		Retailer:     req.Retailer,
		PurchaseDate: req.PurchaseDate,
		PurchaseTime: req.PurchaseTime,
		Total:        req.Total,
		Items:        make([]Structs.PurchasedItem, len(req.Items)),
	}

	// Convert request items to PurchasedItem slice
	for i, item := range req.Items {
		receiptEntity.Items[i] = Structs.PurchasedItem{
			Id:               uuid.New().String(),
			ShortDescription: item.ShortDescription,
			Price:            item.Price,
		}
	}

	// Calculate points based on receipt rules
	points := calculatePoints(&req)

	// Insert into Database using DataService
	err = r.DataService.AddEntity(context.Background(), func(i interface{}) error {
		db, ok := i.(*sql.DB)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Invalid database connection")
		}

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
				r.Logger.Error().Err(err).Msg("Transaction rolled back")
			} else {
				tx.Commit()
				r.Logger.Info().Msg("Transaction committed successfully")
			}
		}()

		// Insert Receipt into 'receipts' table
		queryReceipt := "INSERT INTO receipts (Id, Retailer, PurchaseDate, PurchaseTime, Total, Points) VALUES (?, ?, ?, ?, ?, ?)"
		_, err = tx.ExecContext(context.Background(), queryReceipt,
			receiptEntity.Id,
			receiptEntity.Retailer,
			receiptEntity.PurchaseDate,
			receiptEntity.PurchaseTime,
			receiptEntity.Total,
			points,
		)
		if err != nil {
			return err
		}

		// Insert each PurchasedItem into 'items' table
		queryItem := "INSERT INTO items (Id, ReceiptId, ShortDescription, Price) VALUES (?, ?, ?, ?)"
		for _, item := range receiptEntity.Items {
			_, err = tx.ExecContext(context.Background(), queryItem,
				item.Id,
				receiptEntity.Id,
				item.ShortDescription,
				item.Price,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		r.Logger.Error().Err(err).Msg("Failed to insert receipt into database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process receipt")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Received and processed receipt request. Points awarded: %d", points))
}

func calculatePoints(receipt *Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	points += len(regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(receipt.Retailer, ""))

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil {
		if math.Mod(totalFloat, 1) == 0 {
			points += 50
		}
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if totalFloat > 0 && math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += len(receipt.Items) / 2 * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength > 0 && trimmedLength%3 == 0 {
			priceFloat, _ := strconv.ParseFloat(item.Price, 64)
			additionalPoints := int(math.Ceil(priceFloat * 0.2))
			points += additionalPoints
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	purchaseDateTime, err := time.Parse("2006-01-02 15:04", receipt.PurchaseDate+" "+receipt.PurchaseTime)
	if err == nil {
		day := purchaseDateTime.Day()
		if day%2 != 0 {
			points += 6
		}
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	startTime, _ := time.Parse("15:04", "14:00")
	endTime, _ := time.Parse("15:04", "16:00")
	if purchaseTime.After(startTime) && purchaseTime.Before(endTime) {
		points += 10
	}

	return points
}
