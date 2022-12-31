package common

import (
	"candy_db/common/util"
	"math/rand"
	"sync"
)

type SkipList struct {
	header *Element
	arena  *arena
	mutex  sync.RWMutex
	update []*Element
	max    int
	skip   int
	level  int
	length int32
}

type SkipListImpl interface {
	AddNode(nodes ...*Container)
	FindNode(container ...*Container)
}

func NewSkipList(skip ...int) SkipListImpl {
	list := &SkipList{header: NewElement(nil, MaxLevel)}
	list.arena = NewArena()
	list.max = MaxLevel
	list.skip = 4
	list.level = 0
	list.length = 0
	list.update = make([]*Element, list.max)
	if len(skip) == 1 && skip[0] > 1 {
		list.skip = skip[0]
	}
	return list
}

func (list *SkipList) AddNode(container ...*Container) {
	for i := range container {
		list.add(newNode(list.arena).Dump(container[i]))
	}
}

func (list *SkipList) FindNode(container ...*Container) {
	done := make(chan struct{})
	defer close(done)
	nodes := util.Map(container, list.find)

	for i := 0; i < len(container); i++ {
		if n := <-nodes; n != nil {
			n.Replay(container[i])
		}
	}
	return
}

func (list *SkipList) find(container *Container) *node {
	tarNode := newNode(list.arena).identity(container.key)
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var prev = list.header
	var next *Element
	for i := list.level - 1; i >= 0; i-- {
		next = prev.Level[i]
		for next != nil && next.Data.CompareWith(tarNode, Less) {
			prev = next
			next = prev.Level[i]
		}
	}

	if next != nil && next.Data.CompareWith(tarNode, Equal) {
		return next.Data
	} else {
		return nil
	}
}

func (list *SkipList) add(tarNode *node) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	//获取每层的前驱节点=>list.update
	var prev = list.header
	var next *Element
	for i := list.level - 1; i >= 0; i-- {
		next = prev.Level[i]
		for next != nil && next.Data.CompareWith(tarNode, Less) {
			prev = next
			next = prev.Level[i]
		}
		list.update[i] = prev
	}

	//如果key已经存在
	if next != nil && next.Data.CompareWith(tarNode, Equal) {
		next.Data = tarNode
		return
	}

	//随机生成新结点的层数
	level := list.randomLevel()
	if level > list.level {
		level = list.level + 1
		list.level = level
		list.update[list.level-1] = list.header
	}

	//申请新的结点
	ele := NewElement(tarNode, level)

	//调整next指向
	for i := 0; i < level; i++ {
		ele.Level[i] = list.update[i].Level[i]
		list.update[i].Level[i] = ele
	}
	list.length++
}

func (list *SkipList) randomLevel() int {
	i := 1
	for ; i < list.max; i++ {
		if rand.Int()%list.skip != 0 {
			break
		}
	}
	return i
}
