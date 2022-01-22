package transpiler

import (
	"testing"
)

func TestNative(t *testing.T) {
	tests := []Test{
		{
			name: "print",
			input: TestMain(`x := 2
print(1, "A", x)`, false, false),
			output: `set _main_x 2
print 1
print "A"
print _main_x`,
		},
		{
			name: "println",
			input: TestMain(`x := 2
println(1, "A", x)`, false, false),
			output: `set _main_x 2
print 1
print "A"
print _main_x
print "\n"`,
		},
		{
			name: "float64",
			input: TestMain(`x := float64(1)
print(x)`, false, false),
			output: `set _main_x 1
print _main_x`,
		},
	}
	RunTests(t, tests)
}
