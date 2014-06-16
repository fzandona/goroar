package goroar

import "testing"

func TestNewArrayContainer(t *testing.T) {
	ac := newArrayContainer()
	if ac.cardinality != 0 {
		t.Errorf("Cardinality = %d, want %d", ac.cardinality, 0)
	}
	if cap(ac.content) != arrayContainerInitSize {
		t.Errorf("Content = %d, want %d", cap(ac.content), arrayContainerInitSize)
	}
}

func TestNewArrayContainerWithCapacity(t *testing.T) {
	capacity := 42
	ac := newArrayContainerWithCapacity(capacity)
	if ac.cardinality != 0 {
		t.Errorf("Cardinality = %d, want %d", ac.cardinality, 0)
	}
	if cap(ac.content) != capacity {
		t.Errorf("Content = %d, want %d", cap(ac.content), capacity)
	}
}

func TestNewArrayContainerRunOfOnes(t *testing.T) {
	ac := newArrayContainerRunOfOnes(1, 4)
	if ac.cardinality != 4 {
		t.Errorf("Cardinality = %d, want %d", ac.cardinality, 4)
	}
	if (ac.content[0] != 1) ||
		(ac.content[1] != 2) ||
		(ac.content[2] != 3) ||
		(ac.content[3] != 4) {
		t.Errorf("Content error = %v", ac.content)
	}
}

func testArrayIncreaseCapacity(initialCapacity, expectedCapacity int, t *testing.T) {
	ac := newArrayContainerWithCapacity(initialCapacity)
	ac.increaseCapacity()
	if len(ac.content) != expectedCapacity {
		t.Errorf("Capacity: %d, want: %d", len(ac.content), expectedCapacity)
	}
}

func TestArrayIncreaseCapacity1(t *testing.T) {
	testArrayIncreaseCapacity(4, 8, t)
}

func TestArrayIncreaseCapacity2(t *testing.T) {
	testArrayIncreaseCapacity(1000, 1500, t)
}

func TestArrayIncreaseCapacity3(t *testing.T) {
	testArrayIncreaseCapacity(2000, 2500, t)
}

func TestArrayIncreaseCapacity4(t *testing.T) {
	testArrayIncreaseCapacity(4000, arrayContainerMaxSize, t)
}

func TestArrayContainerAdd(t *testing.T) {
	ac := newArrayContainer()
	ac.add(uint16(0))
	ac.add(uint16(2))

	if ac.content[0] != uint16(0) &&
		ac.content[1] != uint16(2) {
		t.Errorf("Wrong add: %d, %d, want: %d, %d", ac.content[0], ac.content[1],
			0, 2)
	}

	if ac.cardinality != 2 {
		t.Errorf("Cardinality: %d, want %d", ac.cardinality, 2)
	}

	ac.add(uint16(1))
	if ac.content[0] != uint16(0) &&
		ac.content[1] != uint16(1) &&
		ac.content[2] != uint16(2) {
		t.Errorf("Wrong add: %d, %d %d, want: %d, %d %d",
			ac.content[0],
			ac.content[1],
			ac.content[2],
			0, 1, 2)
	}

	if ac.cardinality != 3 {
		t.Errorf("Cardinality: %d, want %d", ac.cardinality, 3)
	}
}

func TestArrayContainerAddLong(t *testing.T) {
	ac := newArrayContainer()

	for i := 0; i < 4000; i++ {
		if i%2 == 0 {
			ac.add(uint16(i))
		}
	}

	for i := 0; i < 4000; i++ {
		if i%2 != 0 {
			ac.add(uint16(i))
		}
	}

	for i := 0; i < 4000; i++ {
		if ac.content[i] != uint16(i) {
			t.Errorf("Add: %d, want: %d", ac.content[i], i)
		}
	}
}

func TestAnd(t *testing.T) {
	ac1 := newArrayContainerWithCapacity(10)
	ac2 := newArrayContainerWithCapacity(1024)

	for i := 0; i < 10; i++ {
		ac1.add(uint16(i * 100))
	}

	for i := 0; i < 1024; i++ {
		ac2.add(uint16(i))
	}

	answer := ac1.andArray(ac2)
	if answer.cardinality != 10 {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, 10)
	}
	for i, v := range answer.content {
		if uint16(i*100) != v {
			t.Errorf("Got: %d, want: %d", v, i)
		}
	}
}

func TestAndBitmap(t *testing.T) {
	ac := newArrayContainerWithCapacity(10)
	bc := newBitmapContainer()

	for i := 0; i < 10; i++ {
		ac.add(uint16(i * 100))
	}

	for i := 0; i < 1024; i++ {
		bc.add(uint16(i))
	}

	answer := ac.andBitmap(bc)
	if answer.cardinality != 10 {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, 10)
	}
	for i, v := range answer.content {
		if uint16(i*100) != v {
			t.Errorf("Got: %d, want: %d", v, i)
		}
	}
}

func TestAndNot(t *testing.T) {
	ac1 := newArrayContainerWithCapacity(10)
	ac2 := newArrayContainerWithCapacity(10)

	for i := 0; i < 10; i++ {
		ac1.add(uint16(i))
	}

	for i := 0; i < 10; i++ {
		ac2.add(uint16(i + 5))
	}

	answer := ac1.andNot(ac2)
	if answer.cardinality != 5 {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, 10)
	}
	for i, v := range answer.content {
		if uint16(i) != v {
			t.Errorf("Got: %d, want: %d", v, i)
		}
	}
}

