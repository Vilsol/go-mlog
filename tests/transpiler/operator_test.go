package transpiler

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/transpiler"
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
			input: TestMain(`if 1 == 2 { print(1) }`, false, false),
			output: `op equal _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "notEqual",
			input: TestMain(`if 1 != 2 { print(1) }`, false, false),
			output: `op notEqual _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "greaterThan",
			input: TestMain(`if 1 > 2 { print(1) }`, false, false),
			output: `op greaterThan _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "greaterThanEq",
			input: TestMain(`if 1 >= 2 { print(1) }`, false, false),
			output: `op greaterThanEq _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "lessThan",
			input: TestMain(`if 1 < 2 { print(1) }`, false, false),
			output: `op lessThan _main_0 1 2
jump 3 equal _main_0 1
jump 4 always
print 1`,
		},
		{
			name:  "lessThanEq",
			input: TestMain(`if 1 <= 2 { print(1) }`, false, false),
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

			test.output = test.output + "\nend"
			testza.AssertEqual(t, test.output, strings.Trim(mlog, "\n"))
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
			name: "add",
			input: TestMain(`x := 1 + 2
print(x)`, false, false),
			output: `op add _main_x 1 2
print _main_x`,
		},
		{
			name: "sub",
			input: TestMain(`x := 1 - 2
print(x)`, false, false),
			output: `op sub _main_x 1 2
print _main_x`,
		},
		{
			name: "mul",
			input: TestMain(`x := 1 * 2
print(x)`, false, false),
			output: `op mul _main_x 1 2
print _main_x`,
		},
		{
			name: "div",
			input: TestMain(`x := 1 / 2
print(x)`, false, false),
			output: `op div _main_x 1 2
print _main_x`,
		},
		{
			name: "mod",
			input: TestMain(`x := 1 % 2
print(x)`, false, false),
			output: `op mod _main_x 1 2
print _main_x`,
		},
		{
			name: "equal",
			input: TestMain(`x := 1 == 2
print(x)`, false, false),
			output: `op equal _main_x 1 2
print _main_x`,
		},
		{
			name: "notEqual",
			input: TestMain(`x := 1 != 2
print(x)`, false, false),
			output: `op notEqual _main_x 1 2
print _main_x`,
		},
		{
			name: "lessThan",
			input: TestMain(`x := 1 < 2
print(x)`, false, false),
			output: `op lessThan _main_x 1 2
print _main_x`,
		},
		{
			name: "lessThanEq",
			input: TestMain(`x := 1 <= 2
print(x)`, false, false),
			output: `op lessThanEq _main_x 1 2
print _main_x`,
		},
		{
			name: "greaterThan",
			input: TestMain(`x := 1 > 2
print(x)`, false, false),
			output: `op greaterThan _main_x 1 2
print _main_x`,
		},
		{
			name: "greaterThanEq",
			input: TestMain(`x := 1 >= 2
print(x)`, false, false),
			output: `op greaterThanEq _main_x 1 2
print _main_x`,
		},
		{
			name: "land",
			input: TestMain(`x := 1 == 2 && 3 == 4
print(x)`, false, false),
			output: `op equal _main_0 1 2
op equal _main_1 3 4
op land _main_x _main_0 _main_1
print _main_x`,
		},
		{
			name: "lor",
			input: TestMain(`x := 1 == 2 || 3 == 4
print(x)`, false, false),
			output: `op equal _main_0 1 2
op equal _main_1 3 4
op or _main_x _main_0 _main_1
print _main_x`,
		},
		{
			name: "shl",
			input: TestMain(`x := 1 << 2
print(x)`, false, false),
			output: `op shl _main_x 1 2
print _main_x`,
		},
		{
			name: "shr",
			input: TestMain(`x := 1 >> 2
print(x)`, false, false),
			output: `op shr _main_x 1 2
print _main_x`,
		},
		{
			name: "or",
			input: TestMain(`x := 1 | 2
print(x)`, false, false),
			output: `op or _main_x 1 2
print _main_x`,
		},
		{
			name: "and",
			input: TestMain(`x := 1 & 2
print(x)`, false, false),
			output: `op and _main_x 1 2
print _main_x`,
		},
		{
			name: "xor",
			input: TestMain(`x := 1 ^ 2
print(x)`, false, false),
			output: `op xor _main_x 1 2
print _main_x`,
		},
		{
			name: "not",
			input: TestMain(`x := !(1 == 2)
print(x)`, false, false),
			output: `op equal _main_0 1 2
op not _main_x _main_0 -1
print _main_x`,
		},
		{
			name: "negative",
			input: TestMain(`x := -5
print(x)`, false, false),
			output: `op mul _main_x 5 -1
print _main_x`,
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

			test.output = test.output + "\nend"
			testza.AssertEqual(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}

func TestFunctionOperator(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name: "Floor",
			input: TestMain(`x := m.Floor(1.2)
print(x)`, true, false),
			output: `op floor _main_x 1.2
print _main_x`,
		},
		{
			name: "Ceil",
			input: TestMain(`x := m.Ceil(1.2)
print(x)`, true, false),
			output: `op ceil _main_x 1.2
print _main_x`,
		},
		{
			name: "Random",
			input: TestMain(`x := m.Random(1.2)
print(x)`, true, false),
			output: `op rand _main_x 1.2
print _main_x`,
		},
		{
			name: "Log10",
			input: TestMain(`x := m.Log10(1.2)
print(x)`, true, false),
			output: `op log10 _main_x 1.2
print _main_x`,
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

			test.output = test.output + "\nend"
			testza.AssertEqual(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
