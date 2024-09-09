// Package tab implements functions for easier creating multicase/table tests.
package tab

import (
	"errors"
	"fmt"
	"testing"
)

type TestFunc func(*testing.T, int)

// Run runs the given function f TestFunc the given number of times
// (num int) and transfer in every call the number of the current iteration.
// If the given t *testing.T or f TestFunc is nil, panic will be caused.
// If the number (num int) of iteration is negative, it causes panic too.
func Run(t *testing.T, num int, f TestFunc) {
	err := errors.Join(checkNil(t, "t"), checkNil(f, "f"))
	if err != nil {
		panic(err)
	}

	if num < 0 {
		panic(GetInvalidNumCycles(num))
	}

	for i := 0; i < num; i++ {
		cn := fmt.Sprintf("case %d", i)
		t.Run(cn, func(t *testing.T) {
			f(t, i)
		})
	}
}
