
`tabtest` - multicase tests quicker and easier
================================

Go code (golang) is a package that provides several exported functions to simplify the process of writing multicase/table tests and reduce amount boilerplate code.

Get started:

  * Install testify with [one line of code](#installation), or [update it with another](#staying-up-to-date)

-----
`tab` package
==============

The `tab` package provides three helpful functions that allow you to write multicase tests quicker and easier.

  * [`Run`](#run)
  * [`RunWithCfgs`](#runwithcfgs)
  * [`RunWithArgs`](#runwithargs)

An example of reducing code:

<p align="center">
  <img alt="Using Run function" src="https://raw.githubusercontent.com/branow/tabtest/main/example.png">
</p>

`Run`
------

Function [`Run`](#run) is the most simple function of the package and lets you only remove a couple of lines that usually contain a for-loop. In the following example of a table test using [`Run`](#run), you can see that it is very similar to the usual look of table tests except that the closure function takes an additional parameter that is the number of current iterations, and [`Run`](#run) function requires the number of all iterations.

If you need a bit more advanced behavior take a look at the next functions [`RunWithCfgs`](#runwithcfgs) and [`RunWithArgs`](#runwithargs)

```go
package yours

import (
	"errors"
	"testing"

	"github.com/branow/tabtest/tab"
)

func TestSomething(t *testing.T) {
	data := []struct{
		in string
		exp string
		eErr error
	}{
		{"valid input 1", "valid output 1", nil},
		{"valid input 2", "valid output 2", nil},
		{"invalid input 1", "", errors.New("some error 1")},
		{"invalid input 2", "", errors.New("some error 2")},
	}
	test := func(t *testing.T, i int)  {
		d := data[i]
		act, aErr := DoSomething(d.in)
		if d.exp != act {
			t.Errorf("expected: %s\nactual: %s", d.exp, act)
		}
		if d.eErr != nil && aErr != nil && d.eErr.Error() != aErr.Error() {
			t.Errorf("expected error: %v\nactual error: %v", d.eErr, aErr)
		} else if d.eErr != nil && aErr == nil {
			t.Errorf("expected error: %v", d.eErr)
		} else if d.eErr == nil && aErr != nil {
			t.Errorf("unexpected error: %v", aErr)
		}
	}
	tab.Run(t, len(data), test)
}
```

`RunWithCfgs`
------

In contrast to the [`Run`](#run), [`RunWithCfgs`](#runwithcfgs) provides the closure function with a struct (or any other type) instead of the number of iterations. 

[`RunWithCfgs`](#runwithcfgs) uses Generics that require a declaration of the type (in the next example it's struct Config). It makes code a bit more verbose compared to the previous example but also it lets you separate test data and test logic. 

As an optional feature, here you can specify a case name in an additional field `CaseName`, it's available only for structs and maps.

This function is worth using when more than one test consumes the same data structure.

```go
package yours

import (
	"errors"
	"testing"

	"github.com/branow/tabtest/tab"
)

func TestSomething(t *testing.T) {
	type Config struct{
		CaseName string //optional field let you specify a case name
		in string
		exp string
		eErr error
	}
	data := []Config{
		{"", "valid input 1", "valid output 1", nil},
		{"", "valid input 2", "valid output 2", nil},
		{"", "invalid input 1", "", errors.New("some error 1")},
		{"", "invalid input 2", "", errors.New("some error 2")},
	}
	test := func(t *testing.T, c Config)  {
		act, aErr := DoSomething(c.in)
		if c.exp != act {
			t.Errorf("expected: %s\nactual: %s", c.exp, act)
		}
		if c.eErr != nil && aErr != nil && c.eErr.Error() != aErr.Error() {
			t.Errorf("expected error: %v\nactual error: %v", c.eErr, aErr)
		} else if c.eErr != nil && aErr == nil {
			t.Errorf("expected error: %v", c.eErr)
		} else if c.eErr == nil && aErr != nil {
			t.Errorf("unexpected error: %v", aErr)
		}
	}
	tab.RunWithCfgs(t, data, test)
}
```

An example with data separating.


```go
package yours

import (
	"errors"
	"testing"

	"github.com/branow/tabtest/tab"
)

func TestSomething(t *testing.T) {
	data := ProvideTestSomething()
	tab.RunWithCfgs(t, data, func(t *testing.T, c Config) {
		act, aErr := DoSomething(c.in)
		if c.exp != act {
			t.Errorf("expected: %s\nactual: %s", c.exp, act)
		}
		if c.eErr != nil && aErr != nil && c.eErr.Error() != aErr.Error() {
			t.Errorf("expected error: %v\nactual error: %v", c.eErr, aErr)
		} else if c.eErr != nil && aErr == nil {
			t.Errorf("expected error: %v", c.eErr)
		} else if c.eErr == nil && aErr != nil {
			t.Errorf("unexpected error: %v", aErr)
		}
	})
}

type Config struct {
	CaseName string //optional field let you specify a case name
	in       string
	exp      string
	eErr     error
}

func ProvideTestSomething() []Config {
	return []Config{
		{"do 1", "valid input 1", "valid output 1", nil},
		{"do 2", "valid input 2", "valid output 2", nil},
		{"nil", "invalid input 1", "", errors.New("some error")},
		{"validation", "invalid input 2", "", errors.New("some error")},
	}
}
```

`RunWithArgs`
------

[`RunWithArgs`](#runwithargs) lets you specify all required variables for the test logic in the parameters of the closure function. The code automatically matches the data with closure parameters based on the order of the parameters and the values in the slice of `tab.Args` using reflection. The matching starts from the second parameter because the first must always be a `*testing.T`.

As an optional feature, you can add an additional value in the arguments slice to specify a case name. The function considers any string value that is at index zero of slice arguments and starts with the prefix `@` as a case name.

Although using [`RunWithArgs`](#runwithargs) looks the most concise, it could become very confusing if test logic requires a lot of arguments. In such situations, it is better to use [`Run`](#run) or [`RunWithCfgs`](#runwithcfgs).

```go
package yours

import (
	"errors"
	"testing"

	"github.com/branow/tabtest/tab"
)

func TestSomething(t *testing.T) {
	data := []tab.Args{
		{"@do", "valid input 1", "valid output 1", nil},
		{"valid input 2", "valid output 2", nil},
		{"@nil", "invalid input 1", "", errors.New("some error")},
		{"@validation", "invalid input 2", "", errors.New("some error")},
	}
	test := func(t *testing.T, in, exp string, eErr error) {
		act, aErr := DoSomething(in)
		if exp != act {
			t.Errorf("expected: %s\nactual: %s", exp, act)
		}
		if eErr != nil && aErr != nil && eErr.Error() != aErr.Error() {
			t.Errorf("expected error: %v\nactual error: %v", eErr, aErr)
		} else if eErr != nil && aErr == nil {
			t.Errorf("expected error: %v", eErr)
		} else if eErr == nil && aErr != nil {
			t.Errorf("unexpected error: %v", aErr)
		}
	}
	tab.RunWithArgs(t, data, test)
}
```


----

Installation
============

To install Testify, use `go get`:

    go get github.com/branow/tabtest

This will then make the following packages available to you:

    github.com/branow/tabtest/tab


------

Staying up to date
==================

To update Testify to the latest version, use `go get -u github.com/branow/tabtest`.

------

Supported go versions
==================

We currently support the most recent major Go versions from 1.23 onward.

------

Contributing
============

Please feel free to submit issues, fork the repository and send pull requests!

------

License
=======

This project is licensed under the terms of the [MIT license](https://github.com/branow/tabtest/blob/main/LICENSE.txt).