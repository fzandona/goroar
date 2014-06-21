package goroar

func binarySearch(array []uint16, length int, k uint16) int {
	low := 0
	high := length - 1

	for low <= high {
		middleIndex := (low + high) >> 1
		middleValue := array[middleIndex]

		switch {
		case middleValue < k:
			low = middleIndex + 1
		case middleValue > k:
			high = middleIndex - 1
		default:
			return middleIndex
		}
	}
	return -(low + 1)
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func intersect2by2(set1 []uint16, length1 int,
	set2 []uint16, length2 int) (int, []uint16) {

	if length1*64 < length2 {
		return oneSidedGallopingIntersect2by2(set1, length1, set2, length2)
	}

	if length2*64 < length1 {
		return oneSidedGallopingIntersect2by2(set2, length2, set1, length1)
	}

	return localIntersect2by2(set1, length1, set2, length2)
}

func localIntersect2by2(set1 []uint16, length1 int,
	set2 []uint16, length2 int) (int, []uint16) {

	if 0 == length1 || 0 == length2 {
		return 0, make([]uint16, 0)
	}

	finalLength := min(length1, length2)
	buffer := make([]uint16, finalLength)
	k1, k2, pos := 0, 0, 0

Mainwhile:
	for {
		if set2[k2] < set1[k1] {
			for {
				k2++
				if k2 == length2 {
					break Mainwhile
				}
				if set2[k2] >= set1[k1] {
					break
				}
			}
		}
		if set1[k1] < set2[k2] {
			for {
				k1++
				if k1 == length1 {
					break Mainwhile
				}
				if set1[k1] >= set2[k2] {
					break
				}
			}
		} else {
			buffer[pos] = set1[k1]
			pos++
			k1++
			if k1 == length1 {
				break
			}
			k2++
			if k2 == length2 {
				break
			}
		}
	}
	return pos, buffer[:pos]
}

func oneSidedGallopingIntersect2by2(
	smallSet []uint16, smallLength int,
	largeSet []uint16, largeLength int) (int, []uint16) {

	if 0 == smallLength {
		return 0, make([]uint16, 0)
	}

	buffer := make([]uint16, smallLength)
	k1, k2, pos := 0, 0, 0

	for {
		if largeSet[k1] < smallSet[k2] {
			k1 = advanceUntil(largeSet, k1, largeLength, smallSet[k2])
			if k1 == largeLength {
				break
			}
		}
		if smallSet[k2] < largeSet[k1] {
			k2++
			if k2 == smallLength {
				break
			}
		} else { // (set2[k2] == set1[k1])
			buffer[pos] = smallSet[k2]
			pos++
			k2++
			if k2 == smallLength {
				break
			}
			k1 = advanceUntil(largeSet, k1, largeLength, smallSet[k2])
			if k1 == largeLength {
				break
			}
		}

	}
	return pos, buffer[:pos]
}

// Find the smallest integer larger than pos such that array[pos]>= min.
// If none can be found, return length. Based on code by O. Kaser.
func advanceUntil(array []uint16, pos, length int, min uint16) int {
	lower := pos + 1

	// special handling for a possibly common sequential case
	if lower >= length || array[lower] >= min {
		return lower
	}

	spansize := 1 // could set larger  bootstrap an upper limit

	for (lower+spansize) < length && array[lower+spansize] < min {
		spansize *= 2
	}
	var upper int
	if lower+spansize < length {
		upper = lower + spansize
	} else {
		upper = length - 1
	}

	// maybe we are lucky (could be common case when the seek ahead
	// expected to be small and sequential will otherwise make us look bad)
	if array[upper] == min {
		return upper
	}

	if array[upper] < min { // means array has no item >= min
		return length
	}

	// we know that the next-smallest span was too small
	lower += (spansize / 2)

	// else begin binary search
	// invariant: array[lower]<min && array[upper]>min
	for lower+1 != upper {
		mid := (lower + upper) / 2
		if array[mid] == min {
			return mid
		} else if array[mid] < min {
			lower = mid
		} else {
			upper = mid
		}
	}
	return upper
}

func difference(
	set1 []uint16, length1 int,
	set2 []uint16, length2 int) (int, []uint16) {

	k1, k2, pos := 0, 0, 0

	if 0 == length2 {
		buffer := make([]uint16, length1)
		copy(buffer, set1)
		return length1, buffer
	}

	if 0 == length1 {
		return 0, make([]uint16, 0)
	}

	buffer := make([]uint16, length1)

	for {
		if set1[k1] < set2[k2] {
			buffer[pos] = set1[k1]
			pos++
			k1++
			if k1 >= length1 {
				break
			}
		} else if set1[k1] == set2[k2] {
			k1++
			k2++
			if k1 >= length1 {
				break
			}
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos++
				}
				break
			}
		} else { // if (val1>val2)
			k2++
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos++
				}
				break
			}
		}
	}
	return pos, buffer[:pos]
}

