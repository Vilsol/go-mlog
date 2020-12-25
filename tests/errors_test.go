package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "InvalidInput",
			input:  `hello world`,
			output: `foo:1:1: expected 'package', found hello`,
		},
		{
			name:   "PackageMustBeMain",
			input:  `package foo`,
			output: `package must be main`,
		},
		{
			name: "NoExternalImports",
			input: `package main
import "time"`,
			output: `unregistered import used: "time"`,
		},
		{
			name: "GlobalScopeVariable",
			input: `package main
var x = 1`,
			output: `global scope may only contain constants not variables`,
		},
		{
			name:   "NoMainFunction",
			input:  `package main`,
			output: `file does not contain a main function`,
		},
		{
			name:   "InvalidOperator",
			input:  TestMain(`x := 1 &^ 1`),
			output: `operator statement cannot use this operation: &^`,
		},
		{
			name:   "NotSupportSelect",
			input:  TestMain(`select {}`),
			output: `statement type not supported: *ast.SelectStmt`,
		},
		{
			name:   "NotSupportSwitch",
			input:  TestMain(`switch {}`),
			output: `statement type not supported: *ast.SwitchStmt`,
		},
		{
			name:   "NotSupportGo",
			input:  TestMain(`go foo()`),
			output: `statement type not supported: *ast.GoStmt`,
		},
		{
			name:   "NotSupportSend",
			input:  TestMain(`foo <- 1`),
			output: `statement type not supported: *ast.SendStmt`,
		},
		{
			name:   "NotSupportDefer",
			input:  TestMain(`defer func() {}()`),
			output: `statement type not supported: *ast.DeferStmt`,
		},
		{
			name:   "NotSupportRange",
			input:  TestMain(`for i := range nums {}`),
			output: `statement type not supported: *ast.RangeStmt`,
		},
		{
			name:   "InvalidAssignment",
			input:  TestMain(`1 = 2`),
			output: `left side variable assignment can only contain identifications`,
		},
		{
			name: "InvalidParamTypeString",
			input: `package main

func main() {
	print(sample1("hello"))
}

func sample1(arg string) int {
	return 1
}`,
			output: `function parameters may only be integers or floating point numbers`,
		},
		{
			name: "InvalidParamTypeOther",
			input: `package main

func main() {
	print(sample1(nil))
}

func sample1(arg hello.world) int {
	return 1
}`,
			output: `function parameters may only be integers or floating point numbers`,
		},
		{
			name:   "CallToUnknownFunction",
			input:  TestMain(`foo()`),
			output: `unknown function: foo`,
		},
		{
			name: "InvalidConstant",
			input: `package main

const x = 1 + 2

func main() {
}`,
			output: `unknown constant type: *ast.BinaryExpr`,
		},
		{
			name:   "EmptyPrintlnError",
			input:  TestMain(`println()`),
			output: `println with 0 arguments`,
		},
		{
			name:   "EmptyPrintError",
			input:  TestMain(`print()`),
			output: `print with 0 arguments`,
		},
		{
			name:   "ErrorWriteToStack",
			input:  TestMain(`m.Write(0, "bank1", 0)`),
			output: `can't read/write to memory cell that is used for the stack: bank1`,
		},
		{
			name:   "ErrorReadFromStack",
			input:  TestMain(`x := m.Read("bank1", 0)`),
			output: `can't read/write to memory cell that is used for the stack: bank1`,
		},
		{
			name: "SingleValueReturns",
			input: `package main

func main() {
	sample()
}

func sample() (int, int) {
	return 1, 2
}`,
			output: `only single value returns are supported`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := transpiler.GolangToMLOG(test.input, transpiler.Options{})
			assert.EqualError(t, err, test.output)
		})
	}
}

func TestRegisterSelectorPanic(t *testing.T) {
	assert.Panics(t, func() {
		transpiler.RegisterSelector("m.RTAny", "any")
	})
}

func TestRegisterFuncTranslationPanic(t *testing.T) {
	assert.Panics(t, func() {
		transpiler.RegisterFuncTranslation("print", transpiler.Translator{})
	})
}
