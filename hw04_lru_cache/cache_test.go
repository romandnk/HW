package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("Set()", func(t *testing.T) {
		c := NewCache(3)

		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3) // 3 2 1

		require.Equal(t, 3, c.(*lruCache).queue.Front().Value)
		require.Equal(t, 2, c.(*lruCache).queue.Back().Prev.Value)
		require.Equal(t, 1, c.(*lruCache).queue.Back().Value)
		require.Equal(t, 3, c.(*lruCache).queue.Len())
	})

	t.Run("Set() with purge", func(t *testing.T) {
		c := NewCache(3)

		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3)     // 3 2 1
		ok := c.Set("one", 0) // 0 3 2
		require.True(t, ok)

		require.Equal(t, 0, c.(*lruCache).queue.Front().Value)
		require.Equal(t, 3, c.(*lruCache).queue.Back().Prev.Value)
		require.Equal(t, 2, c.(*lruCache).queue.Back().Value)
		require.Equal(t, 3, c.(*lruCache).queue.Len())
	})

	t.Run("Set() with purge and offset", func(t *testing.T) {
		c := NewCache(3)

		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3) // 3 2 1
		c.Set("one", 4)
		c.Set("two", 5)
		c.Set("three", 6) // 6 5 4

		require.Equal(t, 6, c.(*lruCache).queue.Front().Value)
		require.Equal(t, 5, c.(*lruCache).queue.Back().Prev.Value)
		require.Equal(t, 4, c.(*lruCache).queue.Back().Value)
		require.Equal(t, 3, c.(*lruCache).queue.Len())
	})

	t.Run("Get()", func(t *testing.T) {
		c := NewCache(3)

		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3) // 3 2 1

		_, ok1 := c.Get("one")
		_, ok2 := c.Get("two")
		_, ok3 := c.Get("three")
		_, ok4 := c.Get("four")
		require.True(t, ok1)
		require.True(t, ok2)
		require.True(t, ok3)
		require.False(t, ok4)
	})

	t.Run("Clear()", func(t *testing.T) {
		c := NewCache(3)

		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3) // 3 2 1

		c.Clear()
		require.Equal(t, 3, c.(*lruCache).capacity)
		require.Equal(t, 0, c.(*lruCache).queue.Len())
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
