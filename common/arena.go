package common

import (
	"sync/atomic"
)

type arena struct {
	buf []byte
	pos uint32
	log LogImpl
}

type ArenaImpl interface {
	createSpace(level int) uint32
	PutData(node node)
	GetElement(pos, offset int) *Element
}

func NewArena() ArenaImpl {
	return &arena{
		buf: make([]byte, 1),
		pos: 1,
		log: Logger().SetField("Arena"),
	}
}

func (a arena) allocate(sz uint32) uint32 {
	newPos := atomic.AddUint32(GenPtr(a.pos), sz)
	if len(a.buf)-int(newPos) < OneElementSize {
		growBy := uint32(len(a.buf))
		if growBy > 1<<30 {
			growBy = 1 << 30
		}
		if growBy < sz {
			growBy = sz
		}
		newBuf := make([]byte, len(a.buf)+int(growBy))
		if !GenEqual(len(a.buf), copy(newBuf, a.buf)) {
			a.log.Log(ERROR, "arena allocate failed")
		}
		a.buf = newBuf
	}
	return newPos - sz
}

func (a arena) createSpace(level int) uint32 {
	unusedSize := (MaxLevel - level) * OffSetSize
	return a.allocate(uint32(OneElementSize - unusedSize))
}

func (a arena) PutData(node node) {
	//TODO implement me
	panic("implement me")
}

func (a arena) GetElement(pos, offset int) *Element {
	//TODO implement me
	panic("implement me")
}
