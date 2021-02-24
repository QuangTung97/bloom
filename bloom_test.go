package bloom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterSetBit(t *testing.T) {
	b := NewFilter(192, 3)
	assert.Equal(t, 3, len(b.bitset))

	b.setBit(10)
	expected := []uint64{0x400, 0, 0}
	assert.Equal(t, expected, b.bitset)

	b.setBit(129)
	expected = []uint64{0x400, 0, 0x2}
	assert.Equal(t, expected, b.bitset)
}

func TestFilterBitIsSet(t *testing.T) {
	b := NewFilter(192, 3)
	b.setBit(10)
	b.setBit(20)

	assert.True(t, b.bitIsSet(10))
	assert.True(t, b.bitIsSet(20))
	assert.False(t, b.bitIsSet(19))
	assert.False(t, b.bitIsSet(21))
}

func TestFilter(t *testing.T) {
	b := NewFilter(65, 3)

	assert.Equal(t, b.m, uint64(65))
	assert.Equal(t, b.k, uint64(3))
	assert.Equal(t, 2, len(b.bitset))

	b.Put([]byte("element-1"))
	b.PutString("element-2")
	b.Put([]byte("element-3"))

	assert.True(t, b.Test([]byte("element-1")))
	assert.True(t, b.Test([]byte("element-2")))
	assert.True(t, b.TestString("element-3"))
	assert.False(t, b.Test([]byte("element-4")))
	assert.False(t, b.TestString("element-4"))
}

func TestEstimateParameters(t *testing.T) {
	m, k := EstimateParameters(1000000, 0.001)
	assert.Equal(t, uint64(14377588), m)
	assert.Equal(t, uint64(10), k)
}