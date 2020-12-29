package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNative(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:  "print",
			input: TestMain(`print(1, "A", x)`),
			output: `print 1
print "A"
print _main_x`,
		},
		{
			name:  "println",
			input: TestMain(`println(1, "A", x)`),
			output: `print 1
print "A"
print _main_x
print "\n"`,
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
