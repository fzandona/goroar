package goroar

type container interface {
	add(x uint16) container
	and(x container) container
	or(x container) container
	andNot(x container) container
	// not(x container) container
	xor(x container) container

	// trim()
	// clone() container
	// clear()
	contains(x uint16) bool

	// deserialize() error
	// serialize() error

	getCardinality() int
	// getSizeInBytes() int32
}
