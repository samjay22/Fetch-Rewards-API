package EDA

import (
	"Fetch-Rewards-API/Interfaces"
	"Fetch-Rewards-API/Structs"
)

// eventBus represents the event bus that handles event subscription and dispatching
type eventBus struct {
	subscribers map[string][]chan<- *Structs.Event
}

// NewEventBus creates a new instance of the event bus
func NewEventBus() Interfaces.EventBus {
	return &eventBus{
		subscribers: make(map[string][]chan<- *Structs.Event),
	}
}

// Subscribe adds a new subscriber for a given event type
func (eb *eventBus) Subscribe(eventType string, subscriber chan<- *Structs.Event) {
	eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber)
}

// Publish sends an event to all subscribers of a given event type
func (eb *eventBus) Publish(event *Structs.Event) {
	subscribers := eb.subscribers[event.Type]
	for _, subscriber := range subscribers {
		subscriber <- event
	}
}
