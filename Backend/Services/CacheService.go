package Services

import (
	Interfaces2 "Fetch-Rewards-API/Backend/Interfaces"
	"errors"
	"sync"
)

type MemoryCacheService struct {
	cache map[string]interface{}
	mu    sync.RWMutex
}

func NewMemoryCacheService() Interfaces2.CacheService {
	return &MemoryCacheService{
		cache: make(map[string]interface{}),
	}
}

func (m *MemoryCacheService) Get(key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if val, ok := m.cache[key]; ok {
		return val, nil
	}
	return nil, errors.New("key not found in cache")
}

func (m *MemoryCacheService) Set(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache[key] = value
	return nil
}

func (m *MemoryCacheService) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.cache, key)
	return nil
}
