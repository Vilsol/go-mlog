package transpiler

import (
	"testing"
)

func TestCustom(t *testing.T) {
	tests := []Test{
		{
			name: "Const",
			input: TestMain(`x := m.Const("@copper")
print(x)`, true, false),
			output: `set _main_x @copper
print _main_x`,
		},
		{
			name: "NestedSelector",
			input: TestMain(`x := m.This.GetX()
print(x)`, true, false),
			output: `sensor _main_x @this @x
print _main_x`,
		},
		{
			name:  "Inline",
			input: TestMain(`print(m.Const("@copper"), m.B("message1"), float64(1))`, true, false),
			output: `print @copper
print message1
print 1`,
		},
	}
	RunTests(t, tests)
}
