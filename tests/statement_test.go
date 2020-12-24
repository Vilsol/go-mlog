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
print _main_i
op add _main_i _main_i 1
jump 1 lessThan _main_i 10`,
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
