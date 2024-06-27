package Interfaces

type ReceiptFilterRule struct {
	Id           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Points       string `json:"points"`
}

// ApplyDefaults sets default values to fields that are empty
func (r *ReceiptFilterRule) ApplyDefaults() {
	if r.Id == "" {
		r.Id = "%"
	}
	if r.Retailer == "" {
		r.Retailer = "%"
	}
	if r.PurchaseDate == "" {
		r.PurchaseDate = "%"
	}
	if r.PurchaseTime == "" {
		r.PurchaseTime = "%"
	}
	if r.Total == "" {
		r.Total = "%"
	}
	if r.Points == "" {
		r.Points = "%"
	}
}
