package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestJumpOperator(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:  "equal",
			input: TestMain(`if 1 == 2 { print(1) }`),
			output: `op equal _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "notEqual",
			input: TestMain(`if 1 != 2 { print(1) }`),
			output: `op notEqual _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "greaterThan",
			input: TestMain(`if 1 > 2 { print(1) }`),
			output: `op greaterThan _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "greaterThanEq",
			input: TestMain(`if 1 >= 2 { print(1) }`),
			output: `op greaterThanEq _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "lessThan",
			input: TestMain(`if 1 < 2 { print(1) }`),
			output: `op lessThan _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "lessThanEq",
			input: TestMain(`if 1 <= 2 { print(1) }`),
			output: `op lessThanEq _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
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

func TestNormalOperator(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "add",
			input:  TestMain(`x := 1 + 2`),
			output: `op add _main_x 1 2`,
		},
		{
			name:   "sub",
			input:  TestMain(`x := 1 - 2`),
			output: `op sub _main_x 1 2`,
		},
		{
			name:   "mul",
			input:  TestMain(`x := 1 * 2`),
			output: `op mul _main_x 1 2`,
		},
		{
			name:   "div",
			input:  TestMain(`x := 1 / 2`),
			output: `op div _main_x 1 2`,
		},
		{
			name:   "mod",
			input:  TestMain(`x := 1 % 2`),
			output: `op mod _main_x 1 2`,
		},
		{
			name:   "equal",
			input:  TestMain(`x := 1 == 2`),
			output: `op equal _main_x 1 2`,
		},
		{
			name:   "notEqual",
			input:  TestMain(`x := 1 != 2`),
			output: `op notEqual _main_x 1 2`,
		},
		{
			name:   "lessThan",
			input:  TestMain(`x := 1 < 2`),
			output: `op lessThan _main_x 1 2`,
		},
		{
			name:   "lessThanEq",
			input:  TestMain(`x := 1 <= 2`),
			output: `op lessThanEq _main_x 1 2`,
		},
		{
			name:   "greaterThan",
			input:  TestMain(`x := 1 > 2`),
			output: `op greaterThan _main_x 1 2`,
		},
		{
			name:   "greaterThanEq",
			input:  TestMain(`x := 1 >= 2`),
			output: `op greaterThanEq _main_x 1 2`,
		},
		{
			name:  "land",
			input: TestMain(`x := 1 == 2 && 3 == 4`),
			output: `op equal _main_0 1 2
op equal _main_1 3 4
op land _main_x _main_0 _main_1`,
		},
		{
			name:  "lor",
			input: TestMain(`x := 1 == 2 || 3 == 4`),
			output: `op equal _main_0 1 2
op equal _main_1 3 4
op or _main_x _main_0 _main_1`,
		},
		{
			name:   "shl",
			input:  TestMain(`x := 1 << 2`),
			output: `op shl _main_x 1 2`,
		},
		{
			name:   "shr",
			input:  TestMain(`x := 1 >> 2`),
			output: `op shr _main_x 1 2`,
		},
		{
			name:   "or",
			input:  TestMain(`x := 1 | 2`),
			output: `op or _main_x 1 2`,
		},
		{
			name:   "and",
			input:  TestMain(`x := 1 & 2`),
			output: `op and _main_x 1 2`,
		},
		{
			name:   "xor",
			input:  TestMain(`x := 1 ^ 2`),
			output: `op xor _main_x 1 2`,
		},
		{
			name:  "not",
			input: TestMain(`x := !(1 == 2)`),
			output: `op equal _main_0 1 2
op not _main_x _main_0 -1`,
		},
		{
			name:   "negative",
			input:  TestMain(`x := -5`),
			output: `op mul _main_x 5 -1`,
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
