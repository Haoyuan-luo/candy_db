package common

import (
	"bytes"
	"sort"
)

type node struct {
	key   []byte
	score float64
	data  int64
}

func calcScore(key []byte) float64 {
	var hash uint64
	l := len(key)
	if l > 8 {
		l = 8
	}
	for i := 0; i < l; i++ {
		shift := uint(64 - 8*(i+1))
		hash |= uint64(key[i]) << shift
	}
	return float64(hash)
}

func NewNode(data int64, key []byte) *node {
	return &node{
		key:   key,
		score: calcScore(key),
		data:  data,
	}
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
