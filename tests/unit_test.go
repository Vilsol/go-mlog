package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnit(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "UnitBind",
			input:  TestMain(`m.UnitBind("A")`),
			output: `ubind "A"`,
		},
		{
			name:   "UnitRadar",
			input:  TestMain(`m.UnitRadar(m.RTAlly, m.RTEnemy, m.RTBoss, 0, m.RSArmor)`),
			output: `radar ally enemy boss armor turret1 0 @return`,
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
