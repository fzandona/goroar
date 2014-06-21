package goroar

import (
	"math/rand"
	"sort"
	"testing"
)

var bitmapSize = 1000000

type Values []uint32

func (v Values) Len() int           { return len(v) }
func (v Values) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Values) Less(i, j int) bool { return v[i] < v[j] }

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
	var pos int
	buffer := getBuffer(bitmapSize, 42)

	b.N = bitmapSize
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Add(buffer[pos])
		pos++
	}
}

func BenchmarkAddSorted(b *testing.B) {
	rb := New()
	var pos int
	buffer := getBuffer(bitmapSize, 42)
	sort.Sort(Values(buffer))

	b.N = bitmapSize
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Add(buffer[pos])
		pos++
	}
}

func BenchmarkIterator(b *testing.B) {
	rb := New()
	var pos int
	buffer := getBuffer(bitmapSize, 42)
	sort.Sort(Values(buffer))

	for _, v := range buffer {
		rb.Add(v)
	}

	b.ResetTimer()
	for v := range rb.Iterator() {
		if v > 0 {
		}
		pos++
	}
}

func BenchmarkContains(b *testing.B) {
	rb := New()
	var pos int
	buffer := getBuffer(bitmapSize, 42)
	sort.Sort(Values(buffer))

	for _, v := range buffer {
		rb.Add(v)
	}

	testBuffer := getBuffer(bitmapSize, 24)
	b.N = bitmapSize
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Contains(testBuffer[pos])
		pos++
	}
}

func getTwoRB(size1, seed1, size2, seed2 int) (*RoaringBitmap, *RoaringBitmap) {
	rb1 := New()
	rb2 := New()

	buffer := getBuffer(size1, seed1)
	sort.Sort(Values(buffer))
	for _, v := range buffer {
		rb1.Add(v)
	}

	buffer = getBuffer(size2, seed2)
	sort.Sort(Values(buffer))
	for _, v := range buffer {
		rb2.Add(v)
	}

	return rb1, rb2
}

func BenchmarkAndSameContent(b *testing.B) {
	rb1, rb2 := getTwoRB(bitmapSize, 42, bitmapSize, 42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.And(rb2)
	}
}

func BenchmarkAndDifferentContent(b *testing.B) {
	rb1, rb2 := getTwoRB(bitmapSize, 42, bitmapSize, 24)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.And(rb2)
	}
}

func BenchmarkOrSameContent(b *testing.B) {
	rb1, rb2 := getTwoRB(bitmapSize, 42, bitmapSize, 42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.Or(rb2)
	}
}

func BenchmarkOrDifferentContent(b *testing.B) {
	rb1, rb2 := getTwoRB(bitmapSize, 42, bitmapSize, 24)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.Or(rb2)
	}
}

func BenchmarkXorSameContent(b *testing.B) {
	rb1, rb2 := getTwoRB(bitmapSize, 42, bitmapSize, 42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.Xor(rb2)
	}
}

func BenchmarkXorDifferentContent(b *testing.B) {
	rb1, rb2 := getTwoRB(bitmapSize, 42, bitmapSize, 24)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.Xor(rb2)
	}
}

func BenchmarkAndNotSameContent(b *testing.B) {
	rb1, rb2 := getTwoRB(bitmapSize, 42, bitmapSize, 42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.AndNot(rb2)
	}
}

func BenchmarkAndNotDifferentContent(b *testing.B) {
	rb1, rb2 := getTwoRB(bitmapSize, 42, bitmapSize, 24)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.AndNot(rb2)
	}
}

func getTwoSimpleRB(size1, size2 int) (*RoaringBitmap, *RoaringBitmap) {
	rb1 := New()
	rb2 := New()

	for i := 0; i < size1; i++ {
		rb1.Add(uint32(i))
	}

	for i := 0; i < size2; i++ {
		rb2.Add(uint32(i))
	}

	return rb1, rb2
}

func BenchmarkAndSimpleContent(b *testing.B) {
	rb1, rb2 := getTwoSimpleRB(bitmapSize, bitmapSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.And(rb2)
	}
}

func BenchmarkOrSimpleContent(b *testing.B) {
	rb1, rb2 := getTwoSimpleRB(bitmapSize, bitmapSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.Or(rb2)
	}
}

func BenchmarkXorSimpleContent(b *testing.B) {
	rb1, rb2 := getTwoSimpleRB(bitmapSize, bitmapSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.Xor(rb2)
	}
}

func BenchmarkAndNotSimpleContent(b *testing.B) {
	rb1, rb2 := getTwoSimpleRB(bitmapSize, bitmapSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb1.AndNot(rb2)
	}
}
