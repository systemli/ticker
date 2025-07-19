package cache

import (
	"sync"
	"time"

	"github.com/systemli/ticker/internal/logger"
)

var log = logger.GetWithPackage("cache")

// Cache is a simple in-memory cache with expiration.
type Cache struct {
	items sync.Map
	close chan struct{}
}

type item struct {
	data    interface{}
	expires int64
}

// NewCache creates a new cache with a cleaning interval.
func NewCache(cleaningInterval time.Duration) *Cache {
	cache := &Cache{
		close: make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(cleaningInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()

				cache.items.Range(func(key, value interface{}) bool {
					item := value.(item)

					if item.expires > 0 && now > item.expires {
						cache.items.Delete(key)
					}

					return true
				})

			case <-cache.close:
				return
			}
		}
	}()

	return cache
}

// Get returns a value from the cache.
func (cache *Cache) Get(key interface{}) (interface{}, bool) {
	obj, exists := cache.items.Load(key)

	if !exists {
		log.WithField("key", key).Debug("cache miss")
		return nil, false
	}

	item := obj.(item)

	if item.expires > 0 && time.Now().UnixNano() > item.expires {
		log.WithField("key", key).Debug("cache expired")
		return nil, false
	}

	log.WithField("key", key).Debug("cache hit")
	return item.data, true
}

// Set stores a value in the cache.
func (cache *Cache) Set(key interface{}, value interface{}, duration time.Duration) {
	var expires int64

	if duration > 0 {
		expires = time.Now().Add(duration).UnixNano()
	}

	cache.items.Store(key, item{
		data:    value,
		expires: expires,
	})
}

// Range loops over all items in the cache.
func (cache *Cache) Range(f func(key, value interface{}) bool) {
	now := time.Now().UnixNano()

	fn := func(key, value interface{}) bool {
		item := value.(item)

		if item.expires > 0 && now > item.expires {
			return true
		}

		return f(key, item.data)
	}

	cache.items.Range(fn)
}

// Delete removes a value from the cache.
func (cache *Cache) Delete(key interface{}) {
	cache.items.Delete(key)
}

// Close stops the cleaning interval and clears the cache.
func (cache *Cache) Close() {
	cache.close <- struct{}{}
	cache.items = sync.Map{}
}
