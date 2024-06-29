package Interfaces

import (
	"Fetch-Rewards-API/Shared/Structs"
	"context"
)

// Makes life easier for page pagniation
type SearchPagePayload struct {
	Receipts []Structs.Receipt `json:"receipts"`
	MaxPages int               `json:"maxPages"`
}

type ReceiptService interface {
	GetPointsForReceiptById(id string) (int64, error)
	ProcessReceipt(receiptEntity *Structs.Receipt) error
	GetReceipts(ctx context.Context, filterBy *ReceiptFilterRule, page int) (*SearchPagePayload, error)
}
