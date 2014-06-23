package goroar

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"text/template"
)

type entry struct {
	key       uint16
	container container
}

type RoaringBitmap struct {
	containers []entry
}

// New creates a new RoaringBitmap
func New() *RoaringBitmap {
	containers := make([]entry, 0, 4)
	return &RoaringBitmap{containers}
}

// Add adds a uint32 value to the RoaringBitmap
func (rb *RoaringBitmap) Add(x uint32) {
	hb, lb := highlowbits(x)

	pos := rb.containerIndex(hb)
	if pos >= 0 {
		container := rb.containers[pos].container
		rb.containers[pos].container = container.add(lb)
	} else {
		ac := newArrayContainer()
		ac.add(lb)
		rb.increaseCapacity()

		loc := -pos - 1

		// insertion : shift the elements > x by one position to
		// the right and put x in it's appropriate place
		rb.containers = rb.containers[:len(rb.containers)+1]
		copy(rb.containers[loc+1:], rb.containers[loc:])
		e := entry{hb, ac}
		rb.containers[loc] = e
	}
}

// Contains checks whether the value in included, which is equivalent to checking
// if the corresponding bit is set (get in BitSet class).
func (rb *RoaringBitmap) Contains(i uint32) bool {
	pos := rb.containerIndex(highbits(i))
	if pos < 0 {
		return false
	}
	return rb.containers[pos].container.contains(lowbits(i))
}

// Cardinality returns the number of distinct integers (uint32) in the bitmap.
func (rb *RoaringBitmap) Cardinality() int {
	var cardinality int
	for _, entry := range rb.containers {
		cardinality = cardinality + entry.container.getCardinality()
	}
	return cardinality
}

// And computes the bitwise AND operation.
// The receiving RoaringBitmap is modified - the input one is not.
func (rb *RoaringBitmap) And(other *RoaringBitmap) {
	pos1 := 0
	pos2 := 0
	length1 := len(rb.containers)
	length2 := len(other.containers)

Main:
	for pos1 < length1 && pos2 < length2 {
		s1 := rb.keyAtIndex(pos1)
		s2 := other.keyAtIndex(pos2)
		for {
			if s1 < s2 {
				rb.removeAtIndex(pos1)
				length1 = length1 - 1
				if pos1 == length1 {
					break Main
				}
				s1 = rb.keyAtIndex(pos1)
			} else if s1 > s2 {
				pos2 = pos2 + 1
				if pos2 == length2 {
					break Main
				}
				s2 = other.keyAtIndex(pos2)
			} else {
				c := rb.containers[pos1].container.and(other.containers[pos2].container)

				if c.getCardinality() > 0 {
					rb.containers[pos1].container = c
					pos1 = pos1 + 1
				} else {
					rb.removeAtIndex(pos1)
					length1 = length1 - 1
				}
				pos2 = pos2 + 1
				if (pos1 == length1) || (pos2 == length2) {
					break Main
				}
				s1 = rb.keyAtIndex(pos1)
				s2 = other.keyAtIndex(pos2)
			}
		}
	}
	rb.resize(pos1)
}

// Or computes the bitwise OR operation.
// The receiving RoaringBitmap is modified - the input one is not.
func (rb *RoaringBitmap) Or(other *RoaringBitmap) {
	pos1, pos2 := 0, 0
	length1 := len(rb.containers)
	length2 := len(other.containers)

main:
	for pos1 < length1 && pos2 < length2 {
		s1 := rb.keyAtIndex(pos1)
		s2 := other.keyAtIndex(pos2)
		for {
			if s1 < s2 {
				pos1++
				if pos1 == length1 {
					break main
				}
				s1 = rb.keyAtIndex(pos1)
			} else if s1 > s2 {
				rb.insertAt(pos1, s2, other.containers[pos2].container)
				pos1++
				length1++
				pos2++
				if pos2 == length2 {
					break main
				}
				s2 = other.containers[pos2].key
			} else {
				rb.containers[pos1].container = rb.containers[pos1].container.or(other.containers[pos2].container)
				pos1++
				pos2++
				if pos1 == length1 || pos2 == length2 {
					break main
				}
				s1 = rb.containers[pos1].key
				s2 = other.containers[pos2].key
			}
		}
	}
	if pos1 == length1 {
		rb.containers = append(rb.containers, other.containers[pos2:length2]...)
	}
}

