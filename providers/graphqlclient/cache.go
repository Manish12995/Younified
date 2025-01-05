package graphqlclient

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

// Cache manages caching for GraphQL queries and mutations
type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	data      []byte
	expiresAt time.Time
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]cacheItem),
	}
}

// Set adds an item to the cache with a default expiration
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	c.items[key] = cacheItem{
		data:      data,
		expiresAt: time.Now().Add(10 * time.Minute), // Default 10-minute expiration
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.data, true
}

// Clear removes expired cache entries
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, v := range c.items {
		if now.After(v.expiresAt) {
			delete(c.items, k)
		}
	}
}

// generateCacheKey creates a unique key for caching
func generateCacheKey(query string, variables map[string]interface{}) string {
	// Combine query and variables
	key := query
	if variables != nil {
		varJson, _ := json.Marshal(variables)
		key += string(varJson)
	}

	// Generate MD5 hash
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}
