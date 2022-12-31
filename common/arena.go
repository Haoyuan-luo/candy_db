package common

import (
	"candy_db/common/util"
	"sync/atomic"
)

type arena struct {
	buf []byte
	pos uint32
	log util.LogImpl
}

func NewArena() *arena {
	return &arena{
		buf: make([]byte, 1),
		pos: 1,
		log: util.Logger().SetField("Arena"),
	}
}

func (a *arena) allocate(sz uint32) uint32 {
	newPos := atomic.AddUint32(util.GenPtr(a.pos), sz)
	if len(a.buf)-int(newPos) < OneNodeSize {
		growBy := uint32(len(a.buf))
		if growBy > 1<<30 {
			growBy = 1 << 30
		}
		if growBy < sz {
			growBy = sz
		}
		newBuf := make([]byte, len(a.buf)+int(growBy))
		if !util.GenEqual(len(a.buf), copy(newBuf, a.buf)) {
			a.log.Log(util.ERROR, "arena allocate failed")
		}
		a.buf = newBuf
	}
	return newPos - sz
}

func (a *arena) CopyBy(fn func() []byte) (startPoint, size uint32) {
	startPoint = a.pos
	src := fn()
	size = uint32(len(src))
	offset := a.allocate(size)
	if !util.GenEqual(len(src), copy(a.buf[offset:], src)) {
		a.log.Log(util.ERROR, "arena copy failed")
	}
	offset += uint32(len(src))
	a.pos = offset
	return
}

func (a *arena) GetBy(offset, size uint32, fn func([]byte)) {
	fn(a.buf[offset : offset+size])
}
