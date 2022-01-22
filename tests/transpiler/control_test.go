package transpiler

import (
	"testing"
)

func TestControl(t *testing.T) {
	tests := []Test{
		{
			name:   "ControlEnabled",
			input:  TestMain(`m.ControlEnabled(m.B("A"), true)`, true, false),
			output: `control enabled A true`,
		},
		{
			name:   "ControlShoot",
			input:  TestMain(`m.ControlShoot(m.B("A"), 3, 4, true)`, true, false),
			output: `control shoot A 3 4 true`,
		},
		{
			name:   "ControlShootP",
			input:  TestMain(`m.ControlShootP(m.B("A"), m.B("B"), true)`, true, false),
			output: `control shootp A B true`,
		},
		{
			name:   "ControlConfigure",
			input:  TestMain(`m.ControlConfigure(m.B("A"), 1)`, true, false),
			output: `control configure A 1`,
		},
	}
	RunTests(t, tests)
}
