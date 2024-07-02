package Services

import (
	Interfaces2 "Fetch-Rewards-API/Backend/Interfaces"
	Structs2 "Fetch-Rewards-API/Shared/Structs"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"sync"
)

type itemService struct {
	dataService  Interfaces2.DatabaseService
	cacheService Interfaces2.CacheService
	logger       *zerolog.Logger
}

// NewItemService creates a new itemService instance
func NewItemService(dataService Interfaces2.DatabaseService, cacheService Interfaces2.CacheService, logger *zerolog.Logger) Interfaces2.ItemService {
	return &itemService{
		dataService:  dataService,
		cacheService: cacheService,
		logger:       logger,
	}
}

// GenerateItemIds generates unique IDs for items that do not have one
func (is *itemService) GenerateItemIds(items []Structs2.PurchasedItem) {
	for i, item := range items {
		if item.Id == "" {
			items[i].Id = uuid.New().String()
		}
	}
}

func (is *itemService) GetItemsForReceipt(db *sql.DB, receiptId string) ([]Structs2.PurchasedItem, error) {
	// Attempt to retrieve items from cache
	cachedItems, err := is.cacheService.Get("Items_" + receiptId)
	if err == nil && cachedItems != nil {
		if items, ok := cachedItems.([]Structs2.PurchasedItem); ok {
			return items, nil
		}
	}

	// Retrieve items from database
	rows, err := db.Query("SELECT Id, ShortDescription, Price FROM Items WHERE ReceiptId = ?", receiptId)
	if err != nil {
		is.logger.Error().Err(err).Msg("Failed to query items from database")
		return nil, err
	}
	defer rows.Close()

	// Channel to collect items and manage concurrency
	itemsCh := make(chan Structs2.PurchasedItem)
	var wg sync.WaitGroup

	// Concurrently fetch items
	for rows.Next() {
		var item Structs2.PurchasedItem
		err := rows.Scan(&item.Id, &item.ShortDescription, &item.Price)
		if err != nil {
			is.logger.Error().Err(err).Msg("Failed to scan item row")
			continue
		}
		wg.Add(1)
		go func(item Structs2.PurchasedItem) {
			defer wg.Done()
			itemsCh <- item
		}(item)
	}

	// Close the channel when all items are processed
	go func() {
		wg.Wait()
		close(itemsCh)
	}()

	// Collect items from the channel
	var items []Structs2.PurchasedItem
	for item := range itemsCh {
		items = append(items, item)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		is.logger.Error().Err(err).Msg("Rows error after scanning")
		return nil, err
	}

	// Cache the retrieved items
	err = is.cacheService.Set("Items_"+receiptId, items)
	if err != nil {
		is.logger.Error().Err(err).Msg("Failed to set items to cache")
	}

	return items, nil
}

// InsertItems inserts items into the database and updates the cache
func (is *itemService) InsertItems(tx *sql.Tx, receiptId string, items []Structs2.PurchasedItem) error {
	stmt, err := tx.PrepareContext(context.Background(), "INSERT INTO items (Id, ReceiptId, ShortDescription, Price) VALUES (?, ?, ?, ?)")
	if err != nil {
		is.logger.Error().Err(err).Msg("Failed to prepare insert statement")
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err = stmt.ExecContext(context.Background(), item.Id, receiptId, item.ShortDescription, item.Price)
		if err != nil {
			is.logger.Error().Err(err).Msg("Failed to execute insert statement")
			return err
		}
	}

	// Update the cache
	err = is.cacheService.Set("Items_"+receiptId, items)
	if err != nil {
		is.logger.Error().Err(err).Msg("Failed to set items to cache")
		return err
	}

	return nil
}
