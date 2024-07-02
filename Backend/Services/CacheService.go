package Services

import (
	Interfaces2 "Fetch-Rewards-API/Backend/Interfaces"
	"errors"
	"sync"
)

// MemoryCacheService provides an in-memory cache implementation
type MemoryCacheService struct {
	cache map[string]interface{} // cache stores key-value pairs
	mu    sync.RWMutex           // mu ensures thread-safe access to the cache
}

// NewMemoryCacheService creates a new instance of MemoryCacheService
func NewMemoryCacheService() Interfaces2.CacheService {
	return &MemoryCacheService{
		cache: make(map[string]interface{}), // initialize the cache map
	}
}

// Get retrieves a value from the cache by key
func (m *MemoryCacheService) Get(key string) (interface{}, error) {
	m.mu.RLock()         // acquire read lock
	defer m.mu.RUnlock() // ensure the lock is released

	// check if the key exists in the cache
	if val, ok := m.cache[key]; ok {
		return val, nil
	}
	return nil, errors.New("key not found in cache") // return error if key not found
}

// Set adds or updates a key-value pair in the cache
func (m *MemoryCacheService) Set(key string, value interface{}) error {
	m.mu.Lock()         // acquire write lock
	defer m.mu.Unlock() // ensure the lock is released

	m.cache[key] = value // set the value in the cache
	return nil
}

// Delete removes a key-value pair from the cache
func (m *MemoryCacheService) Delete(key string) error {
	m.mu.Lock()         // acquire write lock
	defer m.mu.Unlock() // ensure the lock is released

	delete(m.cache, key) // delete the key from the cache
	return nil
}

// Purge clears all key-value pairs from the cache
func (m *MemoryCacheService) Purge() error {
	m.mu.Lock()         // acquire write lock
	defer m.mu.Unlock() // ensure the lock is released

	m.cache = make(map[string]interface{}) // reinitialize the cache map
	return nil
}
