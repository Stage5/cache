package cache

import (
	"errors"
	"time"
)

var (
	ErrKeyExists     = errors.New("key requested to add is already present in cache")
	ErrKeyNotInCache = errors.New("key requested is not present in cache")
)

// NewLocalCache returns an instance of a local cache.
// CacheTTL (Cache Time To Live): defines how long an item should live in the cache
// CleanupInterval: defines how often to clean up expired values from the cache
func NewLocalCache(cacheTTL, cleanupInterval time.Duration) *LocalCache {
	lc := &LocalCache{
		cache:    make(map[string]cachedValue),
		cacheTTL: cacheTTL,
		stop:     make(chan struct{}),
	}

	lc.wg.Add(1)
	go func(cleanupInterval time.Duration) {
		defer lc.wg.Done()
		lc.cleanupLoop(cleanupInterval)
	}(cleanupInterval)

	return lc
}

// InsertOne inserts a single key value into the cache. This function validates that the key does not already exist in the cache.
func (lc *LocalCache) InsertOne(key string, value interface{}) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// ensure value is not in cache
	_, ok := lc.cache[key]
	if ok {
		return ErrKeyExists
	}

	lc.cache[key] = cachedValue{
		value:             value,
		expireAtTimestamp: time.Now().Add(lc.cacheTTL).Unix(),
	}

	return nil
}

// UpdateOne updates a single value in the cache by the key. This function validates that the key exists in the cache.
func (lc *LocalCache) UpdateOne(key string, value interface{}) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// ensure value is in cache
	_, ok := lc.cache[key]
	if !ok {
		return ErrKeyNotInCache
	}

	lc.cache[key] = cachedValue{
		value:             value,
		expireAtTimestamp: time.Now().Add(lc.cacheTTL).Unix(),
	}
	return nil
}

// GetOne retrieves a value from the cache given the key. This function returns an error if no item is found in the cache for the given key.
func (lc *LocalCache) GetOne(key string) (*cachedValue, error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// ensure value is in cache
	cacheValue, ok := lc.cache[key]
	if !ok {
		return nil, ErrKeyNotInCache
	}

	return &cacheValue, nil
}

// DeleteOne deletes an item from the cache given a key if the key exists, else this function is a no-op.
func (lc *LocalCache) DeleteOne(key string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// The delete built-in function deletes the element with the specified key
	// (m[key]) from the map. If m is nil or there is no such element, delete
	// is a no-op.
	delete(lc.cache, key)
}

// GetSize retrieves the size of the cache.
func (lc *LocalCache) GetSize() int {
	return len(lc.cache)
}