// Xor computes the bitwise XOR operation.
// The receiving RoaringBitmap is modified - the input one is not.
func (rb *RoaringBitmap) Xor(other *RoaringBitmap) {
	pos1, pos2 := 0, 0
	length1 := len(rb.containers)
	length2 := len(other.containers)

main:
	for pos1 < length1 && pos2 < length2 {
		s1 := rb.keyAtIndex(pos1)
		s2 := other.keyAtIndex(pos2)
		for {
			if s1 < s2 {
				pos1++
				if pos1 == length1 {
					break main
				}
				s1 = rb.keyAtIndex(pos1)
			} else if s1 > s2 {
				rb.insertAt(pos1, s2, other.containers[pos2].container)
				pos1++
				length1++
				pos2++
				if pos2 == length2 {
					break main
				}
				s2 = other.containers[pos2].key
			} else {
				c := rb.containers[pos1].container.xor(other.containers[pos2].container)
				if c.getCardinality() > 0 {
					rb.containers[pos1].container = c
					pos1++
				} else {
					rb.removeAtIndex(pos1)
					length1--
				}
				pos2++
				if pos1 == length1 || pos2 == length2 {
					break main
				}
				s1 = rb.containers[pos1].key
				s2 = other.containers[pos2].key
			}
		}
	}
	if pos1 == length1 {
		rb.containers = append(rb.containers, other.containers[pos2:length2]...)
	}
}

// AndNot computes the bitwise andNot operation (difference)
// The receiving RoaringBitmap is modified - the input one is not.
func (rb *RoaringBitmap) AndNot(other *RoaringBitmap) {
	pos1, pos2 := 0, 0
	length1 := len(rb.containers)
	length2 := len(other.containers)

main:
	for pos1 < length1 && pos2 < length2 {
		s1 := rb.keyAtIndex(pos1)
		s2 := other.keyAtIndex(pos2)
		for {
			if s1 < s2 {
				pos1++
				if pos1 == length1 {
					break main
				}
				s1 = rb.keyAtIndex(pos1)
			} else if s1 > s2 {
				pos2++
				if pos2 == length2 {
					break main
				}
				s2 = other.containers[pos2].key
			} else {
				c := rb.containers[pos1].container.andNot(other.containers[pos2].container)
				if c.getCardinality() > 0 {
					rb.containers[pos1].container = c
					pos1++
				} else {
					rb.removeAtIndex(pos1)
					length1--
				}
				pos2++
				if pos1 == length1 || pos2 == length2 {
					break main
				}
				s1 = rb.containers[pos1].key
				s2 = other.containers[pos2].key
			}
		}
	}
}

