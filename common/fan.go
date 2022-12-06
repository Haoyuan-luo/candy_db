package common

func Producer[R any](done chan struct{}, integers []R) <-chan R {
	stream := make(chan R, cap(integers))
	go func() {
		defer close(stream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case stream <- i:
			}
		}
	}()
	return stream
}

func Processor[R, T any](done chan struct{}, stream <-chan R, fn func(R) T) <-chan T {
	out := make(chan T, cap(stream))
	go func() {
		defer close(out)
		for v := range stream {
			select {
			case <-done:
				return
			case out <- fn(v):
			}
		}
	}()
	return out
}

func Consumer[T any](done chan struct{}, stream <-chan T) []T {
	var ret []T
	defer close(done)
	for v := range stream {
		select {
		case <-done:
			return ret
		default:
			ret = append(ret, v)
		}
	}
	return ret
}
