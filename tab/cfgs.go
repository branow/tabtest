package tab

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type TestFuncWithCfg[C any] func(*testing.T, C)

// RunWithArgs runs the given function (f any) provided with arguments
//
// RunWithCfgs runs the given function (f [TestFuncWithCfg]) provided
// with configs (cfgs []C) the number of times that equals the number
// of the given configs. If any of the given values is nil, it causes panic.
//
// Config is a variable of any type that transfers to a function call
// as a second argument. Optionally, if the config is a struct that contains
// a field named [CaseNameField] of type string or any other type
// with the kind string that can be converted to a string, this field
// value is used as a part of the name of the current test case (as the
// first argument of [testing.T.Run()]). It also works with maps that contain
// a pair with key [CaseNameField] and value of type string or convertible one.
//
// The examples of test case names: "case 0", "case 0: invalid input"
// Here are 'case 0' a default name and 'invalid input' a custom name taken from args.
func RunWithCfgs[C any](t *testing.T, cfgs []C, f TestFuncWithCfg[C]) {
	// nil check
	err := errors.Join(checkNil(t, "t"), checkNil(cfgs, "cfgs"), checkNil(f, "f"))
	if err != nil {
		panic(err)
	}

	// go through all the configs and run the test functions
	for i, cfg := range cfgs {

		//try to get CaseName from cfg
		cn := fmt.Sprintf("case %d", i)
		ccn, ok := getCaseName(reflect.ValueOf(cfg))
		if ok {
			cn = fmt.Sprintf("%s: %s", cn, ccn)
		}

		t.Run(cn, func(t *testing.T) {
			f(t, cfg)
		})
	}
}

const CaseNameField = "CaseName"

func getCaseName(c reflect.Value) (string, bool) {
	ct := c.Type()
	ck := ct.Kind()
	if ck != reflect.Struct && ck != reflect.Map {
		return "", false
	}
	var val reflect.Value
	if ck == reflect.Struct {
		val = c.FieldByName(CaseNameField)
	} else if ck == reflect.Map {
		if ct.Key().Kind() != reflect.String {
			return "", false
		}
		key := reflect.ValueOf(CaseNameField)
		key = key.Convert(ct.Key())
		val = c.MapIndex(key)
	}
	if val.IsValid() && val.Kind() == reflect.String {
		return val.String(), true
	}
	return "", false
}
