package goroar

import "testing"

func TestBinarySearch(t *testing.T) {
	a := []uint16{0, 1, 2, 3, 4}
	pos := binarySearch(a, 5, 2)
	if pos != 2 {
		t.Errorf("Position: %d, want: 2", pos)
	}
}

func TestLocalIntersect2by2(t *testing.T) {
	set1 := []uint16{1, 2, 3, 4, 5}
	set2 := []uint16{3, 4, 5, 6, 7}

	length, newSet := localIntersect2by2(set1, 5, set2, 5)
	if length != 3 {
		t.Errorf("Length: %d, want: %d", length, 3)
	}

	if newSet[0] != 3 && newSet[1] != 4 && newSet[2] != 5 {
		t.Errorf("Set content: %v, want: [3, 4, 5]", newSet)
	}
}

func TestLocalIntersect2by2_2(t *testing.T) {
	set1 := []uint16{1, 2, 3, 4, 5}
	set2 := []uint16{5}

	length, newSet := localIntersect2by2(set1, 5, set2, 1)
	if length != 1 {
		t.Errorf("Length: %d, want: %d", length, 1)
	}

	if newSet[0] != 5 {
		t.Errorf("Set content: %v, want: [5]", newSet)
	}
}

func TestIntersect2by2(t *testing.T) {
	set1 := make([]uint16, 10)
	set2 := make([]uint16, 1024)

	for i := range set1 {
		set1[i] = uint16(i * 100)
	}

	for i := range set2 {
		set2[i] = uint16(i)
	}

	length, newSet := intersect2by2(set1, len(set1), set2, len(set2))
	if length != 10 {
		t.Errorf("Length: %d, want: %d", length, 10)
	}
	for i, v := range newSet {
		if uint16(i*100) != v {
			t.Errorf("Got: %d, want: %d", v, i)
		}
	}
}

func TestDifference(t *testing.T) {
	set1 := []uint16{1, 2, 3, 4, 5}
	set2 := []uint16{3, 4, 5, 6, 7}

	length, newSet := difference(set1, 5, set2, 5)
	if length != 2 {
		t.Errorf("Length: %d, want: %d", length, 3)
	}

	if newSet[0] != 1 && newSet[1] != 2 {
		t.Errorf("Set content: %v, want: [3, 4, 5]", newSet)
	}
}

func TestCountBits(t *testing.T) {
	x := uint64(0)
	if bits := countBits(x); bits != 0 {
		t.Errorf("Count bits: %d, want 0", bits)
	}

	x = uint64(1)
	if bits := countBits(x); bits != 1 {
		t.Errorf("Count bits: %d, want 1", bits)
	}

	x = ^uint64(0x5555555555555555)
	if bits := countBits(x); bits != 32 {
		t.Errorf("Count bits: %d, want 32", bits)
	}
}

func TestFillArrayAND(t *testing.T) {
	bitmap1 := []uint64{50}
	bitmap2 := []uint64{50}

	answer := fillArrayAND(bitmap1, bitmap2, 3)
	if answer[0] != uint16(1) && answer[1] != uint16(4) &&
		answer[2] != uint16(5) {
		t.Errorf("fillArrayAND error")
	}
}

func TestTrailingZeros(t *testing.T) {
	for i := 0; i < 63; i++ {
		v := trailingZeros(1 << uint(i))
		if v != i {
			t.Errorf("TrailingZeros: %d, want: %d", v, i)
		}
	}
}

func TestUnion2by2_1(t *testing.T) {
	count := 100
	set1 := make([]uint16, 0, count)
	set2 := make([]uint16, 0, count)

	for i := 0; i < count; i += 2 {
		set1 = append(set1, uint16(i))
		set2 = append(set2, uint16(i+1))
	}

	total, buffer := union2by2(set1, len(set1), set2, len(set2), len(set1)+len(set2))
	if total != count {
		t.Errorf("Union total: %d, want: %d", total, count)
	}
	for k, v := range buffer {
		if uint16(k) != v {
			t.Errorf("Union: %d, want: %d", v, k)
			break
		}
	}
}

func TestUnion2by2_2(t *testing.T) {
	count := 100
	set1 := make([]uint16, 0, count)
	set2 := make([]uint16, 0, count)

	for i := 0; i < count; i++ {
		set1 = append(set1, uint16(i))
	}

	total, buffer := union2by2(set1, len(set1), set2, len(set2), len(set1)+len(set2))
	if total != count {
		t.Errorf("Union total: %d, want: %d", total, count)
	}

	for k, v := range buffer {
		if set1[k] != v {
			t.Errorf("Union: %d, want: %d", v, k)
			break
		}
	}
}

func TestUnion2by2_3(t *testing.T) {
	count := 100
	set1 := make([]uint16, 0, count)
	set2 := make([]uint16, 0, count)

	for i := 0; i < count; i++ {
		set2 = append(set2, uint16(i))
	}

	total, buffer := union2by2(set1, len(set1), set2, len(set2), len(set1)+len(set2))
	if total != count {
		t.Errorf("Union total: %d, want: %d", total, count)
	}

	for k, v := range buffer {
		if set2[k] != v {
			t.Errorf("Union: %d, want: %d", v, k)
			break
		}
	}
}

func TestUnion2by2_4(t *testing.T) {
	count := 10
	set1 := make([]uint16, 0, count)
	set2 := make([]uint16, 0, count)

	for i := 0; i < count; i++ {
		set1 = append(set1, uint16(i))
		set2 = append(set2, uint16(i))
	}

	total, buffer := union2by2(set1, len(set1), set2, len(set2), len(set1)+len(set2))
	if total != count {
		t.Errorf("Union total: %d, want: %d", total, count)
	}

	for k, v := range buffer {
		if set2[k] != v {
			t.Errorf("Union: %d, want: %d", v, k)
			break
		}
	}
}

func TestUnion2by2_5(t *testing.T) {
	set1 := make([]uint16, 0)
	set2 := make([]uint16, 0)

	total, buffer := union2by2(set1, len(set1), set2, len(set2), len(set1)+len(set2))
	if total != 0 {
		t.Errorf("Union total: %d, want: %d", total, 0)
	}

	for k, v := range buffer {
		if set2[k] != v {
			t.Errorf("Union: %d, want: %d", v, k)
			break
		}
	}
}
