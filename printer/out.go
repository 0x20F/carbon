package printer

import "fmt"

var out OutputWrapper = &impl{}

type OutputWrapper interface {
	Ln(a ...interface{})
}

type impl struct{}

// Wraps fmt.Println() so we can mock it away or replace
// it easily when needed.
// Does absolutely nothing else, just makes sure that
// everything will be called correctly.
func (i *impl) Ln(a ...interface{}) {
	fmt.Println(a...)
}

// Replaces the default Output instance with a custom
// implementation.
//
// Note that this exists for the sole purpose of unit testing.
// It makes it easy to replace the actual printer with something
// that doesn't output anything during tests.
func WrapStdout(custom OutputWrapper) {
	out = custom
}
