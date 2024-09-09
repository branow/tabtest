package tab_test

import (
	"errors"
	"fmt"
	"runtime/debug"
	"testing"

	. "github.com/branow/tabtest/tab"
)

func TestRunWithArgs(t *testing.T) {
	ctx := "github.com/branow/tabtest/tab_test.TestRunWithArgs"
	mockT := &testing.T{}
	data := []struct {
		casename string
		t        *testing.T
		args     []Args
		testFunc any
		err      error
	}{
		{
			"nil args", nil, nil, nil,
			errors.Join(
				GetNilErr("t"),
				GetNilErr("args"),
				GetNilErr("f")),
		},
		{
			"is not a function",
			mockT,
			[]Args{{1}, {2}},
			"func (t *testing.T, a int)  {}",
			GetInvalidKindErr("string", "func", "string"),
		},
		{
			"func must have at least one param",
			mockT,
			[]Args{},
			func() {},
			GetInvalidFuncLeastParamNumErr(ctx+".func1", 1),
		},
		{
			"invalid first param *testing.T",
			mockT,
			[]Args{{1}, {2}},
			func(a int, b int) {},
			GetInvalidFuncParamErr(0, ctx+".func2", "*testing.T", "int"),
		},
		{
			"not enough args",
			mockT,
			[]Args{{}, {2}, {4}},
			func(t *testing.T, a int) {},
			WrapTestCaseError(0, GetInvalidFuncParamNumErr(ctx+".func3", 1, 2)),
		},
		{
			"too much args",
			mockT,
			[]Args{{1}, {2, 3}},
			func(t *testing.T, a int) {},
			WrapTestCaseError(1, GetInvalidFuncParamNumErr(ctx+".func4", 3, 2)),
		},
		{
			"cannot convert arg",
			mockT,
			[]Args{{2.3}},
			func(t *testing.T, a string) {},
			WrapTestCaseError(0, WrapArgError(1, GetConvertErr("float64", "string"))),
		},
		{
			"valid running",
			t,
			[]Args{{1, 2}, {-45, 651}, {0, 32}},
			func(t *testing.T, a, b int) {},
			nil,
		},
		{
			"valid running with nil arg",
			t,
			[]Args{{"", nil, nil, nil}},
			func(t *testing.T, str string, err error, s []int, p *int) {},
			nil,
		},
		{
			"valid running with converting",
			t,
			[]Args{{1.1234, 2}, {0, 32.9}},
			func(t *testing.T, a, b int) {},
			nil,
		},
		{
			"valid running with case names",
			t,
			[]Args{{"@first", 1.1234, 2}, {"@second", 0, 32.9}},
			func(t *testing.T, a, b int) {},
			nil,
		},
	}
	for i, d := range data {
		cn := fmt.Sprintf("case %d: %s", i, d.casename)
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
				RunWithArgs(d.t, d.args, d.testFunc)
			}()
		})
	}
}
