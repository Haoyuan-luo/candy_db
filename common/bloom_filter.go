package common

import (
	"candy_db/common/util"
	"math"
)

type BloomFilter struct {
	tolerance float64
	account   int
}

type BloomService interface {
	GetBloomArray(key ...[]byte) []int
}

func NewBloomFilter() BloomService {
	return &BloomFilter{tolerance: Tolerance, account: 0}
}

func (bf *BloomFilter) preKeyRatio() float64 {
	size := -1 * float64(bf.account) * math.Log(bf.tolerance) / math.Pow(math.Log(2), 2)
	return math.Ceil(size / float64(bf.account))
}

func (bf *BloomFilter) calcHshNum(ratio float64) uint32 {
	num := uint32(ratio * math.Log(2))
	if num < 1 {
		num = 1
	}
	if num > 32 {
		num = 32
	}
	return num
}

func (bf *BloomFilter) GetBloomArray(key ...[]byte) []int {
	done := make(chan struct{})
	defer close(done)
	preKeyRatio := bf.preKeyRatio()
	hashNum := bf.calcHshNum(preKeyRatio)
	nBits := len(key) * int(preKeyRatio)
	tFnv := util.Hash().Times(hashNum)
	setPoint := func(point []uint64) []int {
		filter := make([]int, nBits)
		for i := range point {
			filter[int(point[i])%nBits] = 1
		}
		return filter
	}
	ret := make([]int, len(key))
	for r := range util.Map(util.Map(key, tFnv.CalcHash).ToSlice(), setPoint) {
		for i, _ := range ret {
			ret[i] = ret[i] | r[i]
		}
	}
	return ret
}
