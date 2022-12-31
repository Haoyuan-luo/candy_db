package client

import (
	"candy_db/common"
	"candy_db/common/util"
	"unsafe"
)

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

type CandyDBClient[T any] struct {
	common.SkipListImpl
}

type Optional[T any] struct {
	Value []byte
}

func (c Optional[T]) Assert() (value T, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	return util.GenDirect(*(**T)(unsafe.Pointer(&c.Value))), true
}

type CandyDBService[T any] interface {
	Add(key []byte, value T)
	Find(key []byte) Optional[T]
}

func NewCandyDBClient[T any]() CandyDBService[T] {
	return &CandyDBClient[T]{
		SkipListImpl: common.NewSkipList(),
	}
}

func (c CandyDBClient[T]) Add(key []byte, value T) {
	Len := unsafe.Sizeof(value)
	data := *(*[]byte)(unsafe.Pointer(&SliceMock{
		addr: uintptr(unsafe.Pointer(&value)),
		cap:  int(Len),
		len:  int(Len),
	}))
	c.AddNode(common.NewContainer(key, data))
}
func (c CandyDBClient[T]) Find(key []byte) Optional[T] {
	container := common.NewContainer(key, nil)
	c.FindNode(container)
	return Optional[T]{container.GetData()}
}
