package Services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	Interfaces2 "Fetch-Rewards-API/Backend/Interfaces"
	Structs2 "Fetch-Rewards-API/Shared/Structs"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type receiptService struct {
	logger      *zerolog.Logger
	dataService Interfaces2.DatabaseService
}

func NewReceiptService(logger *zerolog.Logger, dataService Interfaces2.DatabaseService) Interfaces2.ReceiptService {
	return &receiptService{
		logger:      logger,
		dataService: dataService,
	}
}

func (rt *receiptService) GetPointsForReceiptById(id string) (int64, error) {
	receipt, err := rt.getReceiptById(id)
	if err != nil {
		return 0, err
	}
	return receipt.Points, nil
}

func (rt *receiptService) ProcessReceipt(receiptEntity *Structs2.Receipt) error {
	// Generate IDs for items if not already present
	for i, item := range receiptEntity.Items {
		if item.Id == "" {
			receiptEntity.Items[i].Id = uuid.New().String()
		}
	}

	// Calculate points based on receipt rules
	points := calculatePoints(receiptEntity)

	// Insert receipt and items into the database
	err := rt.insertReceiptAndItems(receiptEntity, points)
	if err != nil {
		return fmt.Errorf("failed to process receipt: %w", err)
	}

	return nil
}

func (rt *receiptService) GetReceipts(ctx context.Context, filterBy *Interfaces2.ReceiptFilterRule, page int) (*Interfaces2.SearchPagePayload, error) {
	const pageSize = 15

	// Calculate offset for pagination
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// Prepare the filter function for GetEntityByFilterRule
	var filterFunc func(interface{}) (interface{}, error)
	filterFunc = rt.buildReceiptsFilterFunc(ctx, filterBy, pageSize, offset)

	// Call GetEntityByFilterRule with the filter function
	var err error
	var data interface{}
	data, err = rt.dataService.GetEntityByFilterRule(ctx, filterFunc)
	if err != nil {
		rt.logger.Info().Err(err)
	}

	// Assert and return the receipts
	receipts, ok := data.(*Interfaces2.SearchPagePayload)
	if !ok {
		return nil, fmt.Errorf("unexpected type for response data")
	}

	return receipts, nil
}

func (rt *receiptService) getReceiptById(id string) (*Structs2.Receipt, error) {
	// Query the database for the receipt with the given ID
	r, err := rt.dataService.GetEntityByFilterRule(context.Background(), func(dbI interface{}) (interface{}, error) {
		receipt := &Structs2.Receipt{}
		db, ok := dbI.(*sql.DB)
		if !ok {
			return nil, fmt.Errorf("unexpected type for database connection")
		}

		row := db.QueryRow("SELECT Id, Retailer, PurchaseDate, PurchaseTime, Total, Points FROM Receipts WHERE Id = ?", id)
		err := row.Scan(&receipt.Id, &receipt.Retailer, &receipt.PurchaseDate, &receipt.PurchaseTime, &receipt.Total, &receipt.Points)
		if err != nil {
			return nil, err
		}

		// Query items related to the receipt
		items, err := rt.getItemsForReceipt(db, id)
		if err != nil {
			return nil, err
		}
		receipt.Items = items

		return receipt, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error retrieving receipt from database: %w", err)
	}

	receipt, ok := r.(*Structs2.Receipt)
	if !ok {
		return nil, fmt.Errorf("unexpected type for receipt")
	}

	return receipt, nil
}

func (rt *receiptService) insertReceiptAndItems(receiptEntity *Structs2.Receipt, points int) error {
	err := rt.dataService.AddEntity(context.Background(), func(i interface{}) error {
		db, ok := i.(*sql.DB)
		if !ok {
			return fmt.Errorf("invalid database connection")
		}

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
				rt.logger.Error().Err(err).Msg("Transaction rolled back")
			} else {
				tx.Commit()
				rt.logger.Info().Msg("Transaction committed successfully")
			}
		}()

		// Insert Receipt into 'receipts' table
		_, err = tx.ExecContext(context.Background(), "INSERT INTO receipts (Id, Retailer, PurchaseDate, PurchaseTime, Total, Points) VALUES (?, ?, ?, ?, ?, ?)",
			receiptEntity.Id, receiptEntity.Retailer, receiptEntity.PurchaseDate, receiptEntity.PurchaseTime, receiptEntity.Total, points)
		if err != nil {
			return err
		}

		// Assuming 'tx' is your database transaction object

		// Prepare the SQL statement for bulk insert
		stmt, err := tx.PrepareContext(context.Background(), "INSERT INTO items (Id, ReceiptId, ShortDescription, Price) VALUES (?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		// Prepare the values to be inserted
		for _, item := range receiptEntity.Items {
			_, err = stmt.ExecContext(context.Background(), item.Id, receiptEntity.Id, item.ShortDescription, item.Price)
			if err != nil {
				return err
			}
		}

		return nil

	})

	if err != nil {
		rt.logger.Error().Err(err).Msg("Failed to insert receipt into database")
		return fmt.Errorf("failed to insert receipt into database: %w", err)
	}

	return nil
}

