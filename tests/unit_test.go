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
			input:  TestMain(`x := m.UnitRadar(m.RTAlly, m.RTEnemy, m.RTBoss, 0, m.RSArmor)`),
			output: `uradar ally enemy boss armor turret1 0 _main_x`,
		},
		{
			name:   "UnitLocateOre",
			input:  TestMain(`x, y, z := m.UnitLocateOre("@copper")`),
			output: `ulocate ore core true @copper _main_x _main_y _main_z null`,
		},
		{
			name:   "UnitLocateBuilding",
			input:  TestMain(`x, y, z, b := m.UnitLocateBuilding(m.BCore, 1)`),
			output: `ulocate building core 1 @copper _main_x _main_y _main_z _main_b`,
		},
		{
			name:   "UnitLocateSpawn",
			input:  TestMain(`x, y, z, b := m.UnitLocateSpawn()`),
			output: `ulocate spawn core true @copper _main_x _main_y _main_z _main_b`,
		},
		{
			name:   "UnitLocateDamaged",
			input:  TestMain(`x, y, z, b := m.UnitLocateDamaged()`),
			output: `ulocate damaged core true @copper _main_x _main_y _main_z _main_b`,
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
