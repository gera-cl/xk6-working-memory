package memory

import (
	"testing"
	"time"
)

func TestCacheInit(t *testing.T) {
	cache := &Cache{}
	cache.Init(5, 10)

	if cache.cache == nil {
		t.Fatalf("expected cache to be initialized")
	}
}

func TestCacheSetAndGet(t *testing.T) {
	cache := &Cache{}
	cache.Init(5, 10)

	// Test setting a value in the cache
	id := "testKey"
	value := "testValue"
	success, err := cache.Set(id, value)
	if err != nil || !success {
		t.Fatalf("expected Set to succeed, got err: %v", err)
	}

	// Test getting the value back from the cache
	got, err := cache.Get(id)
	if err != nil {
		t.Fatalf("expected Get to succeed, got err: %v", err)
	}

	if got != value {
		t.Fatalf("expected Get to return %v, got %v", value, got)
	}
}

func TestCacheExpiration(t *testing.T) {
	cache := &Cache{}
	cache.Init(1, 1) // Set short expiration for testing

	id := "expiringKey"
	value := "expiringValue"
	cache.Set(id, value, 1) // 1-second expiration

	time.Sleep(2 * time.Second) // Wait for the item to expire

	got, err := cache.Get(id)
	if err != nil {
		t.Fatalf("expected Get to succeed after expiration, got err: %v", err)
	}
	if got != nil {
		t.Fatalf("expected Get to return nil after expiration, got %v", got)
	}
}

func TestCacheFlush(t *testing.T) {
	cache := &Cache{}
	cache.Init(5, 10)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.Flush()

	// After flushing, all keys should be removed
	got, err := cache.Get("key1")
	if err != nil {
		t.Fatalf("expected Get to succeed after flush, got err: %v", err)
	}
	if got != nil {
		t.Fatalf("expected Get to return nil after flush, got %v", got)
	}

	got, err = cache.Get("key2")
	if err != nil {
		t.Fatalf("expected Get to succeed after flush, got err: %v", err)
	}
	if got != nil {
		t.Fatalf("expected Get to return nil after flush, got %v", got)
	}
}
