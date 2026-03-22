package commands

import (
	"testing"
	"time"
)

func TestDNSCacheKey(t *testing.T) {
	t.Parallel()
	if dnsCacheKey("Example.COM.", 4) != "example.com#4" {
		t.Fatal(dnsCacheKey("Example.COM.", 4))
	}
	if dnsCacheKey("a.b", 6) != "a.b#6" {
		t.Fatal(dnsCacheKey("a.b", 6))
	}
}

func TestDNSAnswerCachePositiveNegative(t *testing.T) {
	t.Parallel()
	c := newDNSAnswerCache(10)
	now := time.Now()
	key := "x#4"

	if _, hit := c.get(now, key); hit {
		t.Fatal("unexpected hit")
	}

	c.set(now, key, []string{"1.1.1.1", "2.2.2.2"}, false, time.Minute)
	got, hit := c.get(now, key)
	if !hit || len(got) != 2 || got[0] != "1.1.1.1" {
		t.Fatalf("got %v hit=%v", got, hit)
	}

	c.set(now, key, nil, true, time.Minute)
	got2, hit2 := c.get(now, key)
	if !hit2 || len(got2) != 0 {
		t.Fatalf("negative got %v hit=%v", got2, hit2)
	}
}

func TestDNSAnswerCacheExpiry(t *testing.T) {
	t.Parallel()
	c := newDNSAnswerCache(10)
	now := time.Now()
	key := "y#4"
	c.set(now, key, []string{"9.9.9.9"}, false, 50*time.Millisecond)
	if _, hit := c.get(now.Add(100 * time.Millisecond), key); hit {
		t.Fatal("expected miss after expiry")
	}
}
