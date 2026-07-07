package cache

import (
	"sync"
	"time"
)

// Entry — запись в кэше.
type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// Cache — потокобезопасный in-memory кэш с TTL.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]Entry
	ttl     time.Duration
	hits    int64
	misses  int64
}

// New создаёт кэш с заданным TTL.
func New(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]Entry),
		ttl:   ttl,
	}
	// Запускаем фоновую очистку просроченных записей
	go c.cleanup()
	return c
}

// Get возвращает значение из кэша или nil, если запись отсутствует или просрочена.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		c.mu.Lock()
		c.misses++
		c.mu.Unlock()
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.misses++
		c.mu.Unlock()
		return nil, false
	}

	c.mu.Lock()
	c.hits++
	c.mu.Unlock()
	return entry.Value, true
}

// Set сохраняет значение в кэше с TTL по умолчанию.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// SetWithTTL сохраняет значение с индивидуальным TTL.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Delete удаляет запись из кэша.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear очищает весь кэш.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Entry)
	c.hits = 0
	c.misses = 0
}

// Stats возвращает статистику кэша.
type Stats struct {
	Size   int   `json:"size"`
	Hits   int64 `json:"hits"`
	Misses int64 `json:"misses"`
}

func (c *Cache) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Stats{
		Size:   len(c.items),
		Hits:   c.hits,
		Misses: c.misses,
	}
}

// cleanup периодически удаляет просроченные записи.
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.items {
			if now.After(entry.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
