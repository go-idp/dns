package commands

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-idp/dns/cmd/dns/config"
)

// dnsAnswerCache stores final IP lists for queries that were resolved via upstream
// (including config/system alias chains). Keys are normalized name + query type.
type dnsAnswerCache struct {
	mu         sync.RWMutex
	entries    map[string]*dnsCacheEntry
	maxEntries int
}

type dnsCacheEntry struct {
	ips      []string // only used when negative == false
	expires  time.Time
	negative bool // NXDOMAIN / empty success
}

func dnsCacheKey(hostname string, typ int) string {
	h := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(hostname), "."))
	return h + "#" + strconv.Itoa(typ)
}

func newDNSAnswerCache(maxEntries int) *dnsAnswerCache {
	if maxEntries <= 0 {
		maxEntries = config.DNSCacheMaxEntriesDefault
	}
	return &dnsAnswerCache{
		entries:    make(map[string]*dnsCacheEntry),
		maxEntries: maxEntries,
	}
}

// get returns (ips, true) on hit. For negative cache, ips is empty slice.
func (c *dnsAnswerCache) get(now time.Time, key string) ([]string, bool) {
	if c == nil {
		return nil, false
	}
	c.mu.RLock()
	e := c.entries[key]
	c.mu.RUnlock()
	if e == nil {
		return nil, false
	}
	if now.After(e.expires) {
		c.mu.Lock()
		if cur := c.entries[key]; cur == e && now.After(cur.expires) {
			delete(c.entries, key)
		}
		c.mu.Unlock()
		return nil, false
	}
	if e.negative {
		return []string{}, true
	}
	out := make([]string, len(e.ips))
	copy(out, e.ips)
	return out, true
}

func (c *dnsAnswerCache) set(now time.Time, key string, ips []string, negative bool, ttl time.Duration) {
	if c == nil || ttl <= 0 {
		return
	}
	e := &dnsCacheEntry{
		expires:  now.Add(ttl),
		negative: negative,
	}
	if !negative {
		e.ips = make([]string, len(ips))
		copy(e.ips, ips)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = e
	c.evictIfNeededLocked(now)
}

func (c *dnsAnswerCache) evictIfNeededLocked(now time.Time) {
	for k, e := range c.entries {
		if now.After(e.expires) {
			delete(c.entries, k)
		}
	}
	if len(c.entries) <= c.maxEntries {
		return
	}
	over := len(c.entries) - c.maxEntries
	for k := range c.entries {
		if over <= 0 {
			break
		}
		delete(c.entries, k)
		over--
	}
}
