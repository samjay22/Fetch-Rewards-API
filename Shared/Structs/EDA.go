package Structs

import "time"

// Event represents a generic event structure
type Event struct {
	Type      string
	Timestamp time.Time
	Data      interface{}
}
