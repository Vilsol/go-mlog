package transpiler

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBase(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name: "Read",
			input: TestMain(`x := m.Read("cell1", 0)
print(x)`, true, false),
			output: `read _main_x cell1 0
print _main_x`,
		},
		{
			name:   "Write",
			input:  TestMain(`m.Write(1, "cell1", 0)`, true, false),
			output: `write 1 cell1 0`,
		},
		{
			name:   "PrintFlush",
			input:  TestMain(`m.PrintFlush("message1")`, true, false),
			output: `printflush message1`,
		},
		{
			name: "GetLink",
			input: TestMain(`x := m.GetLink(0)
print(x)`, true, false),
			output: `getlink _main_x 0
print _main_x`,
		},
		{
			name: "Radar",
			input: TestMain(`x := m.Radar(m.This, m.RTAlly, m.RTEnemy, m.RTBoss, false, m.RSArmor)
print(x)`, true, false),
			output: `radar ally enemy boss armor @this false _main_x
print _main_x`,
		},
		{
			name: "Sensor",
			input: TestMain(`_, _, _, b := m.UnitLocateDamaged()
x := m.Sensor(b, "B")
print(x)`, true, false),
			output: `ulocate damaged core true @copper @_ @_ @_ _main_b
sensor _main_x _main_b B
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
