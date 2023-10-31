package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}

func (s *CacheTestSuite) TestCache() {
	interval := 100 * time.Microsecond
	c := NewCache(interval)
	defer c.Close()

	s.Run("Set", func() {
		s.Run("adds value to cache", func() {
			c.Set("foo", "bar", 0)
			foo, found := c.Get("foo")
			s.True(found)
			s.Equal("bar", foo)
		})

		s.Run("add value to cache with expiration", func() {
			c.Set("foo", "bar", interval/2)
			foo, found := c.Get("foo")
			s.True(found)
			s.Equal("bar", foo)

			time.Sleep(interval)

			foo, found = c.Get("foo")
			s.False(found)
			s.Empty(foo)
		})
	})

	s.Run("Get", func() {
		s.Run("returns empty value if not found", func() {
			foo, found := c.Get("foo")
			s.False(found)
			s.Empty(foo)
		})

		s.Run("returns empty value if expired", func() {
			c.Set("foo", "bar", interval/2)
			foo, found := c.Get("foo")
			s.True(found)
			s.Equal("bar", foo)

			time.Sleep(interval)

			foo, found = c.Get("foo")
			s.False(found)
			s.Empty(foo)
		})
	})

	s.Run("Delete", func() {
		s.Run("removes value from cache", func() {
			c.Set("foo", "bar", 0)
			foo, found := c.Get("foo")
			s.True(found)
			s.Equal("bar", foo)

			c.Delete("foo")
			foo, found = c.Get("foo")
			s.False(found)
			s.Empty(foo)
		})
	})

	s.Run("Range", func() {
		s.Run("iterates over all values in cache", func() {
			c.Set("foo", "bar", 0)
			c.Set("bar", "baz", 0)

			count := 0
			c.Range(func(key, value interface{}) bool {
				count++
				return true
			})

			s.Equal(2, count)
		})

		s.Run("iterates not over expired values", func() {
			c.Set("foo", "bar", interval/2)
			c.Set("bar", "baz", 0)

			time.Sleep(interval)

			count := 0
			c.Range(func(key, value interface{}) bool {
				count++
				return true
			})

			s.Equal(1, count)
		})
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