func TestAndNot_2(t *testing.T) {
	ac1 := newArrayContainerWithCapacity(10)
	ac2 := newArrayContainer()

	for i := 0; i < 10; i++ {
		ac1.add(uint16(i))
	}

	answer := ac1.andNot(ac2)
	if answer.cardinality != ac1.cardinality {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, ac1.cardinality)
	}
	for i, v := range answer.content {
		if ac1.content[i] != v {
			t.Errorf("Got: %d, want: %d", v, i)
		}
	}
}

func TestAndNot_3(t *testing.T) {
	ac1 := newArrayContainer()
	ac2 := newArrayContainerWithCapacity(10)

	for i := 0; i < 10; i++ {
		ac2.add(uint16(i))
	}

	answer := ac1.andNot(ac2)
	if answer.cardinality != 0 {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, ac1.cardinality)
	}

	if len(answer.content) != 0 {
		t.Errorf("Got: %d, want: %d", answer.content, 0)
	}
}

func TestAndNot_4(t *testing.T) {
	ac1 := newArrayContainerWithCapacity(10)
	ac2 := newArrayContainerWithCapacity(10)

	for i := 1; i <= 10; i++ {
		ac1.add(uint16(i))
	}

	for i := 1; i <= 10; i++ {
		ac2.add(uint16(i * 20))
	}

	answer := ac1.andNot(ac2)
	if answer.cardinality != ac1.cardinality {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, ac1.cardinality)
	}
	for i, v := range answer.content {
		if ac1.content[i] != v {
			t.Errorf("Got: %d, want: %d", v, i)
		}
	}
}

func TestAndNotBitmap(t *testing.T) {
	ac := newArrayContainerWithCapacity(10)
	bc := newBitmapContainer()

	for i := 0; i < 10; i++ {
		ac.add(uint16(i))
	}

	for i := 0; i < 10; i++ {
		bc.add(uint16(i + 5))
	}

	answer := ac.andNotBitmap(bc)
	if answer.cardinality != 5 {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, 10)
	}
	for i, v := range answer.content {
		if uint16(i) != v {
			t.Errorf("Got: %d, want: %d", v, i)
		}
	}
}

func TestAndNotBitmap_2(t *testing.T) {
	ac := newArrayContainerWithCapacity(10)
	bc := newArrayContainer()

	for i := 0; i < 10; i++ {
		ac.add(uint16(i))
	}

	answer := ac.andNot(bc)
	if answer.cardinality != ac.cardinality {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, ac.cardinality)
	}
	for i, v := range answer.content {
		if ac.content[i] != v {
			t.Errorf("Got: %d, want: %d", v, i)
		}
	}
}

func TestAndNotBitmap_3(t *testing.T) {
	ac := newArrayContainer()
	bc := newArrayContainerWithCapacity(10)

	for i := 0; i < 10; i++ {
		bc.add(uint16(i))
	}

	answer := ac.andNot(bc)
	if answer.cardinality != 0 {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, ac.cardinality)
	}

	if len(answer.content) != 0 {

		t.Errorf("Got: %d, want: %d", answer.content, 0)
	}
}

func TestAndNotBitmap_4(t *testing.T) {
	ac := newArrayContainerWithCapacity(10)
	bc := newArrayContainerWithCapacity(10)

	for i := 1; i <= 10; i++ {
		ac.add(uint16(i))
	}

	for i := 1; i <= 10; i++ {
		bc.add(uint16(i * 20))
	}

	answer := ac.andNot(bc)
	if answer.cardinality != ac.cardinality {
		t.Errorf("Cardinality: %d, want: %d", answer.cardinality, ac.cardinality)
	}
	for i, v := range answer.content {
		if ac.content[i] != v {
			t.Errorf("Got: %d, want: %d", v, i)
		}
	}
}

func TestClear(t *testing.T) {
	ac := newArrayContainerRunOfOnes(0, 9)
	ac.clear()
	if ac.cardinality != 0 {
		t.Errorf("Cardinality: %d, want: 0", ac.cardinality)
	}

	if ac.contains(5) {
		t.Errorf("ArrayContainer is not empty.")
	}
}

func TestConstains(t *testing.T) {
	ac := newArrayContainerRunOfOnes(0, 9)

	for i := 10; i < 20; i++ {
		if ac.contains(uint16(i)) {
			t.Errorf("Array constains %d, want: false)", i)
		}
	}

	for i := 0; i < 9; i++ {
		if !ac.contains(uint16(i)) {
			t.Errorf("Array not constains %d, want: true)", i)
		}
	}

	ac = newArrayContainerRunOfOnes(1, 5)

	if ac.contains(uint16(0)) {
		t.Errorf("Array constains zero.")
	}
}

func TestOrArray(t *testing.T) {
	ac1 := newArrayContainer()
	ac2 := newArrayContainer()
	count := 100

	for i := 0; i < count; i += 2 {
		ac1.add(uint16(i))
		ac2.add(uint16(i + 1))
	}

	result := ac1.orArray(ac2)
	if result.getCardinality() != count {
		t.Errorf("Cardinality: %d, want: %d", result.getCardinality(), count)
	}
	for k, v := range result.(*arrayContainer).content {
		if v != uint16(k) {
			t.Errorf("orArray: %d, want: %d", v, k)
		}
	}
}

func TestOrArray_2(t *testing.T) {
	ac1 := newArrayContainer()
	ac2 := newArrayContainer()
	count := 100

	for i := 0; i < count; i++ {
		ac1.add(uint16(i))
	}

	result := ac1.orArray(ac2)
	if result.getCardinality() != count {
		t.Errorf("Cardinality: %d, want: %d", result.getCardinality(), count)
	}
	for i := 0; i < count; i++ {
		if result.(*arrayContainer).content[i] != uint16(i) {
			t.Errorf("orArray: %d, want: %d", result.(*arrayContainer).content[i], i)
		}
	}
}
