package roaring

import (
	"math"
	"math/rand"
	"sort"
	"testing"
)

func TestRoaring(t *testing.T) {
	rb := New()
	for i := 0; i < 4096; i++ {
		rb.Add(uint32(i))
	}
	rb.Add(math.MaxUint16)
	if !rb.Contains(1) {
		t.Error("Contains")
	}
}

var bitmapSize = 1000000

func getBuffer(size, seed int) []uint32 {
	rand.Seed(int64(seed))

	set := make(map[uint32]struct{})
	buffer := make([]uint32, 0, size)

	for len(set) < size {
		set[rand.Uint32()] = struct{}{}
	}
	for v, _ := range set {
		buffer = append(buffer, v)
	}
	return buffer
}

func BenchmarkAdd(b *testing.B) {
	rb := New()
	buffer := getBuffer(bitmapSize, 42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Add(buffer[i%bitmapSize])
	}
}

func BenchmarkAddSorted(b *testing.B) {
	rb := New()
	buffer := getBuffer(bitmapSize, 42)
	sort.Slice(buffer, func(i, j int) bool { return buffer[i] < buffer[j] })

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Add(buffer[i%bitmapSize])
	}
}

func BenchmarkContains(b *testing.B) {
	rb := New()
	buffer := getBuffer(bitmapSize, 42)
	sort.Slice(buffer, func(i, j int) bool { return buffer[i] < buffer[j] })

	for _, v := range buffer {
		rb.Add(v)
	}

	testBuffer := getBuffer(bitmapSize, 24)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Contains(testBuffer[i%bitmapSize])
	}
}
