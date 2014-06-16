package goroar

import (
	"math/rand"
	"testing"
)

func TestAdd(t *testing.T) {
	rb := New()
	rb.Add(42)

	switch ac := rb.containers[0].container.(type) {
	case *arrayContainer:
		if ac.cardinality != 1 || ac.content[0] != 42 {
			t.Errorf("Error on cardinality and content")
		}
	default:
		t.Errorf("Wrong container type: %T", ac)
	}
}

func TestAdd_2(t *testing.T) {
	rb := New()
	// creates 5 containers with keys: 1, 2, 4, 8, 16
	for shift := uint(16); shift <= 20; shift++ {
		rb.Add(uint32((1 << shift) + 1))
	}

	if len(rb.containers) != 5 {
		t.Errorf("Containers length: %d, want: 5", len(rb.containers))
	}

	for _, container := range rb.containers {
		switch ac := container.container.(type) {
		case *arrayContainer:
			if ac.cardinality != 1 || ac.content[0] != 1 {
				t.Errorf("Error on cardinality and/or content")
			}
		default:
			t.Errorf("Wrong container type: %T", ac)
		}
	}
}

func TestAdd_3(t *testing.T) {
	rb := New()
	for i := 0; i < 4096; i++ {
		rb.Add(uint32(i))
	}

	switch ac := rb.containers[0].container.(type) {
	case *arrayContainer:
		if ac.cardinality != 4096 {
			t.Errorf("Cardinality: %d, want: 4096", ac.cardinality)
		}
		for i := 0; i < 4096; i++ {
			if ac.content[i] != uint16(i) {
				t.Errorf("ArrayContent content: %d, want: %d", ac.content[i], i)
			}
		}

	default:
		t.Errorf("Wrong container type: %T", ac)

	}

	rb.Add(uint32(0))
	switch rb.containers[0].container.(type) {
	case *bitmapContainer:
		t.Error("BitmapContainer found, want ArrayContainer")
	}

	rb.Add(uint32(4096))
	switch bc := rb.containers[0].container.(type) {
	case *bitmapContainer:
		if bc.cardinality != 4097 {
			t.Errorf("Cardinality: %d, want: 4097", bc.cardinality)
		}
	default:
		t.Errorf("Wrong container type: %T", bc)
	}
}

func TestAdd_4(t *testing.T) {
	rb := New()
	for i := 0; i < 1000; i++ {
		rb.Add(1234567)
	}
	size := rb.Cardinality()
	if size != 1 {
		t.Errorf("Repeated add: %d, want: 1", size)
	}
}

func TestCardinality(t *testing.T) {
	rb := RoaringBitmap{}
	count := 10000
	for i := 0; i < count; i++ {
		rb.Add(uint32(i))
	}
	size := rb.Cardinality()
	if size != count {
		t.Errorf("Cardinality: %d, want: %d", size, count)
	}
}

func TestRoaringBitmapContains(t *testing.T) {
	rb := New()
	r := rand.New(rand.NewSource(42))
	r.Seed(42)
	iterations := 123456
	buffer := make([]uint32, iterations)

	for i := 0; i < iterations; i++ {
		buffer[i] = uint32(r.Int31())
		rb.Add(buffer[i])
	}

	for i := 0; i < iterations; i++ {
		value := buffer[i]
		if !rb.Contains(value) {
			t.Errorf("Bitmap should contain: %d", value)
			break
		}
	}

	for i := 0; i < 10; i++ {
		value := uint32(r.Int31())
		if rb.Contains(value) {
			t.Errorf("Bitmap should not contain: %d", value)
			break
		}
	}
}

func TestIterator(t *testing.T) {
	rb := New()
	iterations := uint32(99999)

	for i := uint32(0); i < iterations; i += 2 {
		rb.Add(i)
	}

	pos := 0
	for val := range rb.Iterator() {
		if val != uint32(pos) {
			t.Errorf("Iterator: %d, want: %d", val, pos)
			break
		}
		pos = pos + 2
	}

	if iterations != uint32(pos-1) {
		t.Errorf("Wrong iterations: %d, want: %d", pos-1, iterations)
	}
}

