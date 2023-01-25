package common

import (
	"candy_db/common/util"
	"context"
)

// 支持多种混合配置的MemTable
type memTable[K comparable] struct {
	skipList  skipListImpl             // 高性能查找和保存的工具
	bloom     BloomService             // 对STTable进行过滤
	cache     util.CacheService[K]     // 缓存工具
	serialize util.SerializeService[K] // 序列化工具
	log       util.LogImpl             // 日志工具
}

type MemConfig[K comparable] func(memTable *memTable[K]) error

func CustomSkipList[K comparable](skipList skipListImpl) MemConfig[K] {
	return func(memTable *memTable[K]) error {
		memTable.skipList = skipList
		return nil
	}
}

func CustomBloom[K comparable](bloom BloomService) MemConfig[K] {
	return func(memTable *memTable[K]) error {
		memTable.bloom = bloom
		return nil
	}
}

func CustomCache[K comparable](cache util.CacheService[K]) MemConfig[K] {
	return func(memTable *memTable[K]) error {
		memTable.cache = cache
		return nil
	}
}

func CustomSerialize[K comparable](serialize util.SerializeService[K]) MemConfig[K] {
	return func(memTable *memTable[K]) error {
		memTable.serialize = serialize
		return nil
	}
}

type MemService[K any] interface {
	Add(ctx context.Context, key K, value []byte) error
	Find(ctx context.Context, key K) (*Container, error)
}

func NewMemTable[K comparable](config ...MemConfig[K]) (MemService[K], error) {
	// 创建一个带有默认配置的memTable
	mem := &memTable[K]{
		skipList:  newSkipList(),
		bloom:     NewBloomFilter(),
		cache:     util.NewCacheService[K](""),
		serialize: util.NewSerializeService[K](),
		log:       util.Logger().SetField("memTable"),
	}
	// 通过config对memTable进行配置
	for _, c := range config {
		if err := c(mem); err != nil {
			return nil, err
		}
	}
	return mem, nil
}

func (m *memTable[K]) Add(ctx context.Context, key K, value []byte) error {
	byteK, err := m.serialize.Serialize(key)
	if err != nil {
		return err
	}
	m.skipList.AddNode(NewContainer(byteK, value))
	if m.cache != nil {
		m.cache.Set(key, value)
	}
	return nil
}

func (m *memTable[K]) Find(ctx context.Context, key K) (*Container, error) {
	if m.cache != nil {
		m.log.Info("cache hit", key)
		if v, ok := m.cache.Get(key); ok {
			return NewContainer(nil, v), nil
		}
	}
	byteK, err := m.serialize.Serialize(key)
	if err != nil {
		return nil, err
	}
	container := NewContainer(byteK, nil)
	m.skipList.FindNode(container)
	return container, nil
}
