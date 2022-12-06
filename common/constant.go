package common

import "unsafe"

// skip list 常量控制
const (
	MaxLevel = 32 // 跳表最大层数
)

// arena 常量控制
const (
	OneElementSize = int(unsafe.Sizeof(Element{})) // 单个元素的大小
	OffSetSize     = int(unsafe.Sizeof(uint32(0))) // 适配不同机器的偏移量
)

// bloom filter 常量控制
const Tolerance float64 = 0.8 // 布隆过滤器的容错率
