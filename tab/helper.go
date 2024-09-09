package tab

import (
	"reflect"
	"runtime"
	"slices"
)

func convert(v reflect.Value, t reflect.Type) (reflect.Value, error) {
	if !v.IsValid() {
		return reflect.New(t).Elem(), nil
	}
	if v.Type() != t {
		if !v.CanConvert(t) {
			return reflect.Value{}, GetConvertErr(v.Type(), t)
		}
		v = v.Convert(t)
	}
	return v, nil
}

func getFuncName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func checkNil(a any, varName string) error {
	if a == nil || (reflect.ValueOf(a).Kind() == reflect.Ptr && reflect.ValueOf(a).IsNil()) {
		return GetNilErr(varName)
	}
	av := reflect.ValueOf(a)
	kinds := []reflect.Kind{reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice}
	if slices.Contains(kinds, av.Kind()) && av.IsNil() {
		return GetNilErr(varName)
	}
	return nil
}
