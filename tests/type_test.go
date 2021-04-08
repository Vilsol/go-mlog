package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestType(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name: "Radar_This",
			input: TestMain(`x := m.Radar(m.This, m.RTAlly, m.RTEnemy, m.RTBoss, false, m.RSArmor)
				print(x)`, true, false),
			output: `radar ally enemy boss armor @this false _main_x
print _main_x`,
		},
		{
			name: "Sensor_GetHealth",
			input: TestMain(`_, _, _, b := m.UnitLocateDamaged()
x := b.GetHealth()
print(x)`, true, false),
			output: `ulocate damaged core true @copper @_ @_ @_ _main_b
sensor _main_x _main_b @health
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
