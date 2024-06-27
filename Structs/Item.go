package Structs

// Item represents an item in the receipt
type PurchasedItem struct {
	Id               string `json:"id"`
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Special item, is it an item that gives a bonus on top of inital cost, ie, kellogs spend $20 get 1000 points.
type SpecialItem struct {
	PurchasedItem // inherit

	//This maps to the handler that deals with the logic of the special item.
	ItemHandlerId string `json:"itemHandler"`
}
