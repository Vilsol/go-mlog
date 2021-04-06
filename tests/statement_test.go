package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestStatement(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name: "IfElseifElse",
			input: TestMain(`if x := 1; x == 2 {
	print(3)
} else if x == 4 {
	print(5)
} else {
	print(6)
}`),
			output: `set _main_x 1
op equal _main_0 _main_x 2
jump 4 equal _main_0 1
jump 6 always
print 3
jump 12 always
op equal _main_1 _main_x 4
jump 9 equal _main_1 1
jump 11 always
print 5
jump 12 always
print 6`,
		},
		{
			name:  "ForLoop",
			input: TestMain(`for i := 0; i < 10; i++ { print(i) }`),
			output: `set _main_i 0
jump 3 lessThan _main_i 10
jump 6 always
print _main_i
op add _main_i _main_i 1
jump 3 lessThan _main_i 10`,
		},
		{
			name:   "Reassignment",
			input:  TestMain(`x := y`),
			output: `set _main_x _main_y`,
		},
		{
			name:   "VariableBooleans",
			input:  TestMain(`x := false`),
			output: `set _main_x false`,
		},
		{
			name:   "VariableCharacter",
			input:  TestMain(`x := 'A'`),
			output: `set _main_x "A"`,
		},
		{
			name:  "Break",
			input: TestMain(`for i := 0; i < 10; i++ { if i == 5 { break; }; println(i); }`),
			output: `set _main_i 0
jump 3 lessThan _main_i 10
jump 11 always
op equal _main_0 _main_i 5
jump 6 equal _main_0 1
jump 7 always
jump 11 always
print _main_i
print "\n"
op add _main_i _main_i 1
jump 3 lessThan _main_i 10`,
		},
		{
			name:  "Continue",
			input: TestMain(`for i := 0; i < 10; i++ { if i == 5 { continue; }; println(i); }`),
			output: `set _main_i 0
jump 3 lessThan _main_i 10
jump 11 always
op equal _main_0 _main_i 5
jump 6 equal _main_0 1
jump 7 always
jump 9 always
print _main_i
print "\n"
op add _main_i _main_i 1
jump 3 lessThan _main_i 10`,
		},
		{
			name: "Switch",
			input: TestMain(`switch x := 3 + 7; x {
case 0:
	println("0")
case 1:
	println("1")
	fallthrough
case 2:
	println("2")
	fallthrough
case 3, 4:
	println("3, 4")
	break
case 5, 6:
	println("5, 6")
	break
default:
	println("default")
	break
}`),
			output: `op add _main_x 3 7
jump 9 equal _main_x 0
jump 12 equal _main_x 1
jump 14 equal _main_x 2
jump 16 equal _main_x 3
jump 16 equal _main_x 4
jump 19 equal _main_x 5
jump 19 equal _main_x 6
jump 25 always
print "0"
print "\n"
jump 25 always
print "1"
print "\n"
print "2"
print "\n"
print "3, 4"
print "\n"
jump 25 always
print "5, 6"
print "\n"
jump 25 always
print "default"
print "\n"
jump 25 always`,
		},
		{
			name:   "IgnoredVariable",
			input:  TestMain(`_ := false`),
			output: `set @_ false`,
		},
		{
			name:   "OperatorAssign",
			input:  TestMain(`x += 1`),
			output: `op add _main_x _main_x 1`,
		},
		{
			name: "SelectorAssignment",
			input: TestMain(`a := m.RTAny
b := m.RSDistance
c := m.This
d := m.BCore`),
			output: `set _main_a any
set _main_b distance
set _main_c @this
set _main_d core`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mlog, err := transpiler.GolangToMLOG(test.input, transpiler.Options{
				NoStartup: true,
			})

			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
