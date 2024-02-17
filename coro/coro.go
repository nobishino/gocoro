package coro

func New[In, Out any](f func(In, func(Out) In) Out) (resume func(In) (Out, bool)) {
	cin := make(chan In)
	cout := make(chan Out)
	running := true
	resume = func(in In) (out Out, ok bool) {
		if !running {
			return
		}
		cin <- in
		out = <-cout
		return out, running
	}
	yield := func(out Out) In {
		cout <- out
		return <-cin
	}
	go func() {
		out := f(<-cin, yield)
		running = false
		cout <- out
	}()
	return resume
}

func Pull[V any](push func(yield func(V) bool)) (pull func() (V, bool), stop func()) {
	copush := func(more bool, yield func(V) bool) V {
		if more {
			push(yield)
		}
		var zero V
		return zero
	}
	resume := New(copush)
	pull = func() (V, bool) {
		return resume(true)
	}
	stop = func() {
		resume(false)
	}
	return pull, stop
}
