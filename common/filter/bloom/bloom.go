package bloom

import (
	"hash/fnv"
	"math"
	"math/rand"
	"time"
)

// uint64的位数
const uint64Bits = 64

// Filter 布隆过滤器
// https://llimllib.github.io/bloomfilter-tutorial/
// https://github.com/bits-and-blooms/bloom/blob/master/bloom.go
type Filter struct {
	bits   []uint64 // bit数组
	bitCnt uint64   // bit位数
	seeds  []uint64 // 哈希种子
}

func New(capacity uint64, falsePositiveRate float64) *Filter {
	// bit数量
	factor := -math.Log(falsePositiveRate) / (math.Ln2 * math.Ln2)
	bitCnt := uint64(math.Ceil(float64(capacity) * factor))
	// 这里扩大到最后一个uint64大小，避免浪费
	bitCnt = (bitCnt + uint64Bits - 1) / uint64Bits * uint64Bits
	// 哈希函数数量
	seedCnt := int(math.Ceil(math.Ln2 * float64(bitCnt) / float64(capacity)))
	seeds := make([]uint64, seedCnt)
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < seedCnt; i++ {
		seeds[i] = source.Uint64()
	}
	return &Filter{
		bits:   make([]uint64, bitCnt/uint64Bits),
		bitCnt: bitCnt,
		seeds:  seeds,
	}
}

// Add 添加元素
func (f *Filter) Add(hash uint64) {
	for _, seed := range f.seeds {
		index, offset := f.pos(hash, seed)
		f.bits[index] |= 1 << offset
	}
}

// AddBytes 添加元素
func (f *Filter) AddBytes(b []byte) {
	f.Add(f.hash(b))
}

// AddString 添加元素
// 字符串类型
func (f *Filter) AddString(s string) {
	f.AddBytes([]byte(s))
}

// Contains 元素是否存在
// true表示可能存在
func (f *Filter) Contains(hash uint64) bool {
	for _, seed := range f.seeds {
		index, offset := f.pos(hash, seed)
		mask := uint64(1) << offset
		// 判断这一位是否位1
		if (f.bits[index] & mask) != mask {
			return false
		}
	}
	return true
}

// ContainsBytes 元素是否存在
// true表示可能存在
func (f *Filter) ContainsBytes(b []byte) bool {
	return f.Contains(f.hash(b))
}

// ContainsString 元素是否存在
// 字符串类型
func (f *Filter) ContainsString(s string) bool {
	return f.ContainsBytes([]byte(s))
}

// Clear 清空过滤器
func (f *Filter) Clear() {
	for i := range f.bits {
		f.bits[i] = 0
	}
}

// Len bit位数
func (f *Filter) Len() uint64 {
	return f.bitCnt
}

// 获取对应元素下标和偏移
func (f *Filter) pos(h, seed uint64) (uint64, uint64) {
	// 按照位计算的偏移
	bitsIndex := (h ^ seed) % f.bitCnt
	// 因为一个元素64位，因此需要转换
	index := bitsIndex / uint64Bits
	// 在一个元素里面的偏移
	offset := bitsIndex % uint64Bits
	return index, offset
}

// 计算哈希值
func (f *Filter) hash(b []byte) uint64 {
	fnvHash := fnv.New64()
	fnvHash.Write(b)
	return fnvHash.Sum64()
}
