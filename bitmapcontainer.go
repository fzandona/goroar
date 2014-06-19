package goroar

const (
	bitmapContainerMaxCapacity = uint32(1 << 16)
	one                        = uint64(1)
)

type bitmapContainer struct {
	cardinality int
	bitmap      []uint64
}

var _ container = (*bitmapContainer)(nil)

func newBitmapContainer() *bitmapContainer {
	return &bitmapContainer{0, make([]uint64, bitmapContainerMaxCapacity/64)}
}

func (bc *bitmapContainer) loadData(ac *arrayContainer) {
	bc.cardinality = ac.cardinality
	for i := 0; i < ac.cardinality; i++ {
		bc.bitmap[uint32(ac.content[i])/64] |= one << (ac.content[i] % 64)
	}
}

func (bc *bitmapContainer) add(i uint16) {
	x := uint32(i)
	index := x / 64
	mod := x % 64
	previous := bc.bitmap[index]
	bc.bitmap[index] |= one << mod
	bc.cardinality += int((previous ^ bc.bitmap[index]) >> mod)
}

func (bc *bitmapContainer) and(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return bc.andArray(oc)
	case *bitmapContainer:
		return bc.andBitmap(oc)
	}
	return nil
}

func (bc *bitmapContainer) andArray(value2 *arrayContainer) *arrayContainer {
	answer := make([]uint16, value2.cardinality)

	cardinality := 0
	for k := 0; k < value2.cardinality; k++ {
		if bc.contains(value2.content[k]) {
			answer[cardinality] = value2.content[k]
			cardinality++
		}
	}

	return &arrayContainer{cardinality, answer[:cardinality]}
}

func (bc *bitmapContainer) andBitmap(value2 *bitmapContainer) container {
	newCardinality := 0
	for k, v := range bc.bitmap {
		newCardinality += countBits(v & value2.bitmap[k])
	}

	if newCardinality > arrayContainerMaxSize {
		answer := newBitmapContainer()
		for k, v := range bc.bitmap {
			answer.bitmap[k] = v & value2.bitmap[k]
		}
		answer.cardinality = newCardinality
		return answer

	}
	content := fillArrayAND(bc.bitmap, value2.bitmap, newCardinality)
	return &arrayContainer{newCardinality, content}
}

func (bc *bitmapContainer) or(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return bc.orArray(oc)
	case *bitmapContainer:
		return bc.orBitmap(oc)
	}
	return nil
}

func (bc *bitmapContainer) orArray(ac *arrayContainer) *bitmapContainer {
	answer := bc.clone()
	for i := 0; i < ac.cardinality; i++ {
		answer.add(ac.content[i])
	}
	return answer
}

func (bc *bitmapContainer) orBitmap(other *bitmapContainer) container {
	answer := newBitmapContainer()

	for i := 0; i < len(bc.bitmap); i++ {
		answer.bitmap[i] = bc.bitmap[i] | other.bitmap[i]
		answer.cardinality = answer.cardinality + countBits(answer.bitmap[i])
	}
	return answer
}
func (bc *bitmapContainer) xor(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return bc.xorArray(oc)
	case *bitmapContainer:
		return bc.xorBitmap(oc)
	}
	return nil
}

func (bc *bitmapContainer) xorArray(ac *arrayContainer) container {
	answer := bc.clone()
	for i := 0; i < ac.cardinality; i++ {
		v := ac.content[i]
		mod := v % 64
		index := v / 64
		shift := one << v
		answer.cardinality += 1 - 2*int((answer.bitmap[index]&shift)>>mod)
		answer.bitmap[index] ^= shift

	}
	if answer.cardinality <= arrayContainerMaxSize {
		return answer.toArrayContainer()
	}
	return answer
}

func (bc *bitmapContainer) xorBitmap(other *bitmapContainer) container {
	answer := newBitmapContainer()

	for i := 0; i < len(bc.bitmap); i++ {
		answer.bitmap[i] = bc.bitmap[i] ^ other.bitmap[i]
		answer.cardinality = answer.cardinality + countBits(answer.bitmap[i])
	}

	if answer.cardinality <= arrayContainerMaxSize {
		return answer.toArrayContainer()
	}
	return answer
}
func (bc *bitmapContainer) andNot(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return bc.andNotArray(oc)
	case *bitmapContainer:
		return bc.andNotBitmap(oc)
	}
	return nil
}

func (bc *bitmapContainer) andNotArray(ac *arrayContainer) container {
	answer := bc.clone()
	for i := 0; i < ac.cardinality; i++ {
		v := ac.content[i]
		mod := v % 64
		index := v / 64
		shift := one << v
		answer.bitmap[index] = answer.bitmap[index] & (^shift)
		answer.cardinality -= int((answer.bitmap[index] ^ bc.bitmap[index]) >> mod)
	}
	if answer.cardinality <= arrayContainerMaxSize {
		return answer.toArrayContainer()
	}
	return answer
}

func (bc *bitmapContainer) andNotBitmap(other *bitmapContainer) container {
	answer := newBitmapContainer()

	for i := 0; i < len(bc.bitmap); i++ {
		answer.bitmap[i] = bc.bitmap[i] & (^other.bitmap[i])
		answer.cardinality = answer.cardinality + countBits(answer.bitmap[i])
	}

	if answer.cardinality <= arrayContainerMaxSize {
		return answer.toArrayContainer()
	}
	return answer
}

func (bc *bitmapContainer) clone() *bitmapContainer {
	bitmap := make([]uint64, len(bc.bitmap))
	copy(bitmap, bc.bitmap)
	return &bitmapContainer{bc.cardinality, bitmap}
}

func (bc *bitmapContainer) contains(x uint16) bool {
	return bc.bitmap[uint32(x)/64]&(one<<(x%64)) != 0
}

// nextSetBit finds the index of the next set bit greater or equal to i.
// It returns -1 if none is found.
func (bc *bitmapContainer) nextSetBit(i int) int {
	x := i / 64
	if x >= len(bc.bitmap) {
		return -1
	}

	w := bc.bitmap[x]
	w = w >> (uint(i) % 64)
	if w != 0 {
		return i + trailingZeros(w)
	}

	x = x + 1
	for ; x < len(bc.bitmap); x++ {
		if bc.bitmap[x] != 0 {
			return x*64 + trailingZeros(bc.bitmap[x])
		}
	}

	return -1
}

func (bc *bitmapContainer) getCardinality() int {
	return bc.cardinality
}

func (bc *bitmapContainer) toArrayContainer() *arrayContainer {
	container := make([]uint16, bc.cardinality)
	pos := 0
	for k := 0; k < len(bc.bitmap); k++ {
		bitset := bc.bitmap[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k*64 + countBits(t-1)))
			pos++
			bitset ^= t
		}
	}

	return &arrayContainer{bc.cardinality, container}
}
