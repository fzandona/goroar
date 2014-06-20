package test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/fzandona/goroar"
)

var count = 2000000
var stop = 1000000

type Values []uint32

func (v Values) Len() int           { return len(v) }
func (v Values) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Values) Less(i, j int) bool { return v[i] < v[j] }

var roarBitmap *goroar.RoaringBitmap = goroar.New()
var pos int

func getBuffer() []uint32 {
	rand.Seed(42)

	set := make(map[uint32]struct{})
	buffer := make([]uint32, 0, stop)

	for len(set) < stop {
		set[rand.Uint32()] = struct{}{}
	}
	for v, _ := range set {
		buffer = append(buffer, v)
	}
	return buffer
}
func BenchmarkAdd(b *testing.B) {
	buffer := getBuffer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		roarBitmap.Add(buffer[pos%stop])
		pos++
	}
}

var roarBitmapSorted *goroar.RoaringBitmap = goroar.New()
var posSorted int

func BenchmarkAddSorted(b *testing.B) {
	buffer := getBuffer()
	sort.Sort(Values(buffer))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		roarBitmapSorted.Add(buffer[posSorted%stop])
		posSorted++
	}
}
