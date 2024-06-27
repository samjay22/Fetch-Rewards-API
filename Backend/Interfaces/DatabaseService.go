package Interfaces

import "context"

// DatabaseService represents a generic database service interface
type DatabaseService interface {

	// GetEntities Pass filter rule for extended functionality and use generic
	GetEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (interface{}, error)) (interface{}, error)

	// UpdateEntities Update logic as delegate based on generic type
	UpdateEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error)

	// DeleteEntities removes entities from database
	DeleteEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error)

	//Add entity
	AddEntity(ctx context.Context, delegate func(interface{}) error) error
}
