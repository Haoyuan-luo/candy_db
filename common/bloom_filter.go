package common

import (
	"candy_db/common/util"
	"math"
)

type BloomFilter struct {
	tolerance float64
	account   int
}

func NewBloomFilter() *BloomFilter {
	return &BloomFilter{
		tolerance: Tolerance,
		account:   0,
	}
}

func (bf *BloomFilter) preKeyRatio() float64 {
	size := -1 * float64(bf.account) * math.Log(bf.tolerance) / math.Pow(math.Log(2), 2)
	return float64(math.Ceil(size / float64(bf.account)))
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
	tFnv := util.Hash().times(hashNum)
	filter := make([]int, nBits)
	producer := util.Producer[[]byte](done, key)
	processor := util.Processor[[]byte, []float64](done, producer, tFnv.simpleFnv, nil)
	consumer := util.Consumer[[]float64](done, processor)
	setPoint := func(point []float64) struct{} {
		for i := range point {
			filter[int(point[i])%nBits] = 1
		}
		return struct{}{}
	}
	p := <-consumer
	producerV2 := util.Producer[[]float64](done, [][]float64{p})
	processorV2 := util.Processor[[]float64, struct{}](done, producerV2, setPoint, nil)
	_ = util.Consumer[struct{}](done, processorV2)
	return filter
}
