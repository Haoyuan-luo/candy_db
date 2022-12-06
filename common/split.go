package common

type splitter[T any] struct {
	total int
	step  int
	cur   int
	Data  []T
}

type SplitterImpl[T any] interface {
	Get() ([]T, bool)
}

// NewSplitter 获取一个支持并发操作的分片器
func NewSplitter[T any](data []T) SplitterImpl[T] {
	step := 1
	return &splitter[T]{total: len(data), step: step, cur: 0, Data: data}
}

// Get 获取下一个分片
func (i *splitter[T]) Get() (ret []T, ok bool) {
	// 判断还有没有可以分批的数据
	if i.cur >= i.total {
		return nil, false
	}

	// 如果是最后一次分批，那么直接返回剩余的数据
	if i.cur+i.step >= i.total {
		ret = i.Data[i.cur:]
		i.cur = i.total
		return ret, true
	}

	// 如果不是最后一次分批，那么返回当前分批的数据，同时更新当前分批的起始位置
	ret = i.Data[i.cur : i.cur+i.step]
	i.cur += i.step
	i.step += 100
	return ret, true
}
