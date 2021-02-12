package cache

import (
	"testing"
	"time"

	"github.com/x893675/gopkg/clock"
)

func expectEntry(t *testing.T, c *LRUExpireCache, key interface{}, value interface{}) {
	result, ok := c.Get(key)
	if !ok || result != value {
		t.Errorf("Expected cache[%v]: %v, got %v", key, value, result)
	}
}

func expectNotEntry(t *testing.T, c *LRUExpireCache, key interface{}) {
	if result, ok := c.Get(key); ok {
		t.Errorf("Expected cache[%v] to be empty, got %v", key, result)
	}
}

func TestSimpleGet(t *testing.T) {
	c := NewLRUExpireCache(10)
	c.Add("long-lived", "12345", 10*time.Hour)
	expectEntry(t, c, "long-lived", "12345")
}

func TestExpiredGet(t *testing.T) {
	fakeClock := clock.NewFakeClock(time.Now())
	c := NewLRUExpireCacheWithClock(10, fakeClock)
	c.Add("short-lived", "12345", 1*time.Millisecond)
	// ensure the entry expired
	fakeClock.Step(2 * time.Millisecond)
	expectNotEntry(t, c, "short-lived")
}

func TestLRUOverflow(t *testing.T) {
	c := NewLRUExpireCache(4)
	c.Add("elem1", "1", 10*time.Hour)
	c.Add("elem2", "2", 10*time.Hour)
	c.Add("elem3", "3", 10*time.Hour)
	c.Add("elem4", "4", 10*time.Hour)
	c.Add("elem5", "5", 10*time.Hour)
	expectNotEntry(t, c, "elem1")
	expectEntry(t, c, "elem2", "2")
	expectEntry(t, c, "elem3", "3")
	expectEntry(t, c, "elem4", "4")
	expectEntry(t, c, "elem5", "5")
}
