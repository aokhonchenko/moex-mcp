package cache

import (
	"testing"
	"time"
)

func TestCacheSetGet(t *testing.T) {
	c := New(1 * time.Hour)

	c.Set("key1", "value1")

	val, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if val.(string) != "value1" {
		t.Errorf("expected 'value1', got %v", val)
	}
}

func TestCacheMiss(t *testing.T) {
	c := New(1 * time.Hour)

	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestCacheExpiration(t *testing.T) {
	c := New(50 * time.Millisecond)

	c.Set("key1", "value1")

	// Сразу — попадание
	_, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected hit before expiration")
	}

	// Ждём истечения TTL
	time.Sleep(100 * time.Millisecond)

	_, ok = c.Get("key1")
	if ok {
		t.Fatal("expected miss after expiration")
	}
}

func TestCacheSetWithTTL(t *testing.T) {
	c := New(1 * time.Hour)

	c.SetWithTTL("short", "data", 50*time.Millisecond)

	_, ok := c.Get("short")
	if !ok {
		t.Fatal("expected hit")
	}

	time.Sleep(100 * time.Millisecond)

	_, ok = c.Get("short")
	if ok {
		t.Fatal("expected miss after custom TTL")
	}
}

func TestCacheDelete(t *testing.T) {
	c := New(1 * time.Hour)

	c.Set("key1", "value1")
	c.Delete("key1")

	_, ok := c.Get("key1")
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestCacheClear(t *testing.T) {
	c := New(1 * time.Hour)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Clear()

	stats := c.Stats()
	if stats.Size != 0 {
		t.Errorf("expected size 0 after clear, got %d", stats.Size)
	}
	if stats.Hits != 0 {
		t.Errorf("expected hits 0 after clear, got %d", stats.Hits)
	}
}

func TestCacheStats(t *testing.T) {
	c := New(1 * time.Hour)

	c.Set("a", 1)
	c.Get("a") // hit
	c.Get("b") // miss
	c.Get("a") // hit

	stats := c.Stats()
	if stats.Size != 1 {
		t.Errorf("expected size 1, got %d", stats.Size)
	}
	if stats.Hits != 2 {
		t.Errorf("expected hits 2, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("expected misses 1, got %d", stats.Misses)
	}
}

func TestCacheOverwrite(t *testing.T) {
	c := New(1 * time.Hour)

	c.Set("key", "old")
	c.Set("key", "new")

	val, ok := c.Get("key")
	if !ok {
		t.Fatal("expected hit")
	}
	if val.(string) != "new" {
		t.Errorf("expected 'new', got %v", val)
	}
}

func TestCacheConcurrent(t *testing.T) {
	c := New(1 * time.Hour)

	done := make(chan bool)

	// Писатели
	for i := 0; i < 10; i++ {
		go func(n int) {
			for j := 0; j < 100; j++ {
				c.Set("key", n*100+j)
			}
			done <- true
		}(i)
	}

	// Читатели
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				c.Get("key")
			}
			done <- true
		}()
	}

	// Ожидаем завершения
	for i := 0; i < 20; i++ {
		<-done
	}

	// Проверяем, что кэш не сломан
	stats := c.Stats()
	if stats.Size > 1 {
		t.Errorf("expected at most 1 item, got %d", stats.Size)
	}
}
