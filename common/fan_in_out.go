package common

func Producer[R any](done chan struct{}, integers []R) <-chan R {
	stream := make(chan R)
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
	out := make(chan T)
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
	var out []T
	for v := range stream {
		select {
		case <-done:
			return out
		default:
			out = append(out, v)
		}
	}
	return out
}
