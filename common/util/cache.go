package util

const (
	LRU = "lru"
	LFU = "lfu"
	//TinyLfu       = "tiny-lfu"
	//WindowTinyLfu = "window-tiny-lfu"
)

type CacheService[T comparable] interface {
	Get(key T) ([]byte, bool)
	Set(key T, value []byte)
}

// cache工具类
type capacity[T comparable] struct {
	key  T
	numH int // 命中次数
}

type capl[T comparable] []*capacity[T]

func (cl capl[T]) cap() int {
	return cap(cl)
}

func (cl capl[T]) head() *capacity[T] {
	return cl[0]
}

func (cl capl[T]) tail() (ret *capacity[T]) {
	for i := len(cl) - 1; i >= 0; i-- {
		if cl[i] != nil {
			ret = cl[i]
			cl[i] = nil
			return
		}
	}
	return nil
}

func (cl capl[T]) get(key T) *capacity[T] {
	for _, c := range cl {
		if c != nil && c.key == key {
			return c
		}
	}
	return nil
}

func (cl capl[T]) getPos(key T) int {
	for i, c := range cl {
		if c != nil && c.key == key {
			return i
		}
	}
	return -1
}

func (cl capl[T]) getPosWithH(key T) int {
	for i, c := range cl {
		if c != nil && c.key == key {
			cl[i].numH++
			return i
		}
	}
	return -1
}

func (cl capl[T]) swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

func (cl capl[T]) rotate(i int) {
	for j := i; j > 0; j-- {
		cl.swap(j, j-1)
	}
}

func (cl capl[T]) rotateWithH(i int) {
	for j := i; j > 0; j-- {
		if cl[j].numH >= cl[j-1].numH {
			cl.swap(j, j-1)
		} else {
			break
		}
	}
}

func (cl capl[T]) set(key T) {
	for i := 0; i < len(cl); i++ {
		if cl[i] == nil {
			cl[i] = &capacity[T]{key: key, numH: 0}
			return
		}
	}
	cl[len(cl)-1] = &capacity[T]{key: key, numH: 0}
}

// 流式计数方式
// 一条缓存数据的访问次数真的需要int类型这么大的表示范围来统计吗？
// 一个缓存被访问15次已经算是很高的频率了，那么只用4个Bit就可以保存这个数据
// 也就是四个bit 2^4=16 = 0000~1111 正好是半个uint8的长度，那一个uint8就可以储存两个位图
//type cmkRow []byte // 位图（计数基本单位）
//
//func newCmkRow(numContainer int64) cmkRow {
//	return make(cmkRow, numContainer/2) // 除2是因为一个uint8可以存储两个位图
//}
//
//func (cm cmkRow) increment(n uint64) {
//	//定位到第i个Counter
//	i := n / 2 //r[i]
//	//右移距离，偶数为0，奇数为4
//	s := (n & 1) * 4
//	//取前4Bit还是后4Bit
//	v := (cm[i] >> s) & 0x0f //0000, 1111
//	//没有超出最大计数时，计数+1
//	if v < 15 { // 一个8bit最多表示15次
//		cm[i] += uint8(1)
//		cm[i] <<= s
//	}
//}
//
//func (cm cmkRow) get(n uint64) byte {
//	return byte(cm[n/2]>>((n&1)*4)) & 0x0f
//}

// 各种cache的实现

type lruCache[T comparable] struct {
	container map[T][]byte
	capl      capl[T]
}

func (l lruCache[T]) Get(key T) ([]byte, bool) {
	value, ok := l.container[key]
	if ok {
		l.capl.rotate(l.capl.getPos(key))
	}
	return value, ok
}

func (l lruCache[T]) Set(key T, value []byte) {
	if _, ok := l.container[key]; ok { // 如果存在，更新
		l.container[key] = value
		return
	} else {
		if len(l.container) < l.capl.cap() { // 如果不存在，且容量未满，直接插入
			l.container[key] = value
			l.capl.set(key)
		} else { // 如果不存在，且容量已满，删除最近最少使用的
			delete(l.container, l.capl.tail().key)
			l.container[key] = value
			l.capl.set(key)
		}
	}
}

type lfuCache[T comparable] struct {
	container map[T][]byte
	capl      capl[T]
}

func (l lfuCache[T]) Get(key T) ([]byte, bool) {
	value, ok := l.container[key]
	if ok {
		l.capl.rotateWithH(l.capl.getPosWithH(key))
	}
	return value, ok
}

func (l lfuCache[T]) Set(key T, value []byte) {
	if _, ok := l.container[key]; ok { // 如果存在，更新
		l.container[key] = value
		return
	} else {
		if len(l.container) < l.capl.cap() { // 如果不存在，且容量未满，直接插入
			l.container[key] = value
			l.capl.set(key)
		} else { // 如果不存在，且容量已满，删除最近最少使用的
			delete(l.container, l.capl.tail().key)
			l.container[key] = value
			l.capl.set(key)
		}
	}
}

//type tinyLfuCache struct {
//	container map[string][]byte
//	capl      capl
//}
//
//func (t tinyLfuCache) Get(key string) ([]byte, bool) {
//	value, ok := t.container[key]
//	if ok {
//		t.capl.rotateWithH(t.capl.getPosWithH(key))
//	}
//	return value, ok
//}
//
//func (t tinyLfuCache) Set(key string, value []byte) {
//	//TODO implement me
//	panic("implement me")
//}
//
//type windowTinyLfuCache struct {
//	container map[string][]byte
//	capl      capl
//}
//
//func (w windowTinyLfuCache) Get(key string) ([]byte, bool) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (w windowTinyLfuCache) Set(key string, value []byte) {
//	//TODO implement me
//	panic("implement me")
//}

func NewCacheService[T comparable](cacheType string) CacheService[T] {
	switch cacheType {
	case LRU:
		return &lruCache[T]{
			container: make(map[T][]byte),
			capl:      make([]*capacity[T], 5),
		}
	case LFU:
		return &lfuCache[T]{
			container: make(map[T][]byte),
			capl:      make([]*capacity[T], 5),
		}
		//case TinyLfu:
		//	return &tinyLfuCache{
		//		container: make(map[string][]byte),
		//		capl:      make([]*capacity, 5),
		//	}
		//case WindowTinyLfu:
		//	return &windowTinyLfuCache{
		//		container: make(map[string][]byte),
		//		capl:      make([]*capacity, 5),
		//	}
	}
	return nil
}
