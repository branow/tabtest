package tab

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const CaseNamePrefix = "@"

type Args []any

// RunWithArgs runs the given function (f any) provided with arguments
// (args [][Args]) the number of times that equals the length of the argument
// slice.
//
// If any of the given values is nil, it causes panic. If f is not a function,
// panic is caused as well. The first parameter of the f function must
// be of type *[testing.T], and types of all other parameters must match
// the values of the given argument or be able to be converted at every test case.
//
// Optionally, if the first element of arguments for a specific call is a string with the 
// prefix [CaseNamePrefix], it is considered as a case name and isn't transferred 
// to the function. A test case name is a value that is transferred to the  
// [testing.T.Run()]) as a first parameter.  The custom case name is appended to 
// the default case name "case <num>"  and then it is used.
//
// The examples of test case names: "case 0", "case 0: invalid input"
// Here are 'case 0' a default name and 'invalid input' a custom name taken from args.
func RunWithArgs(t *testing.T, args []Args, f any) {
	// nil check
	err := errors.Join(checkNil(t, "t"), checkNil(args, "args"), checkNil(f, "f"))
	if err != nil {
		panic(err)
	}

	// check if the given f any is a function
	ft := reflect.TypeOf(f)
	fv := reflect.ValueOf(f)
	if ft.Kind() != reflect.Func {
		panic(GetInvalidKindErr(ft, reflect.Func, ft.Kind()))
	}

	// check and parse args
	vals, cns := parseArgs(f, t, args)

	// run through all the args and call the test func
	for i, val := range vals {
		//try to get CaseName from cfg
		cn := fmt.Sprintf("case %d", i)
		if cns[i] != "" {
			cn = fmt.Sprintf("%s: %s", cn, cns[i])
		}

		t.Run(cn, func(t *testing.T) {
			fv.Call(val)
		})
	}
}

func parseArgs(f any, t *testing.T, args []Args) ([][]reflect.Value, []string) {
	funcName := getFuncName(f)
	ft := reflect.TypeOf(f)

	// check the number of the params of the func
	num := ft.NumIn()
	if num < 1 {
		panic(GetInvalidFuncLeastParamNumErr(funcName, 1))
	}

	// check the first param of the func
	pt1 := ft.In(0)
	tt := reflect.TypeOf(t)
	if pt1.Kind() != reflect.Pointer ||
		pt1.Elem() != tt.Elem() {
		panic(GetInvalidFuncParamErr(0, funcName, tt, pt1))
	}

	// go through all the args and check if they match the func params
	vals := make([][]reflect.Value, len(args))
	cns := make([]string, len(args))
	errs := []error{}
	for i, arg := range args {
		// check if args contain case name
		var hasCn bool
		if len(arg) != 0 {
			cn, ok := parseCaseName(arg[0])
			hasCn = ok
			if hasCn {
				cns[i] = cn
			}
		}

		// check if num of args match num of params
		argNum := len(arg) + 1 // plus + 1 because of t *testing.T
		if hasCn {
			argNum-- //because one of given args is case name
		}
		if argNum != num {
			err := GetInvalidFuncParamNumErr(funcName, argNum, num)
			err = WrapTestCaseError(i, err)
			errs = append(errs, err)
			continue
		}

		// create slice of args for every call
		val := make([]reflect.Value, 0, num)
		val = append(val, reflect.ValueOf(t))
		if hasCn {
			arg = arg[1:]
		}
		for j, a := range arg {
			av := reflect.ValueOf(a)
			pt := ft.In(j + 1)
			av, err := convert(av, pt)
			if err != nil {
				err := WrapTestCaseError(i, WrapArgError(j+1, err))
				errs = append(errs, err)
				continue
			}
			val = append(val, av)
		}
		vals[i] = val
	}

	// panic if there is at least one mismatch
	err := errors.Join(errs...)
	if err != nil {
		panic(err)
	}

	return vals, cns
}

func parseCaseName(a any) (string, bool) {
	at, av := reflect.TypeOf(a), reflect.ValueOf(a)
	if at != reflect.TypeOf("") {
		return "", false
	}
	val := av.String()
	if !strings.HasPrefix(val, CaseNamePrefix) {
		return "", false
	}
	return val[len(CaseNamePrefix):], true
}
