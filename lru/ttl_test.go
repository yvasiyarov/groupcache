package lru

import (
	"testing"
	"time"
)

var getExpiredTests = []struct {
	name          string
	keyToAdd      interface{}
	keyToGet      interface{}
	expectedOk    bool
	ttl           time.Duration
	waitBeforeGet int64
}{
	{"hit", "myKey", "myKey", true, 0, 0},
	{"miss", "myKey", "nonsense", false, 0, 0},
	{"miss_expired", "myKey", "myKey", false, 1000 * time.Millisecond, 1100},
}

func TestGetExpired(t *testing.T) {
	for _, tt := range getExpiredTests {
		cache := NewTTL(0, tt.ttl)
		cache.Add(tt.keyToAdd, 1234)

		if tt.waitBeforeGet > 0 {
			time.Sleep(time.Duration(tt.waitBeforeGet * int64(time.Millisecond)))
		}

		val, ok := cache.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestRemoveExpired(t *testing.T) {
	cache := NewTTL(0, 1000*time.Millisecond)
	cache.Add("myKey", 1234)
	if val, ok := cache.Get("myKey"); !ok {
		t.Fatal("TestRemoveExpired returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemoveExpired failed.  Expected %d, got %v", 1234, val)
	}

	time.Sleep(1000 * time.Millisecond)
	cache.RemoveExpired()
	if cache.Len() != 0 {
		t.Fatal("Cache still contains expired entry")
	}
}
