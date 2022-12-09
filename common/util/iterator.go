package util

const fanNum = 3

//const (
//	defaultScene string = "default"
//	filterScene  string = "filter"
//)
//
//func filter[T, R any](f func(T) R, in <-chan T, out chan<- R) {
//	for v := range in {
//		out <- f(v)
//	}
//}
//
//type sceneMap[T, R any] map[string]sceneType[any, any]

type Iterator[T, R any] struct {
	done   chan struct{}
	stream <-chan T
	fan    []<-chan R
}

func (i Iterator[T, R]) GetStream() <-chan R {
	return Consumer[R](i.done, i.fan...)
}

func FromSlice[T, R any](slice []T) Iterator[T, R] {
	done := make(chan struct{})
	return Iterator[T, R]{done, Producer[T](done, slice), make([]<-chan R, fanNum)}
}

func FromSliceChan[T, R any](slice <-chan T) Iterator[T, R] {
	done := make(chan struct{})
	return Iterator[T, R]{done, ProducerChan[T](done, slice), make([]<-chan R, fanNum)}
}

func ForEach[T, R any](iter Iterator[T, R], fn func(T) R, scene func(chan R, T, func(T) R)) Iterator[T, R] {
	for i := 0; i < fanNum; i++ {
		iter.fan[i] = Processor[T, R](iter.done, iter.stream, fn, scene)
	}
	return iter
}

func Serial[T, R any](s []T, sc <-chan T, fn func(T) R, scene func(chan R, T, func(T) R), isChan bool) <-chan R {
	if isChan {
		return ForEach[T, R](FromSliceChan[T, R](sc), fn, scene).GetStream()
	}
	return ForEach[T, R](FromSlice[T, R](s), fn, scene).GetStream()
}

func Map[T, R any](s []T, fn func(T) R) <-chan R {
	scene := func(fan chan R, v T, fn func(T) R) {
		fan <- fn(v)
	}
	return Serial(s, nil, fn, scene, false)
}

func MapChan[T, R any](s <-chan T, fn func(T) R) <-chan R {
	scene := func(fan chan R, v T, fn func(T) R) {
		fan <- fn(v)
	}
	return Serial(nil, s, fn, scene, true)
}

type judge[R any] interface {
	IsTrue() bool
	Value() R
}

func Filter[R any, T judge[R]](s []T) <-chan R {
	scene := func(fan chan R, v T, fn func(T) R) {
		if v.IsTrue() {
			fan <- v.Value()
		}
	}
	return Serial(s, nil, nil, scene, false)
}

func FilterChan[R any, T judge[R]](s <-chan T) <-chan R {
	scene := func(fan chan R, v T, fn func(T) R) {
		if v.IsTrue() {
			fan <- v.Value()
		}
	}
	return Serial(nil, s, nil, scene, true)
}
