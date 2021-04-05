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
			input:  TestMain(`x := m.Read("cell1", 0)`),
			output: `read _main_x cell1 0`,
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
			input:  TestMain(`x := m.GetLink(0)`),
			output: `getlink _main_x 0`,
		},
		{
			name:   "Radar",
			input:  TestMain(`x := m.Radar("A", m.RTAlly, m.RTEnemy, m.RTBoss, 0, m.RSArmor)`),
			output: `radar ally enemy boss armor "A" 0 _main_x`,
		},
		{
			name:   "Sensor",
			input:  TestMain(`x := m.Sensor("A", "B")`),
			output: `sensor _main_x A B`,
		},
		{
			name:   "Sensor_GetHealth",
			input:  TestMain(`x := y.GetHealth()`),
			output: `sensor _main_x _main_y @health`,
		},
		{
			name:   "Radar_This",
			input:  TestMain(`x := m.Radar(m.This, m.RTAlly, m.RTEnemy, m.RTBoss, 0, m.RSArmor)`),
			output: `radar ally enemy boss armor @this 0 _main_x`,
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
