package transpiler

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/transpiler"
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
set @return_0 _sampleDynamic_2
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
set _main_3 @return_0
print _main_3`,
		},
		{
			name: "FunctionStatic",
			input: `package main

func main() {
	print(sampleStatic())
}

func sampleStatic() int {
	return 9
}`,
			output: `set @stack 0
jump 4 always
set @return_0 9
read @counter bank1 @stack
op add @stack @stack 1
write 7 bank1 @stack
jump 2 always
op sub @stack @stack 1
set _main_0 @return_0
print _main_0`,
		},
		{
			name: "FunctionVariable",
			input: `package main

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
set @return_0 _sampleVariable_x
read @counter bank1 @stack
op add @stack @stack 1
write 8 bank1 @stack
jump 2 always
op sub @stack @stack 1
set _main_0 @return_0
print _main_0`,
		},
		{
			name: "FunctionNone",
			input: `package main

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
		{
			name: "MultipleReturnValues",
			input: `package main

func main() {
	x, y, z := Hello()
	print(x, y, z)
	print(Hello())
	World(Hello())
}

func Hello() (int, int, int) {
	return 1, 2, 3
}

func World(x int, y int, z int) {
	print(x, y, z)
}`,
			output: `set @stack 0
jump 16 always
set @return_0 1
set @return_1 2
set @return_2 3
read @counter bank1 @stack
op sub _World_0 @stack 1
read _World_z bank1 _World_0
op sub _World_1 @stack 2
read _World_y bank1 _World_1
op sub _World_2 @stack 3
read _World_x bank1 _World_2
print _World_x
print _World_y
print _World_z
read @counter bank1 @stack
op add @stack @stack 1
write 19 bank1 @stack
jump 2 always
op sub @stack @stack 1
set _main_x @return_0
set _main_y @return_1
set _main_z @return_2
print _main_x
print _main_y
print _main_z
op add @stack @stack 1
write 29 bank1 @stack
jump 2 always
op sub @stack @stack 1
set _main_0 @return_0
set _main_1 @return_1
set _main_2 @return_2
print _main_0
print _main_1
print _main_2
op add @stack @stack 1
write 39 bank1 @stack
jump 2 always
op sub @stack @stack 1
set _main_3 @return_0
set _main_4 @return_1
set _main_5 @return_2
op add @stack @stack 1
write _main_3 bank1 @stack
op add @stack @stack 1
write _main_4 bank1 @stack
op add @stack @stack 1
write _main_5 bank1 @stack
op add @stack @stack 1
write 52 bank1 @stack
jump 6 always
op sub @stack @stack 4`,
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

			test.output = test.output + "\nend"
			testza.AssertEqual(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
