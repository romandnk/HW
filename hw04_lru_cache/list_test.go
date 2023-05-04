package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("Front() and Back()", func(t *testing.T) {
		l := NewList()

		l.PushFront(1)
		l.PushBack(2)
		l.PushFront(3)
		l.PushBack(4)
		l.PushFront(5) // 5 3 1 2 4

		require.Equal(t, 5, l.Front().Value)
		require.Equal(t, 4, l.Back().Value)
	})

	t.Run("PushFront() and PushBack()", func(t *testing.T) {
		l := NewList()

		l.PushFront(1)
		l.PushFront(2)
		l.PushBack(3)
		l.PushBack(4)
		l.PushFront(5)
		l.PushBack(6) // 5 2 1 3 4 6

		require.Equal(t, 5, l.Front().Value)
		require.Equal(t, 2, l.Front().Next.Value)
		require.Equal(t, 1, l.Front().Next.Next.Value)
		require.Equal(t, 3, l.Front().Next.Next.Next.Value)
		require.Equal(t, 4, l.Front().Next.Next.Next.Next.Value)
		require.Equal(t, 6, l.Back().Value)
	})

	t.Run("Remove()", func(t *testing.T) {
		l := NewList()

		l.PushFront(1)
		l.PushFront(2)
		l.PushBack(3)
		l.PushBack(4)
		l.PushFront(5)
		l.PushBack(6) // 5 2 1 3 4 6

		l.Remove(l.Front())          // 2 1 3 4 6
		l.Remove(l.Back().Prev.Prev) // 2 1 4 6
		l.Remove(l.Front())          // 1 4 6
		l.Remove(l.Back())           // 1 4

		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 4, l.Front().Next.Value)
		require.Nil(t, l.Back().Next)
		require.Nil(t, l.Front().Prev)
	})

	t.Run("MoveToFront()", func(t *testing.T) {
		l := NewList()

		l.PushFront(1)
		l.PushFront(2)
		l.PushBack(3)
		l.PushFront(5)
		l.PushBack(6) // 5 2 1 3 6

		l.MoveToFront(l.Front())          // 5 2 1 3 6
		l.MoveToFront(l.Front().Next)     // 2 5 1 3 6
		l.MoveToFront(l.Back())           // 6 2 5 1 3
		l.MoveToFront(l.Back().Prev.Prev) // 5 6 2 1 3

		require.Equal(t, 5, l.Front().Value)
		require.Equal(t, 6, l.Front().Next.Value)
		require.Equal(t, 2, l.Front().Next.Next.Value)
		require.Equal(t, 1, l.Back().Prev.Value)
		require.Equal(t, 3, l.Back().Value)
		require.Nil(t, l.Back().Next)
		require.Nil(t, l.Front().Prev)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
