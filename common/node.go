package common

import (
	"bytes"
	"candy_db/common/util"
	"sort"
	"unsafe"
)

// 可自定义数据编码方式
type ctrInterface interface {
	setByteData([]byte)
	encode() []byte
	encodeSize() uint32
}

type nodeArena struct {
	impl        *arena
	valueOffset uint32
	valueSize   uint32
}

func (n *nodeArena) putVal(ctr ctrInterface) {
	n.valueOffset, n.valueSize = n.impl.CopyBy(func() []byte { return ctr.encode() })
}

func (n *nodeArena) pickVal(ctr ctrInterface) {
	n.impl.GetBy(n.valueOffset, n.valueSize, func(buf []byte) {
		ctr.setByteData(buf)
	})
}

// Container 数据出入的容器
type Container struct {
	key       []byte
	data      []byte
	expiresAT uint64
}

func NewContainer(key []byte, data []byte) *Container {
	return &Container{
		key:  key,
		data: data,
	}
}

func (d *Container) GetData() []byte {
	return d.data
}

// 数据实体的编码方式
func (d *Container) encode() (ret []byte) {
	return d.data
}

// 数据实体的编码大小
func (d *Container) encodeSize() uint32 {
	dataSize := int(unsafe.Sizeof(d.data))
	expiresSize := 0
	for {
		expiresSize++
		d.expiresAT >>= 7
		if d.expiresAT == 0 {
			break
		}
	}
	return uint32(dataSize + expiresSize)
}

// 数据实体的保存
func (d *Container) setByteData(buf []byte) {
	d.data = buf
}

type byteKey []byte

func (b byteKey) encodeSize() uint32 {
	return uint32(len(b))
}

type node struct {
	score     uint64
	key       []byte
	expiresAT uint64
	nodeArena
}

func newNode(impl *arena) *node {
	return &node{nodeArena: nodeArena{impl: impl}}
}

func (n *node) identity(key byteKey) *node {
	n.key = key
	n.score = util.Hash().CalcHash(key)
	return n
}

func (n *node) Dump(value *Container) *node {
	n.identity(value.key)
	n.expiresAT = value.expiresAT
	n.nodeArena.putVal(value)
	return n
}

func (n *node) Replay(data *Container) {
	n.nodeArena.pickVal(data)
}

type compareType int

const (
	Equal compareType = 0
	More  compareType = 1
	Less  compareType = -1
)

func (n *node) CompareWith(tar *node, want compareType) bool {
	if n.score == tar.score {
		return bytes.Compare(n.key, tar.key) == int(want)
	}
	if n.score > tar.score {
		return want == More
	}
	return want == Less
}

type sortTable struct {
	nodes []*node
	less  func(i, j *node) bool
}

func (s sortTable) Len() int {
	return len(s.nodes)
}

func (s sortTable) Less(i, j int) bool {
	return s.less(s.nodes[i], s.nodes[j])
}

func (s sortTable) Swap(i, j int) {
	tmp := s.nodes[i]
	s.nodes[i] = s.nodes[j]
	s.nodes[j] = tmp
}

type NodeList []*node

func (nl NodeList) SortBy(fn func(n1, n2 *node) bool) {
	sort.Sort(&sortTable{nl, fn})
}
