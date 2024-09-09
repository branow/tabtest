package tab

import "fmt"


func WrapArgError(i int, err error) error {
	err = fmt.Errorf("arg %d: %w", i, err)
	return wrapTabError(err)
}

func WrapTestCaseError(i int, err error) error {
	err = fmt.Errorf("case %d: %w", i, err)
	return wrapTabError(err)
}

func GetInvalidNumCycles(num int) error {
	err := fmt.Errorf("number of cycles cannot be negative: %d", num)
	return wrapTabError(err)
}

func GetConvertErr(from, to any) error {
	err := fmt.Errorf("cannot convert from %v to %v", from, to)
	return wrapTabError(err)
}

func GetInvalidFuncParamErr(paramNum int, funcName string, exp, act any) error {
	err := fmt.Errorf("%d(th) param of func %s must be %v, but it is %v", paramNum, funcName, exp, act)
	return wrapTabError(err)
}

func GetInvalidFuncParamNumErr(funcName string, expNum, actNum int) error {
	err := fmt.Errorf("func %s must have %d params, but it has %d", funcName, expNum, actNum)
	return wrapTabError(err)
}

func GetInvalidFuncLeastParamNumErr(funcName string, num int) error {
	err := fmt.Errorf("func %s must have at least %d param(s)", funcName, num)
	return wrapTabError(err)
}

func GetInvalidKindErr(typeName, expKind, actKind any) error {
	err := fmt.Errorf("%v must be a %v, but it is a %v", typeName, expKind, actKind)
	return wrapTabError(err)
}

func GetNilErr(varName string) error {
	err := fmt.Errorf("%s is nil", varName)
	return wrapTabError(err)
}

func wrapTabError(err error) error {
	return fmt.Errorf("tab: %w", err)
}
