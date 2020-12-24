package tests

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
			name:   "Read",
			input:  TestMain(`m.Read("cell1", 0)`),
			output: `read @return cell1 0`,
		},
		{
			name:   "Write",
			input:  TestMain(`m.Write(1, "cell1", 0)`),
			output: `write 1 cell1 0`,
		},
		{
			name:   "PrintFlush",
			input:  TestMain(`m.PrintFlush("message1")`),
			output: `printflush message1`,
		},
		{
			name:   "GetLink",
			input:  TestMain(`m.GetLink(0)`),
			output: `getlink @return 0`,
		},
		{
			name:   "Radar",
			input:  TestMain(`m.Radar("A", m.RTAlly, m.RTEnemy, m.RTBoss, 0, m.RSArmor)`),
			output: `radar ally enemy boss armor "A" 0 @return`,
		},
		{
			name:   "Sensor",
			input:  TestMain(`m.Sensor("A", "B")`),
			output: `sensor @return "A" "B"`,
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
