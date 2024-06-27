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
	ProcessReceipt(req *Structs.Receipt) error
	GetPointsForReceiptById(id string) (int, error)
	GetReceipts(ctx context.Context, searchTerm *ReceiptFilterRule, page int) (*SearchPagePayload, error)
}
