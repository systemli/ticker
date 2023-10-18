package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	interval := 100 * time.Microsecond
	c := NewCache(interval)
	defer c.Close()

	c.Set("foo", "bar", 0)
	c.Set("baz", "qux", interval/2)

	baz, found := c.Get("baz")
	assert.True(t, found)
	assert.Equal(t, "qux", baz)

	time.Sleep(interval / 2)

	_, found = c.Get("baz")
	assert.False(t, found)

	time.Sleep(interval)

	_, found = c.Get("404")
	assert.False(t, found)

	foo, found := c.Get("foo")
	assert.True(t, found)
	assert.Equal(t, "bar", foo)
}

func TestDelete(t *testing.T) {
	c := NewCache(time.Minute)
	c.Set("foo", "bar", time.Hour)

	_, found := c.Get("foo")
	assert.True(t, found)

	c.Delete("foo")

	_, found = c.Get("foo")
	assert.False(t, found)
}

func TestRange(t *testing.T) {
	c := NewCache(time.Minute)
	c.Set("foo", "bar", time.Hour)
	c.Set("baz", "qux", time.Hour)

	count := 0
	c.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	assert.Equal(t, 2, count)
}

func TestRangeTimer(t *testing.T) {
	c := NewCache(time.Minute)
	c.Set("foo", "bar", time.Nanosecond)
	c.Set("baz", "qux", time.Nanosecond)

	time.Sleep(time.Microsecond)

	c.Range(func(key, value interface{}) bool {
		assert.Fail(t, "should not be called")
		return true
	})
}

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NewCache(5 * time.Second).Close()
		}
	})
}

func BenchmarkGet(b *testing.B) {
	c := NewCache(5 * time.Second)
	defer c.Close()

	c.Set("foo", "bar", 0)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get("foo")
		}
	})
}

func BenchmarkSet(b *testing.B) {
	c := NewCache(5 * time.Second)
	defer c.Close()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set("foo", "bar", 0)
		}
	})
}

func BenchmarkDelete(b *testing.B) {
	c := NewCache(5 * time.Second)
	defer c.Close()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Delete("foo")
		}
	})
}
