package util

import "sync"

type ConcurMapService[T any] interface {
	Load(key string) (T, bool)
	Store(key string, value T)
	Len() int
	Range() map[string]T
}

type concurMapImpl[T any] struct {
	ConcurrencyMap sync.Map
	dirty          map[string]T
}

func NewConcurMapImpl[T any]() ConcurMapService[T] {
	return &concurMapImpl[T]{
		ConcurrencyMap: sync.Map{},
		dirty:          map[string]T{},
	}
}

func (c *concurMapImpl[T]) Range() map[string]T {
	c.convert()
	return c.dirty
}

func (c *concurMapImpl[T]) Store(key string, value T) {
	c.dirty[key] = value
}

func (c *concurMapImpl[T]) Load(key string) (value T, ok bool) {
	c.convert()
	if value, ok = c.dirty[key]; ok {
		return value, true
	}
	return value, false
}

func (c *concurMapImpl[T]) Len() int {
	c.convert()
	return len(c.dirty)
}

func (c *concurMapImpl[T]) convert() {
	c.ConcurrencyMap.Range(func(key, value interface{}) bool {
		c.dirty[key.(string)] = value.(T)
		return true
	})
}
