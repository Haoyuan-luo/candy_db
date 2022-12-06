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

func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	filter := func(t T) struct{} {
		if f(t) {
			r = append(r, t)
		}
		return struct{}{}
	}
	ForEach(FromSlice(s), filter).GetStream()
	return r
}

func FilterV2[T any](s []T, f func(T) bool) []T {
	wg := sync.WaitGroup{}
	r := make([]T, 0, len(s))
	filter := func(t T) struct{} {
		if f(t) {
			r = append(r, t)
		}
		return struct{}{}
	}
	for sp := NewSplitter(s); sp.HasNext(); sp.Next() {
		wg.Add(1)
		go func() {
			ForEach(FromSlice(sp.Get()), filter)
			wg.Done()
		}()
	}
	wg.Wait()
	return r
}
