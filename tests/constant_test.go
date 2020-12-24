package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestConstant(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name: "Constant",
			input: `package main

const x = 1
const y = x

func main() {
	print(x)
}`,
			output: `set @stack 0
set x 1
set y x
jump 4 always
print x`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mlog, err := transpiler.GolangToMLOG(test.input, transpiler.Options{})

			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
