package transpiler

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
	"testing"
)

func TestUnitControl(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "UnitStop",
			input:  TestMain(`m.UnitStop()`, true, false),
			output: `ucontrol stop`,
		},
		{
			name:   "UnitMove",
			input:  TestMain(`m.UnitMove(1, 2)`, true, false),
			output: `ucontrol move 1 2`,
		},
		{
			name:   "UnitApproach",
			input:  TestMain(`m.UnitApproach(1, 2, 3)`, true, false),
			output: `ucontrol approach 1 2 3`,
		},
		{
			name:   "UnitBoost",
			input:  TestMain(`m.UnitBoost(true)`, true, false),
			output: `ucontrol boost true`,
		},
		{
			name:   "UnitPathfind",
			input:  TestMain(`m.UnitPathfind()`, true, false),
			output: `ucontrol pathfind`,
		},
		{
			name:   "UnitTarget",
			input:  TestMain(`m.UnitTarget(1, 2, true)`, true, false),
			output: `ucontrol target 1 2 true`,
		},
		{
			name: "UnitTargetP",
			input: TestMain(`x := m.Radar(m.This, m.RTAlly, m.RTEnemy, m.RTBoss, false, m.RSArmor)
m.UnitTargetP(x, true)`, true, false),
			output: `radar ally enemy boss armor @this false _main_x
ucontrol targetp _main_x true`,
		},
		{
			name: "UnitItemDrop",
			input: TestMain(`_, _, _, b := m.UnitLocateDamaged()
m.UnitItemDrop(b, 2)`, true, false),
			output: `ulocate damaged core true @copper @_ @_ @_ _main_b
ucontrol itemDrop _main_b 2`,
		},
		{
			name: "UnitItemTake",
			input: TestMain(`_, _, _, b := m.UnitLocateDamaged()
m.UnitItemTake(b, "A", 2)`, true, false),
			output: `ulocate damaged core true @copper @_ @_ @_ _main_b
ucontrol itemTake _main_b "A" 2`,
		},
		{
			name:   "UnitPayloadDrop",
			input:  TestMain(`m.UnitPayloadDrop()`, true, false),
			output: `ucontrol payDrop`,
		},
		{
			name:   "UnitPayloadTake",
			input:  TestMain(`m.UnitPayloadTake(true)`, true, false),
			output: `ucontrol payTake true`,
		},
		{
			name:   "UnitMine",
			input:  TestMain(`m.UnitMine(1, 2)`, true, false),
			output: `ucontrol mine 1 2`,
		},
		{
			name:   "UnitFlag",
			input:  TestMain(`m.UnitFlag(1)`, true, false),
			output: `ucontrol flag 1`,
		},
		{
			name:   "UnitBuild",
			input:  TestMain(`m.UnitBuild(1, 2, "A", 3, 4)`, true, false),
			output: `ucontrol build 1 2 "A" 3 4`,
		},
		{
			name: "UnitGetBlock",
			input: TestMain(`x, y := m.UnitGetBlock(1, 2)
print(x)
print(y)`, true, false),
			output: `ucontrol getBlock 1 2 _main_x _main_y
print _main_x
print _main_y`,
		},
		{
			name: "UnitWithin",
			input: TestMain(`x := m.UnitWithin(1, 2, 3)
print(x)`, true, false),
			output: `ucontrol within 1 2 3 _main_x
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
			testza.AssertEqual(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
