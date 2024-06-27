package Controllers

// Item represents an individual item in the purchase.
type item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Request represents the JSON request structure.
type request struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}
