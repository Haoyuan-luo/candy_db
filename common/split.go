package common

const batch = 1

type splitter[T any] struct {
	total int
	step  int
	cur   int
	final bool
	Data  []T
}

type SplitterImpl[T any] interface {
	Get() []T
	HasNext() (ok bool)
	Next()
}

func NewSplitter[T any](data []T) SplitterImpl[T] {
	step := len(data) / batch
	return &splitter[T]{total: len(data), step: step, cur: 0, Data: data}
}

func (i *splitter[T]) Get() []T {
	if i.final {
		re := i.Data[i.cur:]
		i.cur = i.total
		return re
	} else {
		return i.Data[i.Start():i.End()]
	}
}

func (i *splitter[T]) HasNext() (ok bool) {
	return i.cur < i.total
}

func (i *splitter[T]) Start() int {
	return i.cur
}

func (i *splitter[T]) End() int {
	end := i.cur + i.step
	if end > i.total {
		end = i.total
	}
	return end
}

func (i *splitter[T]) Next() {
	if !i.HasNext() {
		return
	}
	if i.cur+i.step > i.total {
		i.final = true
		return
	}
	i.cur += i.step
}
