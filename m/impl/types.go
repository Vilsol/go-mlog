package impl

import (
	"github.com/Vilsol/go-mlog/m"
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
)

func init() {
	transpiler.RegisterSelector("m.RTAny", m.RTAny)
	transpiler.RegisterSelector("m.RTEnemy", m.RTEnemy)
	transpiler.RegisterSelector("m.RTAlly", m.RTAlly)
	transpiler.RegisterSelector("m.RTPlayer", m.RTPlayer)
	transpiler.RegisterSelector("m.RTAttacker", m.RTAttacker)
	transpiler.RegisterSelector("m.RTFlying", m.RTFlying)
	transpiler.RegisterSelector("m.RTBoss", m.RTBoss)
	transpiler.RegisterSelector("m.RTGround", m.RTGround)

	transpiler.RegisterSelector("m.RSDistance", m.RSDistance)
	transpiler.RegisterSelector("m.RSHealth", m.RSHealth)
	transpiler.RegisterSelector("m.RSShield", m.RSShield)
	transpiler.RegisterSelector("m.RSArmor", m.RSArmor)
	transpiler.RegisterSelector("m.RSMaxHealth", m.RSMaxHealth)

	transpiler.RegisterSelector("m.BCore", m.BCore)
	transpiler.RegisterSelector("m.BStorage", m.BStorage)
	transpiler.RegisterSelector("m.BGenerator", m.BGenerator)
	transpiler.RegisterSelector("m.BTurret", m.BTurret)
	transpiler.RegisterSelector("m.BFactory", m.BFactory)
	transpiler.RegisterSelector("m.BRepair", m.BRepair)
	transpiler.RegisterSelector("m.BRally", m.BRally)
	transpiler.RegisterSelector("m.BBattery", m.BBattery)
	transpiler.RegisterSelector("m.BResupply", m.BResupply)
	transpiler.RegisterSelector("m.BReactor", m.BReactor)
	transpiler.RegisterSelector("m.BUnitModifier", m.BUnitModifier)
	transpiler.RegisterSelector("m.BExtinguisher", m.BExtinguisher)

	transpiler.RegisterSelector("m.This", "@this")
	transpiler.RegisterSelector("m.ThisX", "@thisx")
	transpiler.RegisterSelector("m.ThisXF", "@thisx")
	transpiler.RegisterSelector("m.ThisY", "@thisy")
	transpiler.RegisterSelector("m.ThisYF", "@thisy")
	transpiler.RegisterSelector("m.Ipt", "@ipt")
	transpiler.RegisterSelector("m.Counter", "@counter")
	transpiler.RegisterSelector("m.Links", "@links")
	transpiler.RegisterSelector("m.CurUnit", "@unit")
	transpiler.RegisterSelector("m.Time", "@time")
	transpiler.RegisterSelector("m.Tick", "@tick")
	transpiler.RegisterSelector("m.MapW", "@mapw")
	transpiler.RegisterSelector("m.MapH", "@maph")

	// HealthC's attributes
	transpiler.RegisterFuncTranslation("GetHealth", createSensorFuncTranslation("@health"))
	transpiler.RegisterFuncTranslation("GetName", createSensorFuncTranslation("@name"))
	transpiler.RegisterFuncTranslation("GetX", createSensorFuncTranslation("@x"))
	transpiler.RegisterFuncTranslation("GetY", createSensorFuncTranslation("@y"))

	transpiler.RegisterFuncTranslation("GetTotalItems", createSensorFuncTranslation("@totalItems"))
	transpiler.RegisterFuncTranslation("GetItemCapacity", createSensorFuncTranslation("@itemCapacity"))
	transpiler.RegisterFuncTranslation("GetRotation", createSensorFuncTranslation("@rotation"))
	transpiler.RegisterFuncTranslation("GetShootX", createSensorFuncTranslation("@shootX"))
	transpiler.RegisterFuncTranslation("GetShootY", createSensorFuncTranslation("@shootY"))
	transpiler.RegisterFuncTranslation("IsShooting", createSensorFuncTranslation("@shooting"))

	// Building's attributes
	transpiler.RegisterFuncTranslation("GetTotalLiquids", createSensorFuncTranslation("@totalLiquids"))
	transpiler.RegisterFuncTranslation("GetLiquidCapaticy", createSensorFuncTranslation("@liquidCapaticy"))
	transpiler.RegisterFuncTranslation("GetTotalPower", createSensorFuncTranslation("@totalPower"))
	transpiler.RegisterFuncTranslation("GetPowerCapaticy", createSensorFuncTranslation("@powerCapaticy"))
	transpiler.RegisterFuncTranslation("GetPowerNetStored", createSensorFuncTranslation("@powerNetStored"))
	transpiler.RegisterFuncTranslation("GetPowerNetCapacity", createSensorFuncTranslation("@powerNetCapacity"))
	transpiler.RegisterFuncTranslation("GetPowerNetIn", createSensorFuncTranslation("@powerNetIn"))
	transpiler.RegisterFuncTranslation("GetPowerNetOut", createSensorFuncTranslation("@powerNetOut"))
	transpiler.RegisterFuncTranslation("GetHeat", createSensorFuncTranslation("@heat"))
	transpiler.RegisterFuncTranslation("GetEfficiency", createSensorFuncTranslation("@efficiency"))
	transpiler.RegisterFuncTranslation("IsEnabled", createSensorFuncTranslation("@enabled"))
}

func createSensorFuncTranslation(attribute string) transpiler.Translator {
	return transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "sensor"},
							vars[0],
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
							&transpiler.Value{Value: attribute},
						},
					},
				},
			}, nil
		},
	}
}

func genBasicFuncTranslation(constants []string, nArgs int, nVars int) transpiler.TranslateFunc {
	return func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
		statements := make([]transpiler.Resolvable, len(constants)+nArgs+nVars)

		for i, constant := range constants {
			statements[i] = &transpiler.Value{Value: constant}
		}

		for i := 0; i < nArgs; i++ {
			statements[i+len(constants)] = &transpiler.Value{Value: args[i].GetValue()}
		}

		for i := 0; i < nVars; i++ {
			statements[i+len(constants)+nArgs] = vars[i]
		}

		return []transpiler.MLOGStatement{
			&transpiler.MLOG{
				Statement: [][]transpiler.Resolvable{statements},
			},
		}, nil
	}
}
