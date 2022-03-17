package transpiler

import (
	"testing"
)

func TestType(t *testing.T) {
	tests := []Test{
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
		{
			name: "Access_Unit",
			input: TestMain(
				`flg := m.CurUnit.GetFlag()
print(flg)`,
				true, false),
			output: `sensor _main_flg @unit @flag
print _main_flg`,
		},
		{
			name: "Access_Turret",
			input: TestMain(
				`turret := m.GetTurret("duo1")
capa := turret.GetAmmoCapacity()
print(capa)`,
				true, false),
			output: `set _main_turret duo1
sensor _main_capa _main_turret @ammoCapacity
print _main_capa`,
		},
		{
			name:   "Bind_constant",
			input:  TestMain("m.UnitBind(m.UZenith)", true, false),
			output: "ubind @zenith",
		},
		{
			name: "Compare_ItemTypes",
			input: TestMain(`item := m.CurUnit.GetFirstItem()
if(item == m.ITSilicon) {
	print("Test")
}`, true, false),
			output: `sensor _main_item @unit @firstItem
set _main_0 @silicon
op equal _main_1 _main_item _main_0
jump 5 equal _main_1 1
jump 6 always
print "Test"`,
		},
		{
			name: "Compare_ItemTypes_old",
			input: TestMain(`item := m.SensorStr(m.CurUnit, "@firstItem")
if item == m.Const("@silicon") {
	print("Test")
}`, true, false),
			output: `sensor _main_item @unit @firstItem
set _main_0 @silicon
op equal _main_1 _main_item _main_0
jump 5 equal _main_1 1
jump 6 always
print "Test"`,
		},
	}
	RunTests(t, tests)
}
