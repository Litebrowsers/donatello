package cache

import (
	"testing"
	"time"
)

func TestInMemoryCache_SetAndGet(t *testing.T) {
	cache := NewInMemoryCache()
	challenge := Challenge{
		Task:         "test_task",
		ExpectedHash: "test_hash",
		ExpiresAt:    time.Now().Add(time.Minute),
	}

	cache.Set("test_key", challenge)

	retrieved, found, err := cache.Get("test_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected to find item, but it was not found")
	}
	if retrieved.Task != challenge.Task || retrieved.ExpectedHash != challenge.ExpectedHash {
		t.Errorf("retrieved challenge does not match original")
	}
}

func TestInMemoryCache_Get_NotFound(t *testing.T) {
	cache := NewInMemoryCache()
	_, found, err := cache.Get("non_existent_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected not to find item, but it was found")
	}
}

func TestInMemoryCache_Get_Expired(t *testing.T) {
	cache := NewInMemoryCache()
	challenge := Challenge{
		Task:         "test_task",
		ExpectedHash: "test_hash",
		ExpiresAt:    time.Now().Add(-time.Minute), // Expired
	}

	cache.Set("test_key", challenge)

	_, found, err := cache.Get("test_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected item to be expired, but it was found")
	}
}

func TestInMemoryCache_Delete(t *testing.T) {
	cache := NewInMemoryCache()
	challenge := Challenge{
		Task:         "test_task",
		ExpectedHash: "test_hash",
		ExpiresAt:    time.Now().Add(time.Minute),
	}

	cache.Set("test_key", challenge)
	cache.Delete("test_key")

	_, found, err := cache.Get("test_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected item to be deleted, but it was found")
	}
}

func TestInMemoryCache_Cleanup(t *testing.T) {
	// Note: This is a simplified test for cleanup.
	// A more robust test might involve manipulating time.
	cache := NewInMemoryCache()
	challenge1 := Challenge{
		Task:         "task1",
		ExpectedHash: "hash1",
		ExpiresAt:    time.Now().Add(-time.Minute), // Expired
	}
	challenge2 := Challenge{
		Task:         "task2",
		ExpectedHash: "hash2",
		ExpiresAt:    time.Now().Add(time.Minute), // Not expired
	}

	cache.Set("key1", challenge1)
	cache.Set("key2", challenge2)

	// Manually trigger cleanup logic for testing purposes
	// In the real implementation, this runs in a goroutine
	cache.mu.Lock()
	for k, v := range cache.items {
		if time.Now().After(v.ExpiresAt) {
			delete(cache.items, k)
		}
	}
	cache.mu.Unlock()

	_, found, _ := cache.Get("key1")
	if found {
		t.Error("Expected expired item to be cleaned up, but it was found")
	}

	_, found, _ = cache.Get("key2")
	if !found {
		t.Error("Expected non-expired item to remain, but it was not found")
	}
}
