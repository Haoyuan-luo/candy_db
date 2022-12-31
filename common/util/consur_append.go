package util

type AppendService[T any] interface {
	Get() []T
	Close()
	Ready() func(...T)
}

type appendImpl[T any] struct {
	c    chan struct{}
	ch   chan T // 用来 同步的channel
	data []T    // 存储数据的slice
}

func (s *appendImpl[T]) schedule() {
	for i := range s.ch {
		s.data = append(s.data, i)
	}
}

func (s *appendImpl[T]) Close() {
	close(s.ch)
	<-s.c
}

func (s *appendImpl[T]) Ready() func(...T) {
	return func(vs ...T) {
		for _, v := range vs {
			s.ch <- v
		}
	}
}

func (s *appendImpl[T]) Get() []T {
	return s.data
}

func NewAppendImpl[T any](size int) AppendService[T] {
	c := make(chan struct{})
	s := &appendImpl[T]{
		c:    c,
		ch:   make(chan T, size),
		data: make([]T, 0),
	}
	done := func() { c <- struct{}{} }
	go func() {
		s.schedule()
		done()
	}()
	return s
}
