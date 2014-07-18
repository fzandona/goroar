goroar
======

*goroar* is an implementation of [Roaring Bitmaps](http://roaringbitmap.org) in Golang. Roaring bitmaps is a new form of compressed bitmaps, proposed by Daniel Lemire *et. al.*, which often offers better compression and fast access than other compressed bitmap approaches.

Make sure to check Lemire's paper for a detailed explanation and comparison with WAH and Concise.

Usage
-----
Get the library using `go get`:

    go get github.com/fzandona/goroar

### Quickstart

```go
package main

import (
    "fmt"

    "github.com/fzandona/goroar"
)

func main() {
    rb1 := goroar.BitmapOf(1, 2, 3, 4, 5)
    rb2 := goroar.BitmapOf(2, 3, 4)
    rb3 := goroar.BitmapOf(1)

    fmt.Println("Cardinality: ", rb1.Cardinality())

    fmt.Println("Contains 3? ", rb1.Contains(3))

    rb1.And(rb2)

    rb3.Add(5)

    rb3.Or(rb1)

    // prints 1, 2, 3, 4, 5
    for value := range rb3.Iterator() {
        fmt.Println(value)
    }
}
```

### Documentation

Documentation of the latest code in master is available at [godoc](http://godoc.org/github.com/fzandona/goroar).

TODO
----

* Immutable bitwise operations
* Serialization
* More idiomatic Go
* Test re-factoring & coverage
