package util

type AppendService[T any] struct {
	c    chan struct{}
	ch   chan T // 用来 同步的channel
	data []T    // 存储数据的slice
}

func (s *AppendService[T]) Schedule() {
	// 从 channel 接收数据
	for i := range s.ch {
		s.data = append(s.data, i)
	}
}

func (s *AppendService[T]) Close() {
	// 最后关闭 channel
	close(s.ch)
	<-s.c
}

func (s *AppendService[T]) Ready() func(...T) {
	return func(vs ...T) {
		for _, v := range vs {
			s.ch <- v
		}
	}
}

func NewAppendService[T any](size int) *AppendService[T] {
	c := make(chan struct{})
	s := &AppendService[T]{
		c:    c,
		ch:   make(chan T, size),
		data: make([]T, 0),
	}
	done := func() { c <- struct{}{} }
	go func() {
		// 并发地 append 数据到 slice
		s.Schedule()
		done()
	}()
	return s
}
