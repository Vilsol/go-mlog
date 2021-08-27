package transpiler

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
	"testing"
)

func TestControl(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "ControlEnabled",
			input:  TestMain(`m.ControlEnabled("A", true)`, true, false),
			output: `control enabled "A" true`,
		},
		{
			name:   "ControlShoot",
			input:  TestMain(`m.ControlShoot("A", 3, 4, true)`, true, false),
			output: `control shoot "A" 3 4 true`,
		},
		{
			name:   "ControlShootP",
			input:  TestMain(`m.ControlShootP("A", 5, true)`, true, false),
			output: `control shootp "A" 5 true`,
		},
		{
			name:   "ControlConfigure",
			input:  TestMain(`m.ControlConfigure("A", 1)`, true, false),
			output: `control configure "A" 1`,
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
