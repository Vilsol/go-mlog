package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestStacklessFunction(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name: "FunctionDynamicReturn",
			input: `package main

import (
	"github.com/Vilsol/go-mlog/m"
	"github.com/Vilsol/go-mlog/x"
)

func main() {
	print(sampleDynamic(1 + 2, m.Floor(m.Random(100))))
}

func sampleDynamic(arg1 int, arg2 int) int {
	return arg1 + arg2
}`,
			output: `jump 6 always
set _sampleDynamic_arg2 @funcArg_sampleDynamic_1
set _sampleDynamic_arg1 @funcArg_sampleDynamic_0
op add _sampleDynamic_0 _sampleDynamic_arg1 _sampleDynamic_arg2
set @return _sampleDynamic_0
set @counter @funcTramp_sampleDynamic
op add _main_0 1 2
set @funcArg_sampleDynamic_0 _main_0
op rand _main_1 100
op floor _main_2 _main_1
set @funcArg_sampleDynamic_1 _main_2
set @funcTramp_sampleDynamic 13
jump 1 always
set _main_3 @return
print _main_3`,
		},
		{
			name: "FunctionStatic",
			input: `package main

import (
	"github.com/Vilsol/go-mlog/m"
	"github.com/Vilsol/go-mlog/x"
)

func main() {
	print(sampleStatic())
}

func sampleStatic() int {
	return 9
}`,
			output: `jump 3 always
set @return 9
set @counter @funcTramp_sampleStatic
set @funcTramp_sampleStatic 5
jump 1 always
set _main_0 @return
print _main_0`,
		},
		{
			name: "FunctionVariable",
			input: `package main

import (
	"github.com/Vilsol/go-mlog/m"
	"github.com/Vilsol/go-mlog/x"
)

func main() {
	print(sampleVariable())
}

func sampleVariable() int {
	x := 5
	return x
}`,
			output: `jump 4 always
set _sampleVariable_x 5
set @return _sampleVariable_x
set @counter @funcTramp_sampleVariable
set @funcTramp_sampleVariable 6
jump 1 always
set _main_0 @return
print _main_0`,
		},
		{
			name: "FunctionNone",
			input: `package main

import (
	"github.com/Vilsol/go-mlog/m"
	"github.com/Vilsol/go-mlog/x"
)

func main() {
	sampleNone()
}

func sampleNone() {
	println("hello")
}`,
			output: `jump 4 always
print "hello"
print "\n"
set @counter @funcTramp_sampleNone
set @funcTramp_sampleNone 6
jump 1 always`,
		},
		{
			name: "TreeShake",
			input: `package main

func main() {
	hello()
}

func hello() {
	println("hello")
}

func foo() {
	println("foo")
}

func bar() {
	println("bar")
}`,
			output: `jump 4 always
print "hello"
print "\n"
set @counter @funcTramp_hello
set @funcTramp_hello 6
jump 1 always`,
		},
		{
			name: "IgnoreEmpty",
			input: `package main

func main() {
	hello()
}

func hello() {
	println("hello")
}

func foo() {
}`,
			output: `jump 4 always
print "hello"
print "\n"
set @counter @funcTramp_hello
set @funcTramp_hello 6
jump 1 always`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mlog, err := transpiler.GolangToMLOG(test.input, transpiler.Options{})

			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
