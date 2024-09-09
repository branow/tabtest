package tab_test

import (
	"errors"
	"fmt"
	"runtime/debug"
	"testing"

	. "github.com/branow/tabtest/tab"
)

func TestRunWithCfgs(t *testing.T) {
	type s struct {
		CaseName string
	}
	type str string
	data := []struct {
		casename string
		t        *testing.T
		cfgs     []any
		testFunc func(*testing.T, any)
		err      error
	}{
		{
			"nil args", nil, nil, nil,
			errors.Join(
				GetNilErr("t"),
				GetNilErr("cfgs"),
				GetNilErr("f")),
		},
		{
			"valid running",
			t,
			[]any{2, 3, 4, 5},
			func(t *testing.T, a any) {},
			nil,
		},
		{
			"valid running with converting",
			t,
			[]any{1.1234, 0, -32.9},
			func(t *testing.T, a any) {},
			nil,
		},
		{
			"valid running with map",
			t,
			[]any{map[int]string{1: "cfg"}},
			func(t *testing.T, a any) {},
			nil,
		},
		{
			"valid running with case name struct",
			t,
			[]any{s{"MYCASE I DO"}, s{"YOUCASE YOU DO"}},
			func(t *testing.T, a any) {},
			nil,
		},
		{
			"valid running with case name (map)",
			t,
			[]any{map[string]string{"CaseName": "MYCASE"}},
			func(t *testing.T, a any) {},
			nil,
		},
		{
			"valid running with case name (map)",
			t,
			[]any{map[str]string{"CaseName": "MYCASE"}},
			func(t *testing.T, a any) {},
			nil,
		},
		{
			"valid running with wrong value type case name (map)",
			t,
			[]any{map[string]int{"one": 1}},
			func(t *testing.T, a any) {},
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
				RunWithCfgs[any](d.t, d.cfgs, d.testFunc)
			}()
		})
	}
}
