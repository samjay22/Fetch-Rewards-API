package Structs

// Item represents an item in the receipt
type PurchasedItem struct {
	Id               string `json:"id"`
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}
