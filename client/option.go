package client

import (
	"candy_db/common"
	"candy_db/common/util"
	"context"
	"github.com/hashicorp/go-multierror"
)

type Callback func(ctx context.Context, req CallReq) error

type CallReq struct {
	key   interface{}
	value interface{}
}

// CallOpt 中间件处理
func callOpt(ctx context.Context, req CallReq, fun []Callback) (multiErr error) {
	return multierror.Append(multiErr, util.Map[Callback, error](fun, func(subFun Callback) error {
		return subFun(ctx, req)
	}).ToSlice()...).ErrorOrNil()
}

// WithBeforeCall 增加一个前置中间件
func WithBeforeCall[K comparable, V any](call ...Callback) Option[K, V] {
	return func(client *candyDBClient[K, V]) error {
		client.before = append(client.before, call...)
		return nil
	}
}

// WithAfterCall 增加一个后置中间件
func WithAfterCall[K comparable, V any](call ...Callback) Option[K, V] {
	return func(client *candyDBClient[K, V]) error {
		client.after = append(client.after, call...)
		return nil
	}
}

// WithMemTable 根据配置创建一个MemTable
func WithMemTable[K comparable, V any](config ...common.MemConfig[K]) Option[K, V] {
	return func(client *candyDBClient[K, V]) (err error) {
		client.memTable, err = common.NewMemTable[K](config...)
		return nil
	}
}
