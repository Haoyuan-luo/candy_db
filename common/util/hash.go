package util

const offset32 fnv = 2166136261
const prime32 fnv = 16777619

type fnv uint32

type timesFnv struct {
	fnv
	times uint32
}

type HashImpl interface {
	times(n uint32) timesFnv
	simpleFnv(key []byte) float64
}

func Hash() HashImpl {
	s := offset32
	return &s
}

func (s *fnv) simpleFnv(key []byte) float64 {
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
	return float64(hash)
}

func (s *fnv) times(n uint32) timesFnv {
	return timesFnv{*s, n}
}

func (s *timesFnv) simpleFnv(key []byte) []float64 {
	hash := make([]float64, s.times)
	for i := uint32(0); i < s.times; i++ {
		hash[i] = s.fnv.simpleFnv(append(key, byte(i)))
	}
	return hash
}
