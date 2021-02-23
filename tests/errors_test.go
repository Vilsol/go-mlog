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
			output: `error at 21: unregistered import used: "time"`,
		},
		{
			name: "GlobalScopeVariable",
			input: `package main
var x = 1`,
			output: `error at 14: global scope may only contain constants not variables`,
		},
		{
			name:   "NoMainFunction",
			input:  `package main`,
			output: `file does not contain a main function`,
		},
		{
			name:   "InvalidOperator",
			input:  TestMain(`x := 1 &^ 1`),
			output: `error at 103: operator statement cannot use this operation: &^`,
		},
		{
			name:   "NotSupportSelect",
			input:  TestMain(`select {}`),
			output: `error at 103: statement type not supported: *ast.SelectStmt`,
		},
		{
			name:   "NotSupportGo",
			input:  TestMain(`go foo()`),
			output: `error at 103: statement type not supported: *ast.GoStmt`,
		},
		{
			name:   "NotSupportSend",
			input:  TestMain(`foo <- 1`),
			output: `error at 103: statement type not supported: *ast.SendStmt`,
		},
		{
			name:   "NotSupportDefer",
			input:  TestMain(`defer func() {}()`),
			output: `error at 103: statement type not supported: *ast.DeferStmt`,
		},
		{
			name:   "NotSupportRange",
			input:  TestMain(`for i := range nums {}`),
			output: `error at 103: statement type not supported: *ast.RangeStmt`,
		},
		{
			name:   "InvalidAssignment",
			input:  TestMain(`1 = 2`),
			output: `error at 103: left side variable assignment can only contain identifications`,
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
			output: `error at 53: function parameters may only be basic types`,
		},
		{
			name:   "CallToUnknownFunction",
			input:  TestMain(`foo()`),
			output: `error at 103-108: unknown function: foo`,
		},
		{
			name: "InvalidConstant",
			input: `package main

const x = 1 + 2

func main() {
	println("1")
}`,
			output: `error at 21: unknown constant type: *ast.BinaryExpr`,
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
			output: `error at 70: only single value returns are supported`,
		},
		{
			name:   "ErrorEmptyMain",
			input:  TestMain(``),
			output: `empty main function`,
		},
		{
			name:   "ErrorIncorrectVarCount",
			input:  TestMain(`a, b := m.Read("bank1", 0)`),
			output: `error at 111-129: function requires 1 variables, provided: 2`,
		},
		{
			name:   "ErrorMismatchedSides",
			input:  TestMain(`a, b, c := 1, 2`),
			output: `error at 103: mismatched variable assignment sides`,
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
