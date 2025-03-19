package memory

import (
	"testing"
	"time"
)

func TestCacheDAO_SetAndGet(t *testing.T) {
	cache := NewCacheDAO()
	key := "testKey"
	value := "testValue"
	ttl := int64(1) // 1 second

	// Test Set
	err := cache.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get
	got, err := cache.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != value {
		t.Errorf("Get returned wrong value: got %v, want %v", got, value)
	}

	// Test expiration
	time.Sleep(2 * time.Second)
	_, err = cache.Get(key)
	if err == nil {
		t.Error("Expected error for expired key, got nil")
	}
}

func TestCacheDAO_SetNX(t *testing.T) {
	cache := NewCacheDAO()
	key := "testKeyNX"
	value := "testValueNX"
	ttl := int64(1) // 1 second

	// Test SetNX
	_, err := cache.SetNX(key, value, ttl)
	if err != nil {
		t.Fatalf("SetNX failed: %v", err)
	}

	// Test SetNX with existing key
	_, err = cache.SetNX(key, "newValue", ttl)
	if err == nil {
		t.Error("Expected error for existing key, got nil")
	}
}

func TestCacheDAO_LockAndUnlock(t *testing.T) {
	cache := NewCacheDAO()
	key := "testLockKey"

	// Test Lock
	err := cache.GetDistributeLock(key, 1*time.Second).Lock()
	if err != nil {
		t.Fatalf("Lock failed: %v", err)
	}

	// Test Unlock
	_, err = cache.GetDistributeLock(key, 1*time.Second).Unlock()
	if err != nil {
		t.Fatalf("Unlock failed: %v", err)
	}

	// Test Lock after Unlock to ensure it's unlocked
	err = cache.GetDistributeLock(key, 1*time.Second).Lock()
	if err != nil {
		t.Fatalf("Lock after Unlock failed: %v", err)
	}
}
