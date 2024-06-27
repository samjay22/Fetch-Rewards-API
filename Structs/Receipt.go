package Structs

// Receipt represents the structure of a receipt
type Receipt struct {
	Id           string          `json:"id"`
	Retailer     string          `json:"retailer"`
	PurchaseDate string          `json:"purchaseDate"`
	PurchaseTime string          `json:"purchaseTime"`
	Total        string          `json:"total"`
	Items        []PurchasedItem `json:"items"`
}
