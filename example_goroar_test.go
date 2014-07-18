package goroar_test

import (
	"fmt"

	"github.com/fzandona/goroar"
)

// ExampleGoroar demonstrates how to use the goroar library.
func Example_goroar() {
	rb1 := goroar.BitmapOf(1, 2, 3, 4, 5)
	rb2 := goroar.BitmapOf(2, 3, 4)
	rb3 := goroar.New()

	fmt.Println("Cardinality: ", rb1.Cardinality())

	fmt.Println("Contains 3? ", rb1.Contains(3))

	rb1.And(rb2)

	rb3.Add(1)
	rb3.Add(5)

	rb3.Or(rb1)

	// prints 1, 2, 3, 4, 5
	for value := range rb3.Iterator() {
		fmt.Println(value)
	}
}
