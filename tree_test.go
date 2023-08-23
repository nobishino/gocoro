package coro_test

import (
	"fmt"
	"testing"

	"github.com/nobishino/coro"
)

var (
	e  = (*Tree[int])(nil)
	t1 = T(T(T(e, 1, e), 2, T(e, 3, e)), 4, T(e, 5, e))
	t2 = T(e, 1, T(e, 2, T(e, 3, T(e, 4, T(e, 5, e)))))
	t3 = T(e, 1, T(e, 2, T(e, 3, T(e, 4, T(e, 6, e)))))
)

func TestTreeCmp(t *testing.T) {
	fmt.Println(cmp(t1, t2), cmp(t1, t3))
}

func TestTreeCmpWithoutPull(t *testing.T) {
	fmt.Println(
		cmp_without_pull(t1, t2),
		cmp_without_pull(t1, t3),
	)

}

type Tree[V any] struct {
	Left  *Tree[V]
	Value V
	Right *Tree[V]
}

func T[V any](l *Tree[V], v V, r *Tree[V]) *Tree[V] {
	return &Tree[V]{Left: l, Value: v, Right: r}
}

func (t *Tree[V]) All(yield func(v V) bool) {
	t.all(yield)
}

func (t *Tree[V]) all(yield func(v V) bool) bool {
	return t == nil ||
		t.Left.all(yield) && yield(t.Value) && t.Right.all(yield)
}

// Pullを利用した場合の比較コード
func cmp[V comparable](t1, t2 *Tree[V]) bool {
	next1, stop1 := coro.Pull(t1.All)
	next2, stop2 := coro.Pull(t2.All)
	defer stop1()
	defer stop2()
	for {
		v1, ok1 := next1()
		v2, ok2 := next2()
		if v1 != v2 || ok1 != ok2 {
			return false
		}
		if !ok1 && !ok2 {
			return true
		}
	}
}

// Pullを利用しない場合の比較コード
func cmp_without_pull[V comparable](t1, t2 *Tree[V]) bool {
	resume1 := coro.New(func(more bool, yield func(V) bool) (zero V) {
		if !more {
			return
		}
		t1.All(yield)
		return zero
	})
	resume2 := coro.New(func(more bool, yield func(V) bool) (zero V) {
		if !more {
			return
		}
		t2.All(yield)
		return zero
	})
	// この処理を省いた場合、coroutineがyieldの戻り値を待ってブロックした状態のままleakする
	defer func() {
		z1, ok1 := resume1(false)
		z2, ok2 := resume1(false)
		fmt.Println("Comparison finished")
		fmt.Println(z1, ok1)
		fmt.Println(z2, ok2)
	}()
	for {
		v1, ok1 := resume1(true)
		v2, ok2 := resume2(true)
		if v1 != v2 || ok1 != ok2 {
			return false
		}
		if !ok1 && !ok2 {
			return true
		}
	}

}
