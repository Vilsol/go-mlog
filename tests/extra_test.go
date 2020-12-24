package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestExtra(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:  "Sleep",
			input: TestMain(`x.Sleep(1000)`),
			output: `op add _main_0 @time 1000
jump 1 lessThan @time _main_0`,
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
