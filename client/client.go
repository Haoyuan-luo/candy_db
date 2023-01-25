package client

import (
	"candy_db/common"
	"candy_db/common/util"
	"context"
	"github.com/pkg/errors"
	"unsafe"
)

type Option[K comparable, V any] func(client *candyDBClient[K, V]) error

type CandyResult[V any] struct {
	Value V
	raw   []byte
	err   error
}

func (o *CandyResult[V]) assert() (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	if o.err != nil {
		return false
	}
	o.Value = util.GenDirect(*(**V)(unsafe.Pointer(&o.raw)))
	return true
}

type candyDBClient[K comparable, V any] struct {
	memTable   common.MemService[K]
	immutables []common.MemService[K]
	after      []Callback
	before     []Callback
}

type CandyDBService[K comparable, V any] interface {
	Add(ctx context.Context, key K, value V) error
	Find(ctx context.Context, key K) (CandyResult[V], error)
}

func NewCandyDBClient[K comparable, V any](opt ...Option[K, V]) (CandyDBService[K, V], error) {
	// 创建一个带默认配置的CandyDBClient
	mem, _ := common.NewMemTable[K]()
	client := &candyDBClient[K, V]{memTable: mem}
	// 通过option对client进行配置
	for _, o := range opt {
		if err := o(client); err != nil { // 增加filed
			return nil, err
		}
	}
	return client, nil
}

type sliceMock struct {
	addr uintptr
	len  int
	cap  int
}

func (c *candyDBClient[K, V]) Add(ctx context.Context, key K, value V) error {
	// 异常处理
	defer func() {
		if r := recover(); r != nil {
			// Todo something example insert wal or print log
		}
	}()

	// 前置中间件
	if err := callOpt(ctx, CallReq{key: key, value: value}, c.before); err != nil {
		return err
	}

	Len := unsafe.Sizeof(value)
	data := *(*[]byte)(unsafe.Pointer(&sliceMock{
		addr: uintptr(unsafe.Pointer(&value)),
		cap:  int(Len),
		len:  int(Len),
	}))

	if err := c.memTable.Add(ctx, key, data); err != nil {
		return err
	}

	// 后置中间件
	if err := callOpt(ctx, CallReq{key: key, value: value}, c.after); err != nil {
		return err
	}
	return nil
}

func (c *candyDBClient[K, V]) Find(ctx context.Context, key K) (candyRet CandyResult[V], err error) {

	// 异常处理
	defer func() {
		if r := recover(); r != nil {
			// Todo something example insert wal or print log
		}
	}()

	// 前置中间件
	if err = callOpt(ctx, CallReq{key: key}, c.before); err != nil {
		return candyRet, err
	}

	container, err := c.memTable.Find(ctx, key)
	if err != nil {
		return candyRet, err
	}
	candyRet.raw = container.GetData()

	// 后置中间件
	if err = callOpt(ctx, CallReq{key: key}, c.after); err != nil {
		return candyRet, err
	}

	if ok := candyRet.assert(); !ok {
		return candyRet, errors.New("assert error")
	} else {
		return candyRet, nil
	}
}