func (rt *receiptService) getItemsForReceipt(db *sql.DB, receiptId string) ([]Structs2.PurchasedItem, error) {
	rows, err := db.QueryContext(context.Background(), "SELECT Id, ShortDescription, Price FROM Items WHERE ReceiptId = ?", receiptId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Structs2.PurchasedItem
	for rows.Next() {
		var item Structs2.PurchasedItem
		err := rows.Scan(&item.Id, &item.ShortDescription, &item.Price)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func calculatePoints(receipt *Structs2.Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	points += len(regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(receipt.Retailer, ""))

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(totalFloat, 1) == 0 {
		points += 50
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
	if err == nil && purchaseDateTime.Day()%2 != 0 {
		points += 6
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

func (rt *receiptService) buildReceiptsFilterFunc(ctx context.Context, filterBy *Interfaces2.ReceiptFilterRule, pageSize, offset int) func(interface{}) (interface{}, error) {
	return func(db interface{}) (interface{}, error) {
		dbInstance, ok := db.(*sql.DB)
		if !ok {
			return nil, fmt.Errorf("unexpected database instance type")
		}

		// Initialize WHERE clause and arguments slice
		var args []interface{}
		whereClauses := []string{"1 = 1"} // Start with a placeholder WHERE clause

		// Append conditions based on filter criteria
		if filterBy.Id != "" {
			whereClauses = append(whereClauses, "Id LIKE ?")
			args = append(args, fmt.Sprintf("%s%%", filterBy.Id))
		}
		if filterBy.Retailer != "" {
			whereClauses = append(whereClauses, "Retailer LIKE ?")
			args = append(args, fmt.Sprintf("%s%%", filterBy.Retailer))
		}
		if filterBy.Points != "" {
			whereClauses = append(whereClauses, "Points LIKE ?")
			args = append(args, fmt.Sprintf("%s%%", filterBy.Points))
		}
		if filterBy.PurchaseTime != "" {
			whereClauses = append(whereClauses, "PurchaseTime LIKE ?")
			args = append(args, fmt.Sprintf("%s%%", filterBy.PurchaseTime))
		}
		if filterBy.PurchaseDate != "" {
			whereClauses = append(whereClauses, "PurchaseDate LIKE ?")
			args = append(args, fmt.Sprintf("%s%%", filterBy.PurchaseDate))
		}
		if filterBy.Total != "" {
			whereClauses = append(whereClauses, "Total LIKE ?")
			args = append(args, fmt.Sprintf("%s%%", filterBy.Total))
		}

		// Construct the WHERE clause by joining all conditions with "AND"
		whereClause := strings.Join(whereClauses, " AND ")

		// Construct the SQL query string with pagination
		query := fmt.Sprintf("SELECT Id, Retailer, PurchaseDate, PurchaseTime, Total, Points FROM Receipts WHERE %s ORDER BY Id DESC LIMIT %d OFFSET %d",
			whereClause, pageSize, offset)

		// Execute the query with arguments
		rows, err := dbInstance.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var receipts []Structs2.Receipt

		// Iterate over the query results and scan into Receipt structs
		for rows.Next() {
			var receipt Structs2.Receipt
			err := rows.Scan(&receipt.Id, &receipt.Retailer, &receipt.PurchaseDate, &receipt.PurchaseTime, &receipt.Total, &receipt.Points)
			if err != nil {
				return nil, err
			}

			// Query items related to the receipt
			items, err := rt.getItemsForReceipt(dbInstance, receipt.Id)
			if err != nil {
				return nil, err
			}
			receipt.Items = items

			receipts = append(receipts, receipt)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		// Calculate total rows for pagination
		countQuery := fmt.Sprintf("SELECT COUNT(Id) FROM Receipts WHERE %s", whereClause)
		var totalRows int
		err = dbInstance.QueryRowContext(ctx, countQuery, args...).Scan(&totalRows)
		if err != nil {
			return nil, err
		}

		// Calculate total pages
		totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))

		// Create a payload structure to return
		return &Interfaces2.SearchPagePayload{
			Receipts: receipts,
			MaxPages: totalPages,
		}, nil
	}
}
