package goroar_test

import (
	"fmt"

	"github.com/fzandona/goroar"
)

// ExampleLarge demonstrates how to use the goroar library.
func Example_large() {
	rb1 := goroar.New()
	rb2 := goroar.New()

	for i := 0; i < 1000000; i += 2 {
		rb1.Add(uint32(i))
		rb2.Add(uint32(i + 1))
	}
	fmt.Println(rb1.Cardinality(), rb2.Cardinality())

	rb1.Or(rb2)

	rb1.Stats()
}
