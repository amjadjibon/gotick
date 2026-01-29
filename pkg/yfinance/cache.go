package yfinance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// CacheType represents the type of cache storage
type CacheType int

const (
	// CacheTypeMemory stores cache in memory only
	CacheTypeMemory CacheType = iota
	// CacheTypeDisk stores cache on disk
	CacheTypeDisk
	// CacheTypeBoth stores cache in both memory and disk
	CacheTypeBoth
)

// CacheConfig configures the cache behavior
type CacheConfig struct {
	Type       CacheType
	Directory  string        // For disk cache
	DefaultTTL time.Duration // Default TTL for cache entries
	MaxSize    int           // Maximum number of entries in memory cache
}

// DefaultCacheConfig returns the default cache configuration
func DefaultCacheConfig() CacheConfig {
	homeDir, _ := os.UserHomeDir()
	return CacheConfig{
		Type:       CacheTypeMemory,
		Directory:  filepath.Join(homeDir, ".yfinance_cache"),
		DefaultTTL: 5 * time.Minute,
		MaxSize:    1000,
	}
}

// cacheEntry represents a cached item
type cacheEntry struct {
	Data      []byte    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Cache provides caching functionality for API responses
type Cache struct {
	config  CacheConfig
	memory  map[string]*cacheEntry
	mu      sync.RWMutex
	enabled bool
}

// NewCache creates a new cache with the given configuration
func NewCache(config CacheConfig) *Cache {
	c := &Cache{
		config:  config,
		memory:  make(map[string]*cacheEntry),
		enabled: true,
	}

	// Create disk cache directory if needed
	if config.Type == CacheTypeDisk || config.Type == CacheTypeBoth {
		os.MkdirAll(config.Directory, 0o755)
	}

	return c
}

// defaultCache is the global cache instance
var (
	defaultCache     *Cache
	defaultCacheMu   sync.Mutex
	defaultCacheOnce sync.Once
)

// GetDefaultCache returns the default cache instance
func GetDefaultCache() *Cache {
	defaultCacheOnce.Do(func() {
		defaultCache = NewCache(DefaultCacheConfig())
	})
	return defaultCache
}

// SetDefaultCache sets the default cache instance
func SetDefaultCache(cache *Cache) {
	defaultCacheMu.Lock()
	defer defaultCacheMu.Unlock()
	defaultCache = cache
}

// EnableCache enables or disables the default cache
func EnableCache(enabled bool) {
	cache := GetDefaultCache()
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.enabled = enabled
}

// generateKey creates a cache key from the given parameters
func (c *Cache) generateKey(prefix string, params ...interface{}) string {
	data := fmt.Sprintf("%s:%v", prefix, params)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	if !c.enabled {
		return nil, false
	}

	// Try memory cache first
	if c.config.Type == CacheTypeMemory || c.config.Type == CacheTypeBoth {
		c.mu.RLock()
		entry, ok := c.memory[key]
		c.mu.RUnlock()

		if ok && time.Now().Before(entry.ExpiresAt) {
			return entry.Data, true
		}
	}

	// Try disk cache
	if c.config.Type == CacheTypeDisk || c.config.Type == CacheTypeBoth {
		data, ok := c.getFromDisk(key)
		if ok {
			// Populate memory cache
			if c.config.Type == CacheTypeBoth {
				c.mu.Lock()
				c.memory[key] = &cacheEntry{Data: data, ExpiresAt: time.Now().Add(c.config.DefaultTTL)}
				c.mu.Unlock()
			}
			return data, true
		}
	}

	return nil, false
}

// Set stores a value in the cache
func (c *Cache) Set(key string, data []byte, ttl time.Duration) {
	if !c.enabled {
		return
	}

	if ttl == 0 {
		ttl = c.config.DefaultTTL
	}

	entry := &cacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}

	// Store in memory
	if c.config.Type == CacheTypeMemory || c.config.Type == CacheTypeBoth {
		c.mu.Lock()
		// Evict if at max size
		if len(c.memory) >= c.config.MaxSize {
			c.evictOldest()
		}
		c.memory[key] = entry
		c.mu.Unlock()
	}

	// Store on disk
	if c.config.Type == CacheTypeDisk || c.config.Type == CacheTypeBoth {
		c.saveToDisk(key, entry)
	}
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.memory, key)
	c.mu.Unlock()

	if c.config.Type == CacheTypeDisk || c.config.Type == CacheTypeBoth {
		c.deleteFromDisk(key)
	}
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	c.memory = make(map[string]*cacheEntry)
	c.mu.Unlock()

	if c.config.Type == CacheTypeDisk || c.config.Type == CacheTypeBoth {
		os.RemoveAll(c.config.Directory)
		os.MkdirAll(c.config.Directory, 0o755)
	}
}

// evictOldest removes the oldest entry from memory cache
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.memory {
		if oldestKey == "" || entry.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.ExpiresAt
		}
	}

	if oldestKey != "" {
		delete(c.memory, oldestKey)
	}
}

// getFromDisk retrieves a value from disk cache
func (c *Cache) getFromDisk(key string) ([]byte, bool) {
	path := filepath.Join(c.config.Directory, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		os.Remove(path)
		return nil, false
	}

	return entry.Data, true
}

// saveToDisk saves a value to disk cache
func (c *Cache) saveToDisk(key string, entry *cacheEntry) {
	path := filepath.Join(c.config.Directory, key+".json")
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	os.WriteFile(path, data, 0o644)
}

// deleteFromDisk removes a value from disk cache
func (c *Cache) deleteFromDisk(key string) {
	path := filepath.Join(c.config.Directory, key+".json")
	os.Remove(path)
}

// CacheKey generates a cache key for API requests
func CacheKey(endpoint string, params map[string]string) string {
	cache := GetDefaultCache()
	return cache.generateKey(endpoint, params)
}

// TTL constants for different data types
const (
	TTLQuote      = 1 * time.Minute  // Quotes are short-lived
	TTLHistory    = 1 * time.Hour    // Historical data can be cached longer
	TTLInfo       = 24 * time.Hour   // Company info changes infrequently
	TTLHolders    = 24 * time.Hour   // Holders data is updated quarterly
	TTLAnalysis   = 6 * time.Hour    // Analysis data is updated periodically
	TTLSearch     = 1 * time.Hour    // Search results are fairly stable
	TTLNews       = 15 * time.Minute // News updates frequently
	TTLOptions    = 5 * time.Minute  // Options data is time-sensitive
	TTLFinancials = 24 * time.Hour   // Financial statements are quarterly
)
