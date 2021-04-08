package impl

import (
	"github.com/Vilsol/go-mlog/m"
	"github.com/Vilsol/go-mlog/transpiler"
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
}
