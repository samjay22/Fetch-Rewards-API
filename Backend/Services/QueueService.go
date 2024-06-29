package Services

import (
	Interfaces2 "Fetch-Rewards-API/Backend/Interfaces"
	"fmt"
	"github.com/rs/zerolog"
	"sync"
)

type queueService struct {
	eventHandlers map[string]Interfaces2.EventDelegate
	eventQueue    chan *Interfaces2.Event
	logger        *zerolog.Logger
	mu            sync.RWMutex
}

// NewQueueService creates a new instance of queueService
func NewQueueService(logger *zerolog.Logger) Interfaces2.QueueService {
	return &queueService{
		eventHandlers: make(map[string]Interfaces2.EventDelegate),
		eventQueue:    make(chan *Interfaces2.Event, 500), // Buffer size can be adjusted based on expected load
		logger:        logger,
	}
}

func (qs *queueService) RegisterEventHandler(eventType string, handler Interfaces2.EventDelegate) {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	qs.eventHandlers[eventType] = handler
	qs.logger.Info().Str("eventType", eventType).Msg("Event handler registered")
}

func (qs *queueService) DispatchEvent(eventType string, data interface{}) (interface{}, error) {
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	if handler, ok := qs.eventHandlers[eventType]; ok {
		qs.logger.Info().Str("eventType", eventType).Msg("Dispatching event")
		return handler(data)
	}
	qs.logger.Error().Str("eventType", eventType).Msg("No handler registered for this event")
	return nil, fmt.Errorf("no handler registered for this event")
}

func (qs *queueService) QueueEvent(eventType string, data interface{}) (interface{}, error) {
	event := &Interfaces2.Event{
		Type:     eventType,
		Data:     data,
		Response: make(chan interface{}),
		Error:    make(chan error),
	}
	qs.eventQueue <- event
	qs.logger.Info().Str("eventType", eventType).Msg("Event queued")

	select {
	case res := <-event.Response:
		return res, nil
	case err := <-event.Error:
		return nil, err
	}
}

func (qs *queueService) ProcessQueue() {
	for event := range qs.eventQueue {
		go func(e *Interfaces2.Event) {
			qs.logger.Info().Str("eventType", e.Type).Msg("Processing queued event")
			response, err := qs.DispatchEvent(e.Type, e.Data)
			if err != nil {
				e.Error <- err
				qs.logger.Error().Err(err).Str("eventType", e.Type).Msg("Error processing event")
			} else {
				e.Response <- response
				qs.logger.Info().Str("eventType", e.Type).Msg("Event processed successfully")
			}
		}(event)
	}
}
