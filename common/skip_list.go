package common

import "math/rand"

type SkipList struct {
	header   *Element
	maxLevel int
}

func NewSkipList(oriNode NodeList) *SkipList {
	oriNode.SortBy(func(n1, n2 *node) bool {
		return n1.CompareWith(n2, Less)
	})
	s := &SkipList{header: NewElement(oriNode[0], len(oriNode)), maxLevel: len(oriNode)}

	for _, orin := range oriNode[1:] {
		level := s.randLevel()
		s.header.AddElement(NewElement(orin, level), level)
	}
	return s
}

func (list *SkipList) randLevel() int {
	i := 1
	for ; i < list.maxLevel; i++ {
		if rand.Int31n(2) == 0 {
			return i
		}
	}
	return i
}
