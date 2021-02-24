package bloom

import "math"

// Filter ...
type Filter struct {
	// number of bits
	m uint64
	// number of hash functions
	k uint64
	// bitset data
	bitset []uint64
}

const (
	hashKey0 uint64 = 0xab123456
	hashKey1 uint64 = 0xccddaa32
	hashKey2 uint64 = 0x34ddacdd
	hashKey3 uint64 = 0x1255abdd
)

// NewFilter creates a Bloom Filter
func NewFilter(m uint64, k uint64) *Filter {
	return &Filter{
		m: m,
		k: k,
		// (m + 64) >> 6 = ceil(m / 64)
		bitset: make([]uint64, (m+63)>>6),
	}
}

func (b *Filter) setBit(i uint64) {
	pos := i & 0x3f
	mask := uint64(1 << pos)
	b.bitset[i>>6] |= mask
}

func (b *Filter) bitIsSet(i uint64) bool {
	pos := i & 0x3f
	mask := uint64(1 << pos)
	return b.bitset[i >> 6] & mask != 0
}

// location returns the ith hashed location using the four base hash values
func location(h [4]uint64, i uint64) uint64 {
	return h[i%2] + i*h[2+(((i+(i%2))%4)/2)]
}

func (b *Filter) location(h [4]uint64, i uint64) uint64 {
	return location(h, i) % b.m
}

func computeHashArray(data []byte) [4]uint64 {
	return [4]uint64{
		memhashSliceKey(data, hashKey0),
		memhashSliceKey(data, hashKey1),
		memhashSliceKey(data, hashKey2),
		memhashSliceKey(data, hashKey3),
	}
}

func computeHashArrayString(s string) [4]uint64 {
	return [4]uint64{
		memhashStringKey(s, hashKey0),
		memhashStringKey(s, hashKey1),
		memhashStringKey(s, hashKey2),
		memhashStringKey(s, hashKey3),
	}
}

// EstimateParameters estimates requirements for m and k.
// Based on https://bitbucket.org/ww/bloom/src/829aa19d01d9/bloom.go
// used with permission.
func EstimateParameters(n uint64, p float64) (m uint64, k uint64) {
	m = uint64(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k = uint64(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return
}

// Put element into bloom filter
func (b *Filter) Put(data []byte) {
	hash := computeHashArray(data)

	for i := uint64(0); i < b.k; i++ {
		b.setBit(b.location(hash, i))
	}
}

// PutString put string element into bloom filter
func (b *Filter) PutString(s string) {
	hash := computeHashArrayString(s)

	for i := uint64(0); i < b.k; i++ {
		b.setBit(b.location(hash, i))
	}
}

// Test checks if *data* exists in bloom filter
func (b *Filter) Test(data []byte) bool {
	hash := computeHashArray(data)

	for i := uint64(0); i < b.k; i++ {
		if !b.bitIsSet(b.location(hash, i)) {
			return false
		}
	}
	return true
}

// TestString checks if *data* exists in bloom filter
func (b *Filter) TestString(s string) bool {
	hash := computeHashArrayString(s)

	for i := uint64(0); i < b.k; i++ {
		if !b.bitIsSet(b.location(hash, i)) {
			return false
		}
	}
	return true
}
