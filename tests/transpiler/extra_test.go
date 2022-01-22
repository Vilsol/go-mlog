package transpiler

import (
	"testing"
)

func TestExtra(t *testing.T) {
	tests := []Test{
		{
			name:  "Sleep",
			input: TestMain(`x.Sleep(1000)`, false, true),
			output: `op add _main_0 @time 1000
jump 1 lessThan @time _main_0`,
		},
	}
	RunTests(t, tests)
}
