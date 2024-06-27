package Interfaces

import (
	"Fetch-Rewards-API/Shared/Structs"
)

type EventBus interface {
	Subscribe(eventType string, subscriber chan<- *Structs.Event)
	Publish(event *Structs.Event)
}
