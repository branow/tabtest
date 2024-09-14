package examples

import (
	"testing"

	"github.com/branow/tabtest/tab"
)

func TestSomething_RunWithArgs(t *testing.T) {
	data := []tab.Args{
		{"@addition", 2, 2, "+", 4, ""},
		{"@subtraction", 2, 2, "-", 0, ""},
		{"@multiplication", 2, 2, "*", 4, ""},
		{"@division", 2, 2, "/", 1, ""},
		{"@bad_division", 2, 0, "/", 0, "division by zero"},
		{"@bad_op", 2, 2, "?", 0, "unknown operator ?"},
	}
	test := func(t *testing.T, num1, num2 int, op string, exp int, errMsg string) {
		act, aErr := DoMath(num1, num2, op)
		if exp != act {
			t.Errorf("expected: %v\nactual: %v", exp, act)
		}
		if errMsg != "" && aErr != nil && errMsg != aErr.Error() {
			t.Errorf("expected error: %v\nactual error: %v", errMsg, aErr)
		} else if errMsg != "" && aErr == nil {
			t.Errorf("expected error: %v", errMsg)
		} else if errMsg == "" && aErr != nil {
			t.Errorf("unexpected error: %v", aErr)
		}
	}
	tab.RunWithArgs(t, data, test)
}
