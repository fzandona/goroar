package goroar

const (
	arrayContainerInitSize = 4
	arrayContainerMaxSize  = 4096
)

type arrayContainer struct {
	cardinality int
	content     []uint16
}

var _ container = (*arrayContainer)(nil)

func newArrayContainer() *arrayContainer {
	content := make([]uint16, arrayContainerInitSize)
	return &arrayContainer{0, content}
}

func newArrayContainerWithCapacity(capacity int) *arrayContainer {
	content := make([]uint16, capacity)
	return &arrayContainer{0, content}
}

func newArrayContainerRunOfOnes(firstOfRun, lastOfRun int) *arrayContainer {
	valuesInRange := lastOfRun - firstOfRun + 1
	content := make([]uint16, valuesInRange)
	for i := 0; i < valuesInRange; i++ {
		content[i] = uint16(firstOfRun + i)
	}
	return &arrayContainer{int(valuesInRange), content}
}

// ArrayContainer add returns false if it is time to switch to a
// BitmapContainer - the integer is not added in this case.
// TODO: check performance against just moving to BC before
//      verifying the integer is distinct or not.
func (ac *arrayContainer) add(x uint16) bool {
	if ac.cardinality == 0 || x > ac.content[ac.cardinality-1] {
		if ac.cardinality >= arrayContainerMaxSize {
			return false
		}
		if ac.cardinality >= len(ac.content) {
			ac.increaseCapacity()
		}
		ac.content[ac.cardinality] = x
		ac.cardinality++
		return true
	}

	loc := binarySearch(ac.content, ac.cardinality, x)
	if loc < 0 {
		if ac.cardinality >= arrayContainerMaxSize {
			return false
		}
		if ac.cardinality >= len(ac.content) {
			ac.increaseCapacity()
		}
		loc = -loc - 1
		// insertion : shift the elements > x by one position to
		// the right and put x in it's appropriate place
		copy(ac.content[loc+1:], ac.content[loc:])
		ac.content[loc] = x
		ac.cardinality++
	}
	return true
}

func (ac *arrayContainer) and(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return ac.andArray(oc)
	case *bitmapContainer:
		return ac.andBitmap(oc)
	}
	return nil
}

func (ac *arrayContainer) andArray(value2 *arrayContainer) *arrayContainer {
	value1 := ac

	cardinality, content := intersect2by2(value1.content,
		value1.cardinality, value2.content,
		value2.cardinality)

	return &arrayContainer{cardinality, content}
}

func (ac *arrayContainer) andBitmap(bc *bitmapContainer) *arrayContainer {
	return bc.andArray(ac)
}

func (ac *arrayContainer) or(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return ac.orArray(oc)
	case *bitmapContainer:
		return ac.orBitmap(oc)
	}
	return nil
}

func (ac *arrayContainer) orArray(other *arrayContainer) container {
	totalCardinality := ac.cardinality + other.cardinality
	if totalCardinality > arrayContainerMaxSize {
		bc := newBitmapContainer()
		for i := 0; i < other.cardinality; i++ {
			bc.add(other.content[i])
		}
		for i := 0; i < ac.cardinality; i++ {
			bc.add(ac.content[i])
		}
		if bc.cardinality <= arrayContainerMaxSize {
			return bc.toArrayContainer()
		}
		return bc
	}
	answer := arrayContainer{}
	pos, content := union2by2(ac.content, ac.cardinality, other.content, other.cardinality, totalCardinality)
	answer.cardinality = pos
	answer.content = content
	return &answer
}

func (ac *arrayContainer) orBitmap(bc *bitmapContainer) container {
	return bc.or(ac)
}

func (ac *arrayContainer) andNot(value2 *arrayContainer) *arrayContainer {
	cardinality, content := difference(ac.content, ac.cardinality,
		value2.content, value2.cardinality)

	return &arrayContainer{cardinality, content}
}

func (ac *arrayContainer) andNotBitmap(value2 *bitmapContainer) *arrayContainer {
	content := make([]uint16, ac.cardinality)

	pos := 0
	for k := 0; k < ac.cardinality; k++ {
		if !value2.contains(ac.content[k]) {
			content[pos] = ac.content[k]
			pos++
		}
	}

	return &arrayContainer{pos, content[:pos]}
}

func (ac *arrayContainer) contains(x uint16) bool {
	return binarySearch(ac.content, ac.cardinality, x) >= 0
}

func (ac *arrayContainer) clear() {
	ac.content = make([]uint16, arrayContainerInitSize)
	ac.cardinality = 0
}

func (ac *arrayContainer) toBitmapContainer() *bitmapContainer {
	bc := newBitmapContainer()
	bc.loadData(ac)
	return bc
}

func (ac *arrayContainer) getCardinality() int {
	return ac.cardinality
}

func (ac *arrayContainer) arraySizeInBytes() int {
	return ac.cardinality * 2
}

func (ac *arrayContainer) increaseCapacity() {
	length := len(ac.content)
	var newLength int
	switch {
	case length < 64:
		newLength = length * 2
	case length < 1024:
		newLength = length * 3 / 2
	default:
		newLength = length * 5 / 4
	}
	if newLength > arrayContainerMaxSize {
		newLength = arrayContainerMaxSize
	}
	newSlice := make([]uint16, newLength)
	copy(newSlice, ac.content)
	ac.content = newSlice
}