func TestIterator_2(t *testing.T) {
	rb := New()

	iterations := 123456
	for i := 0; i < iterations; i++ {
		rb.Add(uint32(i))
	}

	answer := 0
	for val := range rb.Iterator() {
		if val != uint32(answer) {
			t.Errorf("Iterator: %d, want: %d", val, answer)
		}
		answer = answer + 1
	}

	if iterations != answer {
		t.Errorf("Wrong iterations: %d, want: %d", answer, iterations)
	}
}

func TestResize(t *testing.T) {
	rb := New()
	containers := make([]entry, 10)
	for i := 0; i < 10; i++ {
		containers[i].key = uint16(i)
	}

	rb.containers = containers
	pos := 5
	rb.resize(pos)
	if len(rb.containers) != pos {
		t.Errorf("Length: %d, want: %d", len(rb.containers), pos-1)
	}
	for i := 0; i < len(rb.containers); i++ {
		if rb.containers[i].key != uint16(i) {
			t.Errorf("Resize error, key %d, want: %d", rb.containers[i].key, i)
		}
	}
}

func TestRemoveAtIndex(t *testing.T) {
	rb := New()

	// creates 5 containers with keys: 1, 2, 4, 8, 16
	for shift := uint(16); shift <= 20; shift++ {
		rb.Add(uint32((1 << shift) + 1))
	}

	rb.removeAtIndex(2)
	results := []uint16{1, 2, 8, 16}
	for i := 0; i < 4; i++ {
		if rb.containers[i].key != results[i] {
			t.Error("Remove at index: %d, want: %d", rb.containers[i].key, results[i])
		}
	}

	rb.removeAtIndex(3)
	results = []uint16{1, 2, 8}
	for i := 0; i < 3; i++ {
		if rb.containers[i].key != results[i] {
			t.Error("Remove at index: %d, want: %d", rb.containers[i].key, results[i])
		}
	}

	rb.removeAtIndex(0)
	results = []uint16{2, 8}
	for i := 0; i < 2; i++ {
		if rb.containers[i].key != results[i] {
			t.Error("Remove at index: %d, want: %d", rb.containers[i].key, results[i])
		}
	}
}

func TestInPlaceAnd_1(t *testing.T) {
	rb1 := New()
	rb2 := New()

	for i := 0; i < 100; i += 2 {
		rb1.Add(uint32(i))
		rb2.Add(uint32(i + 1))
	}

	rb1.And(rb2)

	if len(rb1.containers) != 0 {
		t.Errorf("And length: %d, want: 0", len(rb1.containers))
	}

	for i := 0; i < 100; i++ {
		rb1.Add(uint32(i))
		rb2.Add(uint32(i))
	}
	rb1.And(rb2)
	for i := 0; i < 100; i++ {
		if !rb1.Contains(uint32(i)) {
			t.Errorf("RB should contain: %d", i)
		}
	}
}

func TestInPlaceAnd_2(t *testing.T) {
	rb1 := New()
	rb2 := New()
	iterations := 500

	for i := 0; i < iterations; i++ {
		rb1.Add(uint32(i))
		rb2.Add(uint32(i))
	}

	rb1.And(rb2)

	for i := 0; i < iterations; i++ {
		if !rb1.Contains(uint32(i)) {
			t.Errorf("RB should contain: %d", i)
		}
	}
}

func TestInPlaceAnd_3(t *testing.T) {
	rb1 := New()
	rb2 := New()
	var buffer []uint32
	count := 10
	// creates 5 containers with keys: 1, 2, 4, 8, 16
	for shift := uint(16); shift <= 20; shift++ {
		base := uint32(1 << shift)
		for i := 0; i < count; i++ {
			value := base + uint32(i)
			rb1.Add(value)
			rb2.Add(value)
			buffer = append(buffer, value)
		}
	}

	rb1.And(rb2)

	pos := 0
	for v := range rb1.Iterator() {
		if v != buffer[pos] {
			t.Errorf("And: %d, want: %d", v, buffer[pos])
		}
		pos = pos + 1
	}
}

func TestInPlaceOr_1(t *testing.T) {
	rb1 := New()
	rb2 := New()
	count := 100

	for i := 0; i < count; i += 2 {
		rb1.Add(uint32(i))
		rb2.Add(uint32(i + 1))
	}

	rb1.Or(rb2)

	var pos uint32
	for v := range rb1.Iterator() {
		if v != pos {
			t.Errorf("Or: %d, want: %d", v, pos)
			break
		}
		pos++
	}
}

