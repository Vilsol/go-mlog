package transpiler

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
			assert.Equal(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