// Iterator returns an iterator over the RoaringBitmap which can be used with "for range".
func (rb *RoaringBitmap) Iterator() <-chan uint32 {
	ch := make(chan uint32)
	go func() {
		// iterate over data
		for _, entry := range rb.containers {
			hs := uint32(entry.key) << 16
			switch typedContainer := entry.container.(type) {
			case *arrayContainer:
				pos := 0
				for pos < typedContainer.cardinality {
					ls := typedContainer.content[pos]
					pos = pos + 1
					ch <- (hs | uint32(ls))
				}
			case *bitmapContainer:
				i := typedContainer.nextSetBit(0)
				for i >= 0 {
					ch <- (hs | uint32(i))
					i = typedContainer.nextSetBit(i + 1)
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (rb *RoaringBitmap) String() string {
	var buffer bytes.Buffer
	name := []byte("RoaringBitmap[")

	buffer.Write(name)
	for val := range rb.Iterator() {
		buffer.WriteString(strconv.Itoa(int(val)))
		buffer.WriteString(", ")
	}
	if buffer.Len() > len(name) {
		buffer.Truncate(buffer.Len() - 2) // removes the last ", "
	}
	buffer.WriteString("]")
	return buffer.String()
}

// Stats prints statistics about the Roaring Bitmap's internals.
func (rb *RoaringBitmap) Stats() {
	const output = `* Roaring Bitmap Stats *
Cardinality: {{.Cardinality}}
Size uncompressed: {{.UncompressedSize}} bytes
Size compressed: {{.CompressedSize}} bytes
Number of containers: {{.TotalContainers}}
    {{.TotalAC}} ArrayContainers
    {{.TotalBC}} BitmapContainers
Average entries per ArrayContainer: {{.AverageAC}}
Max entries per ArrayContainer: {{.MaxAC}}
`
	type stats struct {
		Cardinality, TotalContainers, TotalAC, TotalBC int
		AverageAC, MaxAC                               string
		CompressedSize, UncompressedSize               int
	}

	var totalAC, totalBC, totalCardinalityAC int
	var maxAC int

	for _, c := range rb.containers {
		switch typedContainer := c.container.(type) {
		case *arrayContainer:
			if typedContainer.cardinality > maxAC {
				maxAC = typedContainer.cardinality
			}
			totalCardinalityAC += typedContainer.cardinality
			totalAC++
		case *bitmapContainer:
			totalBC++
		default:
		}
	}

	s := new(stats)
	s.Cardinality = rb.Cardinality()
	s.TotalContainers = len(rb.containers)
	s.TotalAC = totalAC
	s.TotalBC = totalBC
	s.CompressedSize = rb.SizeInBytes()
	s.UncompressedSize = rb.Cardinality() * 4

	if totalCardinalityAC > 0 {
		s.AverageAC = string(totalCardinalityAC / totalAC)
		s.MaxAC = string(maxAC)
	} else {
		s.AverageAC = "--"
		s.MaxAC = "--"
	}

	t := template.Must(template.New("stats").Parse(output))
	if err := t.Execute(os.Stdout, s); err != nil {
		log.Println("RoaringBitmap stats: ", err)
	}
}

func (rb *RoaringBitmap) SizeInBytes() int {
	size := 8
	for _, c := range rb.containers {
		size += 2 + c.container.sizeInBytes()
	}
	return size
}

func (rb *RoaringBitmap) resize(newLength int) {
	for i := newLength; i < len(rb.containers); i++ {
		rb.containers[i] = entry{}
	}
	rb.containers = rb.containers[:newLength]
}

func (rb *RoaringBitmap) keyAtIndex(pos int) uint16 {
	return rb.containers[pos].key
}

func (rb *RoaringBitmap) removeAtIndex(i int) {
	copy(rb.containers[i:], rb.containers[i+1:])
	rb.containers[len(rb.containers)-1] = entry{}
	rb.containers = rb.containers[:len(rb.containers)-1]
}

func (rb *RoaringBitmap) insertAt(i int, key uint16, c container) {
	rb.containers = append(rb.containers, entry{})
	copy(rb.containers[i+1:], rb.containers[i:])
	rb.containers[i] = entry{key, c}
}

func (rb *RoaringBitmap) containerIndex(key uint16) int {
	length := len(rb.containers)

	if length == 0 || rb.containers[length-1].key == key {
		return length - 1
	}

	return searchContainer(rb.containers, length, key)
}

func searchContainer(containers []entry, length int, key uint16) int {
	low := 0
	high := length - 1

	for low <= high {
		middleIndex := (low + high) >> 1
		middleValue := containers[middleIndex].key

		switch {
		case middleValue < key:
			low = middleIndex + 1
		case middleValue > key:
			high = middleIndex - 1
		default:
			return middleIndex
		}
	}
	return -(low + 1)
}

// increaseCapacity increases the slice capacity keeping the same length.
func (rb *RoaringBitmap) increaseCapacity() {
	length := len(rb.containers)
	if length+1 > cap(rb.containers) {
		var newCapacity int
		if length < 1024 {
			newCapacity = 2 * (length + 1)
		} else {
			newCapacity = 5 * (length + 1) / 4
		}

		newSlice := make([]entry, length, newCapacity)
		copy(newSlice, rb.containers)

		// increasing the length by 1
		rb.containers = newSlice
	}
}

// And computes the bitwise AND operation on two RoaringBitmaps.
// The input bitmaps are not modified.
func And(x1, x2 *RoaringBitmap) *RoaringBitmap {
	panic("Not implemented")
}
