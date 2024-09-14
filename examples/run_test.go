package examples

import (
	"testing"

	"github.com/branow/tabtest/tab"
)

func TestSomething_Run(t *testing.T) {
	data := []struct {
		num1   int
		num2   int
		op     string
		exp    int
		errMsg string
	}{
		{2, 2, "+", 4, ""},
		{2, 2, "-", 0, ""},
		{2, 2, "*", 4, ""},
		{2, 2, "/", 1, ""},
		{2, 0, "/", 0, "division by zero"},
		{2, 2, "?", 0, "unknown operator ?"},
	}
	test := func(t *testing.T, i int) {
		d := data[i]
		act, aErr := DoMath(d.num1, d.num2, d.op)
		if d.exp != act {
			t.Errorf("expected: %v\nactual: %v", d.exp, act)
		}
		if d.errMsg != "" && aErr != nil && d.errMsg != aErr.Error() {
			t.Errorf("expected error: %v\nactual error: %v", d.errMsg, aErr)
		} else if d.errMsg != "" && aErr == nil {
			t.Errorf("expected error: %v", d.errMsg)
		} else if d.errMsg == "" && aErr != nil {
			t.Errorf("unexpected error: %v", aErr)
		}
	}
	tab.Run(t, len(data), test)
}
