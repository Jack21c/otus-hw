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

	t.Run("custom", func(t *testing.T) {
		l := NewList()

		l.PushBack(20)  // [20]
		l.PushFront(10) // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		last := l.Back() // 30
		l.Remove(last)   // [10, 20]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushBack(v)
			} else {
				l.PushFront(v)
			}
		} // [70, 50, 10, 20, 40, 60, 80]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 70, l.Front().Value)
		require.Equal(t, 80, l.Back().Value)

		preLast := l.Back().Prev // 60
		second := l.Front().Next // 50
		l.MoveToFront(preLast)   // [60, 70, 50, 10, 20, 40, 80]
		l.MoveToFront(second)    // [50, 60, 70, 10, 20, 40, 80]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{50, 60, 70, 10, 20, 40, 80}, elems)
	})

	t.Run("one_element", func(t *testing.T) {
		l := NewList()

		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		value := "test"
		el := l.PushFront(value)
		require.Equal(t, value, l.Front().Value)
		require.Equal(t, value, l.Back().Value)

		l.Remove(el)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
}
