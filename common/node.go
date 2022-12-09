package common

import (
	"bytes"
	"candy_db/common/util"
	"sort"
)

type nodeArena struct {
	keyOffset   uint32
	keySize     uint32
	valueOffset uint32
	valueSize   uint32
}

func (n *nodeArena) setKey(offset uint32, size uint32) {
	n.keyOffset = offset
	n.keySize = size
}

func (n *nodeArena) setVal(offset uint32, size uint32) {
	n.valueOffset = offset
	n.valueSize = size
}

type node struct {
	key   []byte
	score float64
	data  interface{}
	nodeArena
}

func NewNode(key []byte) *node {
	return &node{
		key:   key,
		score: util.Hash().simpleFnv(key),
	}
}

func (n *node) ToUse(a *arena) {

}

func (n *node) SetData(value interface{}) *node {
	n.data = value
	return n
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
