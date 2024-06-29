package Interfaces

// CacheService represents a generic caching service interface
type CacheService interface {
	Get(key string) (interface{}, error)     // Get retrieves a cached value by key
	Set(key string, value interface{}) error // Set stores a value in the cache with the specified key
	Delete(key string) error                 // Delete removes a cached value by key
}
