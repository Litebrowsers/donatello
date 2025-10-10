package cache

import (
	"sync"
	"time"
)

type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]Challenge
}

func NewInMemoryCache() *InMemoryCache {
	c := &InMemoryCache{items: make(map[string]Challenge)}
	go c.cleanup()
	return c
}

func (c *InMemoryCache) Set(key string, ch Challenge) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = ch
	return nil
}

func (c *InMemoryCache) Get(key string) (Challenge, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, ok := c.items[key]
	if !ok || time.Now().After(ch.ExpiresAt) {
		return Challenge{}, false, nil
	}
	return ch, true, nil
}

func (c *InMemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *InMemoryCache) cleanup() {
	for {
		time.Sleep(time.Minute)
		c.mu.Lock()
		for k, v := range c.items {
			if time.Now().After(v.ExpiresAt) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