// http://en.wikipedia.org/wiki/Hamming_weight
func countBits(i uint64) int {
	i = i - ((i >> 1) & 0x5555555555555555)
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
	result := (((i + (i >> 4)) & 0xF0F0F0F0F0F0F0F) * 0x101010101010101) >> 56
	return int(result)
}

func highbits(x uint32) uint16 {
	return uint16(x >> 16)
}

func lowbits(x uint32) uint16 {
	return uint16(x & 0xFFFF)
}

func highlowbits(x uint32) (uint16, uint16) {
	return highbits(x), lowbits(x)
}

func fillArrayAND(bitmap1, bitmap2 []uint64, newCardinality int) []uint16 {
	pos := 0

	if len(bitmap1) != len(bitmap2) {
		panic("Bitmaps have different length - not supported.")
	}

	container := make([]uint16, newCardinality)
	for k := 0; k < len(bitmap1); k++ {
		bitset := bitmap1[k] & bitmap2[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k*64 + countBits(t-1)))
			pos++
			bitset ^= t
		}
	}

	return container
}

func fillArrayXOR(bitmap1, bitmap2 []uint64, newCardinality int) []uint16 {
	pos := 0

	if len(bitmap1) != len(bitmap2) {
		panic("Bitmaps have different length - not supported.")
	}

	container := make([]uint16, newCardinality)
	for k := 0; k < len(bitmap1); k++ {
		bitset := bitmap1[k] ^ bitmap2[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k*64 + countBits(t-1)))
			pos++
			bitset ^= t
		}
	}

	return container
}

// http://graphics.stanford.edu/~seander/bithacks.html#ZerosOnRightBinSearch
func trailingZeros(v uint64) int {
	if v&0x1 == 1 {
		return 0
	}

	c := 1

	if (v & 0xFFFFFFFF) == 0 {
		v = v >> 32
		c = c + 32
	}

	if (v & 0xFFFF) == 0 {
		v = v >> 16
		c = c + 16
	}
	if (v & 0xFF) == 0 {
		v = v >> 8
		c = c + 8
	}
	if (v & 0xF) == 0 {
		v = v >> 4
		c = c + 4
	}
	if (v & 0x3) == 0 {
		v = v >> 2
		c = c + 2
	}

	return c - int(v&0x1)
}

// Unite two sorted lists
func union2by2(set1 []uint16, length1 int,
	set2 []uint16, length2, bufferSize int) (int, []uint16) {

	if 0 == length2 {
		buffer := make([]uint16, length1)
		copy(buffer, set1)
		return length1, buffer
	}

	if 0 == length1 {
		buffer := make([]uint16, length2)
		copy(buffer, set2)
		return length2, buffer
	}

	buffer := make([]uint16, bufferSize)

	k1, k2, pos := 0, 0, 0

	for {
		if set1[k1] < set2[k2] {
			buffer[pos] = set1[k1]
			pos = pos + 1
			k1 = k1 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
		} else if set1[k1] == set2[k2] {
			buffer[pos] = set1[k1]
			pos = pos + 1
			k1 = k1 + 1
			k2 = k2 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		} else {
			buffer[pos] = set2[k2]
			pos = pos + 1
			k2 = k2 + 1
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		}
	}
	return pos, buffer[:pos]
}

// Compute the exclusive union of two sorted lists
func exclusiveUnion2by2(set1 []uint16, length1 int,
	set2 []uint16, length2, bufferSize int) (int, []uint16) {

	if 0 == length2 {
		buffer := make([]uint16, length1)
		copy(buffer, set1)
		return length1, buffer
	}

	if 0 == length1 {
		buffer := make([]uint16, length2)
		copy(buffer, set2)
		return length2, buffer
	}

	buffer := make([]uint16, bufferSize)

	k1, k2, pos := 0, 0, 0

	for {
		if set1[k1] < set2[k2] {
			buffer[pos] = set1[k1]
			pos = pos + 1
			k1 = k1 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
		} else if set1[k1] == set2[k2] {
			buffer[pos] = set1[k1]
			k1 = k1 + 1
			k2 = k2 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		} else {
			buffer[pos] = set2[k2]
			pos = pos + 1
			k2 = k2 + 1
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		}
	}
	return pos, buffer[:pos]
}
