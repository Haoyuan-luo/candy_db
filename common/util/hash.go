package util

const offset32 fnv = 2166136261
const prime32 fnv = 16777619

type fnv uint64

// 一个计算多次的Fnv工具
type timesFnv struct {
	fnv
	times uint32
}

type HashImpl interface {
	Times(n uint32) timesFnv
	CalcHash(key []byte) uint64
}

func Hash() HashImpl {
	s := offset32
	return &s
}

func (s *fnv) CalcHash(key []byte) uint64 {
	sampled := make([][]byte, 0, 3)
	total := len(key) - 1
	step := total / 5
	sampled = append(sampled, key[0:step], key[step*2:step*3], key[total-step:total])
	hash := *s
	for _, v := range sampled {
		for _, b := range v {
			hash *= prime32
			hash ^= fnv(b)
		}
	}
	return uint64(hash)
}

func (s *fnv) Times(n uint32) timesFnv {
	return timesFnv{*s, n}
}

func (s *timesFnv) CalcHash(key []byte) []uint64 {
	hash := make([]uint64, s.times)
	for i := uint32(0); i < s.times; i++ {
		hash[i] = s.fnv.CalcHash(append(key, byte(i)))
	}
	return hash
}