func TestInPlaceOr_2(t *testing.T) {
	rb1 := New()
	rb2 := New()
	var buffer []uint32
	count := 5000
	// creates 5 containers with keys: 1, 2, 4, 8, 16
	for shift := uint(16); shift <= 20; shift++ {
		base := uint32(1 << shift)
		for i := 0; i < count; i++ {
			value := base + uint32(i)
			rb1.Add(value)
			rb2.Add(value)
			buffer = append(buffer, value)
		}
	}

	rb1.Or(rb2)

	pos := 0
	for v := range rb1.Iterator() {
		if v != buffer[pos] {
			t.Errorf("And: %d, want: %d", v, buffer[pos])
		}
		pos = pos + 1
	}
}

func TestInPlaceOr_3(t *testing.T) {
	rb1 := New()
	rb2 := New()
	var buffer []uint32
	count := 500
	// creates 5 containers with keys: 1, 2, 4, 8, 16
	for shift := uint(16); shift <= 20; shift++ {
		base := uint32(1 << shift)
		for i := 0; i < count; i++ {
			value := base + uint32(i)
			rb1.Add(value)
			buffer = append(buffer, value)
		}
	}

	rb1.Or(rb2)

	pos := 0
	for v := range rb1.Iterator() {
		if v != buffer[pos] {
			t.Errorf("And: %d, want: %d", v, buffer[pos])
		}
		pos = pos + 1
	}
}

func TestInPlaceOr_4(t *testing.T) {
	rb1 := New()
	rb2 := New()
	var buffer []uint32
	count := 500
	// creates 5 containers with keys: 1, 2, 4, 8, 16
	for shift := uint(16); shift <= 20; shift++ {
		base := uint32(1 << shift)
		for i := 0; i < count; i++ {
			value := base + uint32(i)
			rb1.Add(value)
			buffer = append(buffer, value)
		}
	}

	rb2.Add(0)
	rb1.Or(rb2)

	buffer = append(buffer, 0)
	copy(buffer[1:], buffer)
	buffer[0] = 0

	pos := 0
	for v := range rb1.Iterator() {
		if v != buffer[pos] {
			t.Errorf("And: %d, want: %d", v, buffer[pos])
		}
		pos = pos + 1
	}
}

func TestInPlaceOr_5(t *testing.T) {
	rb1 := New()
	rb2 := New()

	rb1.Add(uint32(0))
	for i := 1; i < 5000; i++ {
		rb2.Add(uint32(i))
	}

	rb1.Or(rb2)

	if rb1.Cardinality() != 5000 {
		t.Errorf("Cardinality: %d, want: %d", rb1.Cardinality(), 5000)
	}

	pos := 0
	for v := range rb1.Iterator() {
		if int(v) != pos {
			t.Errorf("And: %d, want: %d", v, pos)
		}
		pos = pos + 1
	}
}

func TestString(t *testing.T) {
	rb := New()
	iterations := 5

	str := "RoaringBitmap[]"
	if rb.String() != str {
		t.Errorf("String: %s, want: %s", rb.String(), str)
	}

	for i := 0; i < iterations; i++ {
		rb.Add(uint32(i))
	}

	str = "RoaringBitmap[0, 1, 2, 3, 4]"
	if rb.String() != str {
		t.Errorf("String: %s, want: %s", rb.String(), str)
	}
}

func TestInsertAt(t *testing.T) {
	rb1 := New()
	count := 10
	// creates 5 containers with keys: 1, 2, 4, 8, 16
	for shift := uint(16); shift <= 20; shift++ {
		base := uint32(1 << shift)
		for i := 0; i < count; i++ {
			value := base + uint32(i)
			rb1.Add(value)
		}
	}

	ac := newArrayContainer()
	ac.add(uint16(123))
	key := uint16(0)

	rb1.insertAt(0, key, ac)
	if e := rb1.containers[0]; e.key != 0 && e.container.(*arrayContainer).content[0] != uint16(123) {
		t.Errorf("insertAt error")
	}

	rb1.insertAt(len(rb1.containers), key, ac)
	if e := rb1.containers[len(rb1.containers)-1]; e.key != 0 && e.container.(*arrayContainer).content[0] != uint16(123) {
		t.Errorf("insertat error")
	}

	rb1.insertAt(2, key, ac)
	if e := rb1.containers[2]; e.key != 0 && e.container.(*arrayContainer).content[0] != uint16(123) {
		t.Errorf("insertAt error")
	}
}
