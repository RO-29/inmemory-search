package inmemorycache

import "sync"

// Cache represents a cache.
type Cache interface {
	// Set sets a value to the cache.
	Set(key string, value interface{})
	// Get looks up a key's value from the cache.
	Get(key string) (value interface{}, found bool)
	// PrefixMatch looks up all key matching the prefix / suffix pattern
	PrefixMatch(prefix string, total int) []string
}

// New returns a new Cache.
// All methods are by default concurrent safe by *mutex lock*.
func New() Cache {
	t := &trie{}
	return &lock{
		Cache: t,
	}
}

func NewDIProvider() func() Cache {
	var c Cache
	var mu sync.Mutex
	return func() Cache {
		mu.Lock()
		defer mu.Unlock()
		if c == nil {
			c = New()
		}
		return c
	}
}

type lock struct {
	mu sync.Mutex
	Cache
}

func (c *lock) Set(key string, value interface{}) {
	c.mu.Lock()
	c.Cache.Set(key, value)
	c.mu.Unlock()
}

func (c *lock) Get(key string) (value interface{}, found bool) {
	c.mu.Lock()
	value, found = c.Cache.Get(key)
	c.mu.Unlock()
	return value, found
}

func (c *lock) PrefixMatch(prefix string, total int) []string {
	c.mu.Lock()
	value := c.Cache.PrefixMatch(prefix, total)
	c.mu.Unlock()
	return value
}
