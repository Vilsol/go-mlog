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
	}
	RunTests(t, tests)
}
