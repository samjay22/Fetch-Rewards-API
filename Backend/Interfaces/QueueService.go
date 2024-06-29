package Interfaces

type Event struct {
	Type     string
	Data     interface{}
	Response chan interface{}
	Error    chan error
}

type EventDelegate func(data interface{}) (interface{}, error)

type QueueService interface {
	RegisterEventHandler(eventType string, handler EventDelegate)
	DispatchEvent(eventType string, data interface{}) (interface{}, error)
	QueueEvent(eventType string, data interface{}) (interface{}, error)
	ProcessQueue()
}
