package common

import "sync"

type Iterator[T any] struct {
	done   chan struct{}
	stream <-chan T
}

func (i Iterator[T]) GetStream() []T {
	return Consumer[T](i.done, i.stream)
}

func FromSlice[T any](slice []T) Iterator[T] {
	done := make(chan struct{})
	return Iterator[T]{done, Producer[T](done, slice)}
}

func ForEach[T, R any](iter Iterator[T], fn func(T) R) Iterator[R] {
	return Iterator[R]{iter.done, Processor[T, R](iter.done, iter.stream, fn)}
}

func Serial[T, R any](s []T, fn func(T) R) (ret []R) {
	ret = make([]R, len(s))
	sp := NewSplitter(s)
	for {
		got, ok := sp.Get()
		if !ok {
			break
		}
		copy(ret, ForEach(FromSlice(got), fn).GetStream())
	}
	return ret
}

func Parallel[T, R any](s []T, fn func(T) R) (ret []R) {
	ret = make([]R, len(s))
	wg := &sync.WaitGroup{}
	sp := NewSplitter(s)
	for {
		got, ok := sp.Get()
		if !ok {
			break
		}
		wg.Add(1)
		go func() {
			copy(ret, ForEach(FromSlice(got), fn).GetStream())
			wg.Done()
		}()
	}
	wg.Wait()
	return ret
}

func Map[T, R any](s []T, fn func(T) R) []R {
	return Parallel(s, fn)
}

func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	filter := func(t T) struct{} {
		if f(t) {
			r = append(r, t)
		}
		return struct{}{}
	}
	Serial(s, filter)
	return r
}
