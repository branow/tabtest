package examples

import (
	"testing"

	"github.com/branow/tabtest/tab"
)

func TestSomething_RunWithCfgsSep(t *testing.T) {
	test := func(t *testing.T, c Config) {
		act, aErr := DoMath(c.num1, c.num2, c.op)
		if c.exp != act {
			t.Errorf("expected: %v\nactual: %v", c.exp, act)
		}
		if c.errMsg != "" && aErr != nil && c.errMsg != aErr.Error() {
			t.Errorf("expected error: %v\nactual error: %v", c.errMsg, aErr)
		} else if c.errMsg != "" && aErr == nil {
			t.Errorf("expected error: %v", c.errMsg)
		} else if c.errMsg == "" && aErr != nil {
			t.Errorf("unexpected error: %v", aErr)
		}
	}
	tab.RunWithCfgs(t, ProvideTest(), test)
}

type Config struct {
	CaseName string //optional field that let you specify a case name
	num1     int
	num2     int
	op       string
	exp      int
	errMsg   string
}

func ProvideTest() []Config {
	return []Config{
		{"addition", 2, 2, "+", 4, ""},
		{"subtraction", 2, 2, "-", 0, ""},
		{"multiplication", 2, 2, "*", 4, ""},
		{"division", 2, 2, "/", 1, ""},
		{"bad_division", 2, 0, "/", 0, "division by zero"},
		{"bad_op", 2, 2, "?", 0, "unknown operator ?"},
	}
}
