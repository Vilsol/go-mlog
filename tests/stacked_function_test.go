package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestStackedFunction(t *testing.T) {
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
			output: `set @stack 0
jump 9 always
op sub _sampleDynamic_0 @stack 1
read _sampleDynamic_arg2 bank1 _sampleDynamic_0
op sub _sampleDynamic_1 @stack 2
read _sampleDynamic_arg1 bank1 _sampleDynamic_1
op add _sampleDynamic_2 _sampleDynamic_arg1 _sampleDynamic_arg2
set @return _sampleDynamic_2
read @counter bank1 @stack
op add _main_0 1 2
op add @stack @stack 1
write _main_0 bank1 @stack
op rand _main_1 100
op floor _main_2 _main_1
op add @stack @stack 1
write _main_2 bank1 @stack
op add @stack @stack 1
write 19 bank1 @stack
jump 2 always
op sub @stack @stack 3
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
			output: `set @stack 0
jump 4 always
set @return 9
read @counter bank1 @stack
op add @stack @stack 1
write 7 bank1 @stack
jump 2 always
op sub @stack @stack 1
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
			output: `set @stack 0
jump 5 always
set _sampleVariable_x 5
set @return _sampleVariable_x
read @counter bank1 @stack
op add @stack @stack 1
write 8 bank1 @stack
jump 2 always
op sub @stack @stack 1
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
			output: `set @stack 0
jump 5 always
print "hello"
print "\n"
read @counter bank1 @stack
op add @stack @stack 1
write 8 bank1 @stack
jump 2 always
op sub @stack @stack 1`,
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
			output: `set @stack 0
jump 5 always
print "hello"
print "\n"
read @counter bank1 @stack
op add @stack @stack 1
write 8 bank1 @stack
jump 2 always
op sub @stack @stack 1`,
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
			output: `set @stack 0
jump 5 always
print "hello"
print "\n"
read @counter bank1 @stack
op add @stack @stack 1
write 8 bank1 @stack
jump 2 always
op sub @stack @stack 1`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mlog, err := transpiler.GolangToMLOG(test.input, transpiler.Options{
				Stacked: "bank1",
			})

			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
