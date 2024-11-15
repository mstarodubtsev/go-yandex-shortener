package storage

import "sync"

// Map to store URLs with thread safety
type Map struct {
	mu sync.RWMutex
	m  map[string]string
}

// AddURL adds a URL to the map
func (m *Map) AddURL(hash, url string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[hash] = url
}

// GetURL retrieves a URL from the map by its hash
func (m *Map) GetURL(hash string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	url, ok := m.m[hash]
	return url, ok
}

// GetAll retrieves a copy of all URLs in the map
func (m *Map) GetAll() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Return a copy to avoid exposing internal state
	mCopy := make(map[string]string, len(m.m))
	for k, v := range m.m {
		mCopy[k] = v
	}
	return mCopy
}

// NewMap creates a new thread-safe map
func NewMap() *Map {
	return &Map{
		m: make(map[string]string),
	}
}
