package transpiler

import (
	"testing"
)

func TestUnit(t *testing.T) {
	tests := []Test{
		{
			name:   "UnitBind",
			input:  TestMain(`m.UnitBind("A")`, true, false),
			output: `ubind A`,
		},
		{
			name: "UnitRadar",
			input: TestMain(`x := m.UnitRadar(m.RTAlly, m.RTEnemy, m.RTBoss, false, m.RSArmor)
print(x)`, true, false),
			output: `uradar ally enemy boss armor turret1 false _main_x
print _main_x`,
		},
		{
			name: "UnitLocateOre",
			input: TestMain(`x, y, z := m.UnitLocateOre("@copper")
print(x)
print(y)
print(z)`, true, false),
			output: `ulocate ore core true @copper _main_x _main_y _main_z null
print _main_x
print _main_y
print _main_z`,
		},
		{
			name: "UnitLocateBuilding",
			input: TestMain(`x, y, z, b := m.UnitLocateBuilding(m.BCore, true)
print(x)
print(y)
print(z)
print(b)`, true, false),
			output: `ulocate building core true @copper _main_x _main_y _main_z _main_b
print _main_x
print _main_y
print _main_z
print _main_b`,
		},
		{
			name: "UnitLocateSpawn",
			input: TestMain(`x, y, z, b := m.UnitLocateSpawn()
print(x)
print(y)
print(z)
print(b)`, true, false),
			output: `ulocate spawn core true @copper _main_x _main_y _main_z _main_b
print _main_x
print _main_y
print _main_z
print _main_b`,
		},
		{
			name: "UnitLocateDamaged",
			input: TestMain(`x, y, z, b := m.UnitLocateDamaged()
print(x)
print(y)
print(z)
print(b)`, true, false),
			output: `ulocate damaged core true @copper _main_x _main_y _main_z _main_b
print _main_x
print _main_y
print _main_z
print _main_b`,
		},
	}
	RunTests(t, tests)
}
