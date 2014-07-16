package lru

import "time"

// CacheTTL extends LRU cache with expiration support.
// Now records are wrapped with metadata which contains TTL value.
// Records that are considered as expired will never return to caller,
// they will be discarded instead
type CacheTTL struct {
	*Cache

	ttl time.Duration
}

// NewTTL creates a new TTL Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func NewTTL(maxEntries int, ttl time.Duration) *CacheTTL {
	return &CacheTTL{
		Cache: New(maxEntries),
		ttl:   ttl,
	}
}

// Get looks up a key's value from the cache.
func (c *CacheTTL) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}

	if ele, hit := c.cache[key]; hit {

		val := ele.Value.(*entry)

		if c.isValueExpired(val) {
			c.removeElement(ele)

			return
		}

		c.ll.MoveToFront(ele)
		return val.value, true
	}

	return
}

// RemoveExpired removes expired records from cache
func (c *CacheTTL) RemoveExpired() {

	for _, ele := range c.cache {
		if c.isValueExpired(ele.Value.(*entry)) {
			c.removeElement(ele)
		}
	}

}

// isValueExpired checks if given cache entry is expired according cache's TTL
// If TTL is set to 0, value is never expired
func (c *CacheTTL) isValueExpired(e *entry) bool {
	return c.ttl > 0 && (e.created+int64(c.ttl) < time.Now().UnixNano())
}
