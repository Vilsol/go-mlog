package m

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterSelector("m.RTAny", RTAny)
	transpiler.RegisterSelector("m.RTEnemy", RTEnemy)
	transpiler.RegisterSelector("m.RTAlly", RTAlly)
	transpiler.RegisterSelector("m.RTPlayer", RTPlayer)
	transpiler.RegisterSelector("m.RTAttacker", RTAttacker)
	transpiler.RegisterSelector("m.RTFlying", RTFlying)
	transpiler.RegisterSelector("m.RTBoss", RTBoss)
	transpiler.RegisterSelector("m.RTGround", RTGround)

	transpiler.RegisterSelector("m.RSDistance", RSDistance)
	transpiler.RegisterSelector("m.RSHealth", RSHealth)
	transpiler.RegisterSelector("m.RSShield", RSShield)
	transpiler.RegisterSelector("m.RSArmor", RSArmor)
	transpiler.RegisterSelector("m.RSMaxHealth", RSMaxHealth)

	transpiler.RegisterSelector("m.BCore", BCore)
	transpiler.RegisterSelector("m.BStorage", BStorage)
	transpiler.RegisterSelector("m.BGenerator", BGenerator)
	transpiler.RegisterSelector("m.BTurret", BTurret)
	transpiler.RegisterSelector("m.BFactory", BFactory)
	transpiler.RegisterSelector("m.BRepair", BRepair)
	transpiler.RegisterSelector("m.BRally", BRally)
	transpiler.RegisterSelector("m.BBattery", BBattery)
	transpiler.RegisterSelector("m.BResupply", BResupply)
	transpiler.RegisterSelector("m.BReactor", BReactor)
	transpiler.RegisterSelector("m.BUnitModifier", BUnitModifier)
	transpiler.RegisterSelector("m.BExtinguisher", BExtinguisher)
}

type RadarTarget = string

const (
	// Target anything
	RTAny      = RadarTarget("any")
	RTEnemy    = RadarTarget("enemy")
	RTAlly     = RadarTarget("ally")
	RTPlayer   = RadarTarget("player")
	RTAttacker = RadarTarget("attacker")
	RTFlying   = RadarTarget("flying")
	RTBoss     = RadarTarget("boss")
	RTGround   = RadarTarget("ground")
)

type RadarSort = string

const (
	RSDistance  = RadarSort("distance")
	RSHealth    = RadarSort("health")
	RSShield    = RadarSort("shield")
	RSArmor     = RadarSort("armor")
	RSMaxHealth = RadarSort("maxHealth")
)

type Link = interface{}

type Unit = interface{}

type HealthC = interface{}

type Building = interface{}

type BlockFlag = string

const (
	BCore         = BlockFlag("core")
	BStorage      = BlockFlag("storage")
	BGenerator    = BlockFlag("generator")
	BTurret       = BlockFlag("turret")
	BFactory      = BlockFlag("factory")
	BRepair       = BlockFlag("repair")
	BRally        = BlockFlag("rally")
	BBattery      = BlockFlag("battery")
	BResupply     = BlockFlag("resupply")
	BReactor      = BlockFlag("reactor")
	BUnitModifier = BlockFlag("unitModifier")
	BExtinguisher = BlockFlag("extinguisher")
)
