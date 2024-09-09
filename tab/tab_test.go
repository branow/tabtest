package tab_test

import (
	"errors"
	"fmt"
	"runtime/debug"
	"testing"

	. "github.com/branow/tabtest/tab"
)

func TestRun(t *testing.T) {
	data := []struct {
		casename string
		t        *testing.T
		num      int
		testFunc func(*testing.T, int)
		err      error
	}{
		{
			"nil args", nil, 1, nil,
			errors.Join(
				GetNilErr("t"),
				GetNilErr("f")),
		},
		{
			"invalid number of cycles",
			t,
			-5,
			func(t *testing.T, i int) {},
			GetInvalidNumCycles(-5),
		},
		{
			"invalid number of cycles",
			t,
			3,
			func(t *testing.T, i int) {},
			nil,
		},
	}
	for i, d := range data {
		cn := fmt.Sprintf("MAIN CASE %d: %s", i, d.casename)
		t.Run(cn, func(t *testing.T) {
			func() {
				defer func() {
					err := recover()
					if d.err == nil && err != nil {
						t.Errorf("unexpected panic: %s\nstack trace:\n%s", err, string(debug.Stack()))
					} else if d.err != nil && err == nil {
						t.Error("expected panic: ", d.err.Error())
					} else if d.err != nil && err != nil {
						m1, m2 := d.err.Error(), err.(error).Error()
						if m1 != m2 {
							t.Errorf("invalid panic:\nexpected:\n'%s'\nactual:\n'%s'\nstack trace:\n%s", m1, m2, string(debug.Stack()))
						}
					}
				}()
				Run(d.t, d.num, d.testFunc)
			}()
		})
	}
}
