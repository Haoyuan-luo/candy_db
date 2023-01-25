package util

const fanNum = 3

type iterator[T, R any] struct {
	done   chan struct{}
	stream <-chan T
	fan    []<-chan R
}

func (i iterator[T, R]) getStream() <-chan R {
	return Consumer[R](i.done, i.fan...)
}

func fromSlice[T, R any](slice []T) iterator[T, R] {
	done := make(chan struct{})
	return iterator[T, R]{done, Producer[T](done, slice), make([]<-chan R, fanNum)}
}

func forEach[T, R any](iter iterator[T, R], fn func(T) R, scene func(chan R, T, func(T) R)) iterator[T, R] {
	for i := 0; i < fanNum; i++ {
		iter.fan[i] = Processor[T, R](iter.done, iter.stream, fn, scene)
	}
	return iter
}

func work[T, R any](s []T, fn func(T) R, scene func(chan R, T, func(T) R)) <-chan R {
	return forEach[T, R](fromSlice[T, R](s), fn, scene).getStream()
}

type IterRet[T any] <-chan T

// Range 迭代
func (i IterRet[T]) Range() <-chan T {
	return i
}

// ToSlice 转换为切片
func (i IterRet[T]) ToSlice() (ret []T) {
	for v := range i {
		ret = append(ret, v)
	}
	return
}

func Map[T, R any](s []T, fn func(T) R) IterRet[R] {
	scene := func(fan chan R, v T, fn func(T) R) {
		fan <- fn(v)
	}
	return work(s, fn, scene)
}

type judge[T any] interface {
	IsTrue() bool
	Value() T
}

func Filter[T any, R judge[T]](s []R) IterRet[T] {
	scene := func(fan chan T, v R, fn func(R) T) {
		if v.IsTrue() {
			fan <- v.Value()
		}
	}
	return work(s, nil, scene)
}
