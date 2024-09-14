# `tabtest` - multicase tests quicker and easier

[![Build Status](https://github.com/branow/tabtest/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/branow/tabtest/actions/workflows/go.yml) 
[![Go Report Card](https://goreportcard.com/badge/github.com/branow/tabtest)](https://goreportcard.com/report/github.com/branow/tabtest) 
[![PkgGoDev](https://pkg.go.dev/badge/github.com/branow/tabtest)](https://pkg.go.dev/github.com/branow/tabtest)

Go code (golang) is a package that provides several exported functions to simplify the process of writing multicase/table tests and reduce amount boilerplate code.

## Installation

To install `tabtest`, use `go get`:

    go get github.com/branow/tabtest

This will then make the following packages available to you:

    github.com/branow/tabtest/tab

## Staying up to date

To update `tabtest` to the latest version, use `go get -u github.com/branow/tabtest`.


## Supported go versions

We currently support the most recent major Go versions from `1.23` onward.

## `tab` package

The `tab` package provides three helpful functions that allow you to write multicase tests quicker and easier.

  * [`Run`](#run)
  * [`RunWithCfgs`](#runwithcfgs)
  * [`RunWithArgs`](#runwithargs)

An example of reducing code:

<p align="center">
  <img alt="Using Run function" src="https://raw.githubusercontent.com/branow/tabtest/main/example.png">
</p>

### `Run`

Function [`Run`](#run) is the most simple function of the package and lets you only remove a couple of lines that usually contain a for-loop. In the following example of a table test using [`Run`](#run), you can see that it is very similar to the usual look of table tests except that the closure function takes an additional parameter that is the number of current iterations, and [`Run`](#run) function requires the number of all iterations.

If you need a bit more advanced behavior take a look at the next functions [`RunWithCfgs`](#runwithcfgs) and [`RunWithArgs`](#runwithargs)

```go
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

```
[The example file.](https://github.com/branow/tabtest/blob/main/examples/run_test.go)

### `RunWithCfgs`

In contrast to the [`Run`](#run), [`RunWithCfgs`](#runwithcfgs) provides the closure function with a struct (or any other type) instead of the number of iterations. 

[`RunWithCfgs`](#runwithcfgs) uses Generics that require a declaration of the type (in the next example it's struct Config). It makes code a bit more verbose compared to the previous example but also it lets you separate test data and test logic. 

As an optional feature, here you can specify a case name in an additional field `CaseName`, it's available only for structs and maps.

This function is worth using when more than one test consumes the same data structure.

```go
package examples

import (
	"testing"

	"github.com/branow/tabtest/tab"
)

func TestSomething_RunWithCfgs(t *testing.T) {
	type Config struct {
		CaseName string //optional field that let you specify a case name
		num1     int
		num2     int
		op       string
		exp      int
		errMsg   string
	}
	data := []Config{
		{"addition", 2, 2, "+", 4, ""},
		{"subtraction", 2, 2, "-", 0, ""},
		{"multiplication", 2, 2, "*", 4, ""},
		{"division", 2, 2, "/", 1, ""},
		{"bad_division", 2, 0, "/", 0, "division by zero"},
		{"bad_op", 2, 2, "?", 0, "unknown operator ?"},
	}
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
	tab.RunWithCfgs(t, data, test)
}

```
[The example file.](https://github.com/branow/tabtest/blob/main/examples/run_with_cfgs_test.go)

An example with data separating.


```go
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

```
[The example file.](https://github.com/branow/tabtest/blob/main/examples/run_with_cfgs_sep_test.go)

### `RunWithArgs`

[`RunWithArgs`](#runwithargs) lets you specify all required variables for the test logic in the parameters of the closure function. The code automatically matches the data with closure parameters based on the order of the parameters and the values in the slice of `tab.Args` using reflection. The matching starts from the second parameter because the first must always be a `*testing.T`.

As an optional feature, you can add an additional value in the arguments slice to specify a case name. The function considers any string value that is at index zero of slice arguments and starts with the prefix `@` as a case name.

Although using [`RunWithArgs`](#runwithargs) looks the most concise, it could become very confusing if test logic requires a lot of arguments. In such situations, it is better to use [`Run`](#run) or [`RunWithCfgs`](#runwithcfgs).

```go
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

```
[The example file.](https://github.com/branow/tabtest/blob/main/examples/run_with_args_test.go)


## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the terms of the [MIT license](https://github.com/branow/tabtest/blob/main/LICENSE.txt).