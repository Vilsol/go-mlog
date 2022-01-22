package transpiler

import (
	"testing"
)

func TestStatement(t *testing.T) {
	tests := []Test{
		{
			name: "IfElseifElse",
			input: TestMain(`if x := 1; x == 2 {
	print(3)
} else if x == 4 {
	print(5)
} else {
	print(6)
}`, false, false),
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
			input: TestMain(`for i := 0; i < 10; i++ { print(i) }`, false, false),
			output: `set _main_i 0
jump 3 lessThan _main_i 10
jump 6 always
print _main_i
op add _main_i _main_i 1
jump 3 lessThan _main_i 10`,
		},
		{
			name: "Reassignment",
			input: TestMain(`y := 1
x := y
print(x)`, false, false),
			output: `set _main_y 1
set _main_x _main_y
print _main_x`,
		},
		{
			name: "VariableBooleans",
			input: TestMain(`x := false
print(x)`, false, false),
			output: `set _main_x false
print _main_x`,
		},
		{
			name: "VariableCharacter",
			input: TestMain(`x := 'A'
print(x)`, false, false),
			output: `set _main_x "A"
print _main_x`,
		},
		{
			name:  "Break",
			input: TestMain(`for i := 0; i < 10; i++ { if i == 5 { break; }; println(i); }`, false, false),
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
			input: TestMain(`for i := 0; i < 10; i++ { if i == 5 { continue; }; println(i); }`, false, false),
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
}`, false, false),
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
			input:  TestMain(`_ = false`, false, false),
			output: `set @_ false`,
		},
		{
			name: "OperatorAssign",
			input: TestMain(`x := 1
x += 1`, false, false),
			output: `set _main_x 1
op add _main_x _main_x 1`,
		},
		{
			name: "SelectorAssignment",
			input: TestMain(`a := m.RTAny
b := m.RSDistance
c := m.This
d := m.BCore
print(a)
print(b)
print(c)
print(d)`, true, false),
			output: `set _main_a any
set _main_b distance
set _main_c @this
set _main_d core
print _main_a
print _main_b
print _main_c
print _main_d`,
		},
		{
			name:  "MainReturn",
			input: TestMain(`x := 11; if x > 10 { return }; print(x)`, false, false),
			output: `set _main_x 11
op greaterThan _main_0 _main_x 10
jump 4 equal _main_0 1
jump 5 always
end
print _main_x`,
		},
		{
			name: "Labels",
			input: TestMain(`print(1)
	goto loop

test:
	print(2)
	goto end

loop:
	for i := 0; i < 10; i++ {
		print(3)
	}

	print(4)
	goto test
end:
	print(5)`, false, false),
			output: `print 1
jump loop
test:
print 2
jump end
loop:
set _main_i 0
jump 7 lessThan _main_i 10
jump 10 always
print 3
op add _main_i _main_i 1
jump 7 lessThan _main_i 10
print 4
jump test
end:
print 5`,
		},
	}
	RunTests(t, tests)
}
