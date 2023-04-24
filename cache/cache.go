package cache

import (
	"errors"
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

// Pull pulls the elements from the given Cache into this; if the two Caches
// have some elements in common, the incoming elements replace the existing ones.
func (c *Cache[K, V]) Pull(other *Cache[K, V]) error {
	if other == nil {
		if c.logger != nil {
			c.logger.Error("merging with nil cache")
		}
		return errors.New("invalid cache")
	}

	if c.logger != nil {
		c.logger.Debug("pulling other caches elements into this")
	}

	keys := other.Keys()
	for _, k := range keys {
		v, _ := other.Get(k)
		c.Put(k, v)
	}
	if c.logger != nil {
		c.logger.Debug("dne pulling other caches elements into this")
	}
	return nil
}

// Merge pulls the elements from the given Cache into this; if the two Caches
// have some elements in common, the existing ones are preserved.
func (c *Cache[K, V]) Merge(other *Cache[K, V]) error {
	if other == nil {
		if c.logger != nil {
			c.logger.Error("merging with nil cache")
		}
		return errors.New("invalid cache")
	}

	if c.logger != nil {
		c.logger.Debug("pulling other caches elements into this")
	}

	keys := other.Keys()
	for _, k := range keys {
		v, _ := other.Get(k)
		c.Put(k, v)
	}
	if c.logger != nil {
		c.logger.Debug("dne pulling other caches elements into this")
	}
	return nil
}

func (c *Cache[K, V]) Store() error {
	if c.logger != nil {
		c.logger.Debug("persisting cache")
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.storeNoLock(true); err != nil {
		if c.logger != nil {
			c.logger.Error("error persisting cache", "error", err)
		}
		return err
	}
	if c.logger != nil {
		c.logger.Debug("done persisting cache")
	}
	return nil
}

func (c *Cache[K, V]) Load() error {
	if c.logger != nil {
		c.logger.Debug("loading cache")
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	return c.loadNoLock()
}

// Put stores an element in the cache; if ana element already exists, it
// does not replace it and keeps the previous value.
func (c *Cache[K, V]) Put(k K, v V) bool {
	if c.logger != nil {
		c.logger.Debug("putting value into cache", "key", k, "value", v)
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.store[k]; !ok {
		c.store[k] = v
		if c.logger != nil {
			c.logger.Debug("value stored into cache", "key", k, "value", v)
		}
		c.storeNoLock(false)
		return true
	}
	return false
}

// Replace stores an element in the cache, possibly replacing an existing
// one under the same key; it returns whether an elements was already
// present in the Cache under the same key and, if so, its value.
func (c *Cache[K, V]) Replace(k K, v V) (V, bool) {
	if c.logger != nil {
		c.logger.Debug("putting value into cache", "key", k, "value", v)
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	old, ok := c.store[k]
	c.store[k] = v
	c.storeNoLock(false)
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
	err := c.storeNoLock(false)
	if c.logger != nil {
		c.logger.Debug("removed value from cache", "present", ok, "key", k, "value", v, "error", err)
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
	err := c.storeNoLock(false)
	if c.logger != nil {
		c.logger.Debug("cache clear", "error", err)
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

// storeNoLock persists the cache without acquiring the read lock,
// which should be held by the caller; not acquiring the lock before
// calling this method can result in unexpected behaviour.
func (c *Cache[K, V]) storeNoLock(force bool) error {
	if c.logger != nil {
		c.logger.Debug("storing the cache without acquiring the lock")
	}

	if !force && !c.policy.Trigger() {
		if c.logger != nil {
			c.logger.Debug("neither policy not user requie the cache to be stored")

		}
		return nil
	}

	data, err := c.encoding.Encode(c.store)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("error encoding cache", "error", err)
		}
		return err
	}

	err = c.persistence.Write(data)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("error persisting cache", "error", err)
		}
		return err
	}

	if c.logger != nil {
		c.logger.Debug("cache stored with no lock acquired")
	}
	return nil
}

// loadNoLock read back the cache without acquiring the write lock,
// which should be held by the caller; not acquiring the lock before
// calling this method can result in unexpected behaviour.
func (c *Cache[K, V]) loadNoLock() error {
	if c.logger != nil {
		c.logger.Debug("loading the cache without acquiring the lock")
	}

	data, err := c.persistence.Read()
	if err != nil {
		if c.logger != nil {
			c.logger.Error("error reading cache data from persistence", "error", err)
		}
		return err
	}

	if c.logger != nil {
		c.logger.Debug("data read, decoding...")
	}

	m, err := c.encoding.Decode((data))
	if err != nil {
		if c.logger != nil {
			c.logger.Error("error decoding the cache from data", "error", err)
		}
		return err
	}

	c.store = m

	if c.logger != nil {
		c.logger.Debug("cache loaded with no lock acquired")
	}
	return nil
}
