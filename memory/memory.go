/*
Package memory provides a caching module for use with k6, an open-source load testing tool. This package uses
an in-memory cache with configurable expiration and cleanup intervals, allowing for temporary data storage across
virtual users (VUs).

Import Path:
	const ImportPath = "k6/x/working-memory"

Package memory utilizes github.com/patrickmn/go-cache for cache management, allowing for item expiration and cleanup.
*/

package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"go.k6.io/k6/js/modules"
)

const ImportPath = "k6/x/working-memory"

// globalCacheInstance is a singleton instance of Cache, ensuring a single cache is used across VUs.
var (
	globalCacheInstance *Cache
	once                sync.Once
)

// Cache struct encapsulates an in-memory cache with a mutex for concurrent access management.
type Cache struct {
	cache *cache.Cache
	mutex sync.Mutex
}

// Init initializes the Cache with a default expiration time and cleanup interval.
// Parameters:
// - defaultExpiration: Cache expiration time in seconds.
// - cleanupInterval: Interval in seconds at which expired items are removed from the cache.
func (c *Cache) Init(defaultExpiration, cleanupInterval int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expiration := time.Duration(defaultExpiration) * time.Second
	cleanup := time.Duration(cleanupInterval) * time.Second
	c.cache = cache.New(expiration, cleanup)
}

// New creates and returns a new instance of rootModule, implementing the k6 modules.Module interface.
func New() modules.Module {
	once.Do(func() {
		globalCacheInstance = &Cache{}
	})
	return new(rootModule)
}

// rootModule serves as the main entry point for the k6 module, exposing cache functions.
type rootModule struct{}

// NewModuleInstance creates and returns a new module instance for each VU.
func (*rootModule) NewModuleInstance(_ modules.VU) modules.Instance {
	instance := &module{
		exports: modules.Exports{
			Default: globalCacheInstance,
			Named: map[string]interface{}{
				"init":  globalCacheInstance.Init,
				"set":   globalCacheInstance.Set,
				"get":   globalCacheInstance.Get,
				"flush": globalCacheInstance.Flush,
			},
		},
	}
	return instance
}

// module defines a structure with exported functions for use in k6 scripts.
type module struct {
	exports modules.Exports
}

// Exports returns the module's exported functions.
func (mod *module) Exports() modules.Exports {
	return mod.exports
}

// Set stores a value in the cache under the specified id with an optional expiration.
// Parameters:
// - id: Unique identifier for the cache entry.
// - value: Value to store in the cache.
// - expiration: Optional expiration time in seconds. Defaults to the cache's default expiration if omitted.
// Returns:
// - A boolean indicating if the value was successfully set.
// - An error if the cache is not initialized or another issue occurs.
func (c *Cache) Set(id, value string, expiration ...int) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.cache == nil {
		return false, errors.New("cache not initialized: please call init() first")
	}
	var exp time.Duration
	if len(expiration) > 0 {
		exp = time.Duration(expiration[0]) * time.Second
	} else {
		exp = cache.DefaultExpiration
	}
	c.cache.Set(id, value, exp)
	_, found := c.cache.Get(id)
	return found, nil
}

// Get retrieves a value from the cache by its id.
// Parameters:
// - id: Unique identifier of the cache entry.
// Returns:
// - The cached value if found, or nil if the entry does not exist.
// - An error if the cache is not initialized or another issue occurs.
func (c *Cache) Get(id string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.cache == nil {
		return nil, errors.New("cache not initialized: please call init() first")
	}
	value, found := c.cache.Get(id)
	if found {
		return value.(string), nil
	}
	return nil, nil
}

// Flush clears all items from the cache, effectively resetting it.
func (c *Cache) Flush() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.cache == nil {
		return errors.New("cache not initialized: please call init() first")
	}
	c.cache.Flush()
	return nil
}
