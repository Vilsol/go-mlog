package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
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
			input:  TestMain(`m.ControlEnabled("A", true)`),
			output: `control enabled "A" true`,
		},
		{
			name:   "ControlShoot",
			input:  TestMain(`m.ControlShoot("A", 3, 4, true)`),
			output: `control shoot "A" 3 4 true`,
		},
		{
			name:   "ControlShootP",
			input:  TestMain(`m.ControlShootP("A", 5, true)`),
			output: `control shootp "A" 5 true`,
		},
		{
			name:   "ControlConfigure",
			input:  TestMain(`m.ControlConfigure("A", 1)`),
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

			assert.Equal(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
