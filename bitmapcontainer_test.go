package goroar

import (
	"math/rand"
	"testing"
)

func TestContains(t *testing.T) {
	bc := newBitmapContainer()

	for i := 0; i < 10; i++ {
		bc.add(uint16(i))
	}

	if bc.contains(uint16(100)) {
		t.Errorf("Contains: %v, want: %v", true, false)
	}

	for i := 0; i < 10; i++ {
		contains := bc.contains(uint16(i))
		if !contains {
			t.Errorf("Contains: %v, want: %v", contains, true)
		}
	}
}

func TestContains_2(t *testing.T) {
	bc := newBitmapContainer()

	rand.Seed(42)
	for i := 0; i < 1000; i++ {
		bc.add(uint16(rand.Int31n(1 << 16)))
	}

	rand.Seed(42)
	for i := 0; i < 1000; i++ {
		contains := bc.contains(uint16(rand.Int31n(1 << 16)))
		if !contains {
			t.Errorf("Contains: %v, want: %v", contains, true)
		}
	}
}

func TestContains_3(t *testing.T) {
	bc := newBitmapContainer()

	for i := 0; i < 10000; i += 2 {
		bc.add(uint16(i))
	}

	for i := 0; i < 10000; i += 2 {
		if !bc.contains(uint16(i)) && bc.contains(uint16(i+1)) {
			t.Error("Wrong content")
		}
	}
}

func TestContains_4(t *testing.T) {
	bc := newBitmapContainer()

	for i := 0; i < (1 << 16); i += 2 {
		bc.add(uint16(i))
	}

	for i := 0; i < (1 << 16); i += 2 {
		if !bc.contains(uint16(i)) && bc.contains(uint16(i+1)) {
			t.Error("Contains fails at: %d), i")
		}
	}
}
func TestAndBitmapBitmap(t *testing.T) {
	bt1 := newBitmapContainer()
	bt2 := newBitmapContainer()

	for i := 0; i < 4000; i++ {
		bt1.add(uint16(i))
		bt2.add(uint16(i))
	}

	answer := bt1.andBitmap(bt2)
	switch ac := answer.(type) {
	case *arrayContainer:
		if ac.cardinality != 4000 {
			t.Errorf("Cardinality: %d, want: 4000", ac.cardinality)
		}
		for i := 0; i < 4000; i++ {
			if ac.content[i] != uint16(i) {
				t.Errorf("AndBitmap: %d, want: %d", ac.content[i], i)
				break
			}
		}
	default:
		t.Errorf("Wrong container type: %T", ac)
	}
}

func TestAndBitmapBitmapEmpty(t *testing.T) {
	bt1 := newBitmapContainer()
	bt2 := newBitmapContainer()

	for i := 0; i < 4000; i++ {
		bt1.add(uint16(i))
		bt2.add(uint16(i + 4000))
	}

	answer := bt1.andBitmap(bt2)
	switch ac := answer.(type) {
	case *arrayContainer:
		if ac.cardinality != 0 {
			t.Errorf("Cardinality: %d, want: 0", ac.cardinality)
		}
	default:
		t.Errorf("Wrong container type: %T", ac)
	}
}

func TestAndBitmapBitmap_2(t *testing.T) {
	bt1 := newBitmapContainer()
	bt2 := newBitmapContainer()

	for i := 0; i < (1 << 16); i++ {
		bt1.add(uint16(i))
		bt2.add(uint16(i))
	}

	answer := bt1.andBitmap(bt2)
	switch c := answer.(type) {
	case *bitmapContainer:
		if c.cardinality != (1 << 16) {
			t.Errorf("Cardinality: %d, want: %d", c.cardinality, 1<<16)
		}
		for i := 0; i < 1<<16; i++ {
			if !c.contains(uint16(i)) {
				t.Errorf("BitmapContainer does not contain: %d", i)
				break
			}
		}
	default:
		t.Errorf("Wrong container type: %T", c)
	}
}

func TestAndBitmapBitmap_3(t *testing.T) {
	bt1 := newBitmapContainer()
	bt2 := newBitmapContainer()

	for i := 0; i < (1 << 16); i++ {
		bt1.add(uint16(i))
	}

	answer := bt1.andBitmap(bt2)
	switch c := answer.(type) {
	case *arrayContainer:
		if c.cardinality != 0 {
			t.Errorf("Cardinality: %d, want: %d", c.cardinality, 0)
		}
	default:
		t.Errorf("Wrong container type: %T", c)
	}
}

func TestAndBitmapArray(t *testing.T) {
	bt := newBitmapContainer()
	ac := newArrayContainer()

	for i := 0; i < (1 << 16); i++ {
		bt.add(uint16(i))
	}
	for i := 0; i < 4096; i++ {
		ac.add(uint16(i))
	}

	answer := bt.andArray(ac)
	if answer.cardinality != 4096 {
		t.Errorf("Cardinality: %d, want: 4096", answer.cardinality)
	}
	for i := 0; i < 4096; i++ {
		if answer.content[i] != uint16(i) {
			t.Errorf("AndBitmapArray: %d, want: %d", answer.content[i], i)
			break
		}
	}
}

func TestNextBitSet(t *testing.T) {
	bc := newBitmapContainer()

	if bc.nextSetBit(0) != -1 {
		t.Errorf("NextBitSet: %d, want: -1", bc.nextSetBit(0))
	}

	result := [6]int{1, 2, 4, 8, 16, 32}

	for i := 0; i < 6; i++ {
		bc = newBitmapContainer()
		bc.add(1 << uint(i))
		v := bc.nextSetBit(i)
		if v != result[i] {
			t.Errorf("NextBitSet: %d, want: %d, %b", v, result[i], bc.bitmap[0])
		}
	}
}

func TestToArrayContainer(t *testing.T) {
	bc := newBitmapContainer()

	for i := 0; i < 4096; i++ {
		bc.add(uint16(i))
	}

	ac := bc.toArrayContainer()

	if ac.cardinality != 4096 {
		t.Errorf("Cardinality: %d, want: %d", ac.cardinality, 4096)
	}
	for i := 0; i < 4096; i++ {
		if ac.content[i] != uint16(i) {
			t.Errorf("Content: %d, want: %d", ac.content[i], i)
			break
		}
	}
}
