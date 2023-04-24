package cache

import (
	"sync"

	"golang.org/x/exp/slog"
)

type Cache[K comparable, V any] struct {
	store       map[K]V
	lock        sync.RWMutex
	persistence Persistence
	policy      Policy
	encoding    Encoding[K, V]
	logger      *slog.Logger
}

// Option is the type for functional options.
type Option[K comparable, V any] func(*Cache[K, V])

// New creates a new Cache object, applying all the provided functional options.
func New[K comparable, V any](options ...Option[K, V]) *Cache[K, V] {
	c := &Cache[K, V]{
		store:       map[K]V{},
		persistence: &Discard{},
		policy:      &Never{},
		encoding:    &GOB[K, V]{},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// WithPersistence applies the persistence option to the Cache, which governs
// how the cache writes its contents to persistent storage.
func WithPersistence[K comparable, V any](p Persistence) Option[K, V] {
	return func(c *Cache[K, V]) {
		if p != nil {
			c.persistence = p
		}
	}
}

// WithPolicy applies the policy option to the Cache, which governs how often
// the Cache write its contents to persistent storage.
func WithPolicy[K comparable, V any](p Policy) Option[K, V] {
	return func(c *Cache[K, V]) {
		if p != nil {
			c.policy = p
		}
	}
}

// WithEncoding applies the encoding option to the Cache, which governs how
// the Cache encodes its contents before sending them to persistent storage.
func WithEncoding[K comparable, V any](e Encoding[K, V]) Option[K, V] {
	return func(c *Cache[K, V]) {
		if e != nil {
			c.encoding = e
		}
	}
}

// WithLogger applies the logger option to the Cache.
func WithLogger[K comparable, V any](l *slog.Logger) Option[K, V] {
	return func(c *Cache[K, V]) {
		if l != nil {
			c.logger = l
		}
	}
}

func (c *Cache[K, V]) Store() error {
	// if c.persistence != nil {
	// 	data, err := c.encoding.Encode(c.store)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return c.persistence.Persist(data)
	// }
	return nil
}

func (c *Cache[K, V]) Load() error {
	// TODO: load
	return nil
}

// Put stores an element in the cache, possibly replacing and existing
// one under the same key; it returns whether an elements was already
// present in the Cache under the same key and, if so, its value.
func (c *Cache[K, V]) Put(k K, v V) (V, bool) {
	if c.logger != nil {
		c.logger.Debug("putting value into cache", "key", k, "value", v)
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	old, ok := c.store[k]
	c.store[k] = v
	c.Store()
	if c.logger != nil {
		c.logger.Debug("returning previous value from cache", "present", ok, "key", k, "value", old)
	}
	return old, ok
}

// Get retrieves an element from the cache, returning whether it is
// presents and its value.
func (c *Cache[K, V]) Get(k K) (V, bool) {
	if c.logger != nil {
		c.logger.Debug("getting value from cache", "key", k)
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	v, ok := c.store[k]
	if c.logger != nil {
		c.logger.Debug("returning value from cache", "present", ok, "key", k, "value", v)
	}
	return v, ok
}

// Delete removes an element from the Cache given its key; it returns
// whether the element was present in the Cache and, if so, its value.
func (c *Cache[K, V]) Delete(k K) (V, bool) {
	if c.logger != nil {
		c.logger.Debug("removing value from cache", "key", k)
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	v, ok := c.store[k]
	delete(c.store, k)
	if c.logger != nil {
		c.logger.Debug("removed value from cache", "present", ok, "key", k, "value", v)
	}
	return v, ok
}

// Size returns the size of the cache.
func (c *Cache[K, V]) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	size := len(c.store)
	if c.logger != nil {
		c.logger.Debug("returning cache size", "size", size)
	}
	return size
}

// Clear removes all elements from the cache.
func (c *Cache[K, V]) Clear() {
	if c.logger != nil {
		c.logger.Debug("clearing value cache")
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.store = map[K]V{}
	if c.logger != nil {
		c.logger.Debug("cache clear")
	}
}

// Keys returns the current set of keys in the Cache.
func (c *Cache[K, V]) Keys() []K {
	keys := []K{}
	c.lock.RLock()
	defer c.lock.RUnlock()
	for k := range c.store {
		keys = append(keys, k)
	}
	if c.logger != nil {
		c.logger.Debug("returning cache keys", "keys", keys, "size", len(keys))
	}
	return keys
}
