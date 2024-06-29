package Interfaces

import (
	Structs2 "Fetch-Rewards-API/Shared/Structs"
	"database/sql"
)

type ItemService interface {
	GenerateItemIds(items []Structs2.PurchasedItem)
	GetItemsForReceipt(db *sql.DB, receiptId string) ([]Structs2.PurchasedItem, error)
	InsertItems(tx *sql.Tx, receiptId string, items []Structs2.PurchasedItem) error
}
