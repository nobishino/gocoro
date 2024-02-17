package iter_test

import (
	"testing"

	"github.com/nobishino/gocoro/iter"
)

func TestPull(t *testing.T) {
	next, stop := iter.Pull(func(yield func(int) bool) {
		for i := 0; i < 3; i++ {
			if !yield(i) {
				return
			}
		}
	})
	defer stop()
	var i int
	for v, ok := next(); ok; v, ok = next() {
		if v != i {
			t.Errorf("got %d, want %d", v, i)
		}
		if !ok {
			t.Errorf("got %t, want %t", ok, true)
		}
		i++
	}
}

func TestPull_Merge(t *testing.T) {

}

func TestPull2(t *testing.T) {
	ss := []string{"hello", "world", "!"}
	next, stop := iter.Pull2(func(yield func(int, string) bool) {
		for i := 0; i < 3; i++ {
			if !yield(i, ss[i]) {
				return
			}
		}
	})
	defer stop()
	var i int
	for k, v, ok := next(); ok; k, v, ok = next() {
		if k != i || v != ss[i] {
			t.Errorf("got (%d, %q), want (%d, %q)", k, v, i, ss[i])
		}
		if !ok {
			t.Errorf("got %t, want %t", ok, true)
		}
		i++
	}
}
