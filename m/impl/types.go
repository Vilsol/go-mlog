package impl

import (
	"github.com/Vilsol/go-mlog/m"
	"github.com/Vilsol/go-mlog/transpiler"
	"strconv"
	"strings"
)

func init() {
	registerConstants()
	registerUnitEnum()
	registerItemEnum()
	registerNativeVariables()
	registerArchtypeMethods()
}

func registerArchtypeMethods() {
	transpiler.RegisterFuncTranslation("IsEnabled", createSensorFuncTranslation("@enabled"))
	transpiler.RegisterFuncTranslation("GetHealth", createSensorFuncTranslation("@health"))
	transpiler.RegisterFuncTranslation("GetMaxHealth", createSensorFuncTranslation("@maxHealth"))
	transpiler.RegisterFuncTranslation("GetTotalItems", createSensorFuncTranslation("@totalItems"))
	transpiler.RegisterFuncTranslation("GetFirstItem", createSensorFuncTranslation("@firstItem"))
	transpiler.RegisterFuncTranslation("GetItemCapacity", createSensorFuncTranslation("@itemCapacity"))
	transpiler.RegisterFuncTranslation("GetTotalLiquids", createSensorFuncTranslation("@totalLiquids"))
	transpiler.RegisterFuncTranslation("GetLiquidCapacity", createSensorFuncTranslation("@liquidCapacity"))
	transpiler.RegisterFuncTranslation("GetTotalPower", createSensorFuncTranslation("@totalPower"))
	transpiler.RegisterFuncTranslation("GetPowerCapacity", createSensorFuncTranslation("@powerCapacity"))
	transpiler.RegisterFuncTranslation("GetPowerNetStored", createSensorFuncTranslation("@powerNetStored"))
	transpiler.RegisterFuncTranslation("GetPowerNetCapacity", createSensorFuncTranslation("@powerNetCapacity"))
	transpiler.RegisterFuncTranslation("GetPowerNetIn", createSensorFuncTranslation("@powerNetIn"))
	transpiler.RegisterFuncTranslation("GetPowerNetOut", createSensorFuncTranslation("@powerNetOut"))
	transpiler.RegisterFuncTranslation("GetAmmo", createSensorFuncTranslation("@ammo"))
	transpiler.RegisterFuncTranslation("GetAmmoCapacity", createSensorFuncTranslation("@ammoCapacity"))
	transpiler.RegisterFuncTranslation("IsShooting", createSensorFuncTranslation("@shooting"))
	transpiler.RegisterFuncTranslation("GetHeat", createSensorFuncTranslation("@heat"))
	transpiler.RegisterFuncTranslation("GetEfficiency", createSensorFuncTranslation("@Efficiency"))
	transpiler.RegisterFuncTranslation("GetX", createSensorFuncTranslation("@x"))
	transpiler.RegisterFuncTranslation("GetY", createSensorFuncTranslation("@y"))
	transpiler.RegisterFuncTranslation("GetSize", createSensorFuncTranslation("@size"))
	transpiler.RegisterFuncTranslation("GetRotation", createSensorFuncTranslation("@rotation"))
	transpiler.RegisterFuncTranslation("GetShootX", createSensorFuncTranslation("@shootX"))
	transpiler.RegisterFuncTranslation("GetShootY", createSensorFuncTranslation("@shootY"))
	transpiler.RegisterFuncTranslation("GetShootPosition", createSensorFuncTranslation("@shootPosition"))
	transpiler.RegisterFuncTranslation("GetControlled", createSensorFuncTranslation("@controlled"))
	transpiler.RegisterFuncTranslation("GetController", createSensorFuncTranslation("@controller"))
	transpiler.RegisterFuncTranslation("GetType", createSensorFuncTranslation("@type"))
	transpiler.RegisterFuncTranslation("IsDead", createSensorFuncTranslation("@dead"))
	transpiler.RegisterFuncTranslation("IsBoosting", createSensorFuncTranslation("@boosting"))
	transpiler.RegisterFuncTranslation("GetFlag", createSensorFuncTranslation("@flag"))
	transpiler.RegisterFuncTranslation("GetTeam", createSensorFuncTranslation("@team"))
	transpiler.RegisterFuncTranslation("GetRange", createSensorFuncTranslation("@range"))
	transpiler.RegisterFuncTranslation("GetMineX", createSensorFuncTranslation("@mineX"))
	transpiler.RegisterFuncTranslation("GetMineY", createSensorFuncTranslation("@mineY"))
	transpiler.RegisterFuncTranslation("IsMining", createSensorFuncTranslation("@mining"))
	transpiler.RegisterFuncTranslation("GetName", createSensorFuncTranslation("@name"))
	transpiler.RegisterFuncTranslation("GetConfig", createSensorFuncTranslation("@configure"))
	transpiler.RegisterFuncTranslation("GetPayloadType", createSensorFuncTranslation("@payloadType"))
	transpiler.RegisterFuncTranslation("GetPayloadCount", createSensorFuncTranslation("@payloadCount"))
}

func registerItemEnum() {
	transpiler.RegisterSelector("m.ITCopper", m.ITCopper)
	transpiler.RegisterSelector("m.ITLead", m.ITLead)
	transpiler.RegisterSelector("m.ITMetaGlass", m.ITMetaGlass)
	transpiler.RegisterSelector("m.ITGraphite", m.ITGraphite)
	transpiler.RegisterSelector("m.ITSand", m.ITSand)
	transpiler.RegisterSelector("m.ITCoal", m.ITCoal)
	transpiler.RegisterSelector("m.ITTitanium", m.ITTitanium)
	transpiler.RegisterSelector("m.ITThorium", m.ITThorium)
	transpiler.RegisterSelector("m.ITScrap", m.ITScrap)
	transpiler.RegisterSelector("m.ITSilicon", m.ITSilicon)
	transpiler.RegisterSelector("m.ITPlastanium", m.ITPlastanium)
	transpiler.RegisterSelector("m.ITPhaseFabric", m.ITPhaseFabric)
	transpiler.RegisterSelector("m.ITSurgeAlloy", m.ITSurgeAlloy)
	transpiler.RegisterSelector("m.ITSporePod", m.ITSporePod)
	transpiler.RegisterSelector("m.ITBlastCompound", m.ITBlastCompound)
	transpiler.RegisterSelector("m.ITPyratite", m.ITPyratite)
	transpiler.RegisterSelector("m.ITNone", m.ITNone)

	transpiler.RegisterSelector("m.FWater", m.FWater)
	transpiler.RegisterSelector("m.FOil", m.FOil)
	transpiler.RegisterSelector("m.FSlag", m.FSlag)
	transpiler.RegisterSelector("m.FCryofluid", m.FCryofluid)
}

func registerUnitEnum() {
	transpiler.RegisterSelector("m.UDagger", m.UDagger)
	transpiler.RegisterSelector("m.UMace", m.UMace)
	transpiler.RegisterSelector("m.UFortress", m.UFortress)
	transpiler.RegisterSelector("m.UScepter", m.UScepter)
	transpiler.RegisterSelector("m.UReign", m.UReign)
	transpiler.RegisterSelector("m.UNova", m.UNova)
	transpiler.RegisterSelector("m.UPulsar", m.UPulsar)
	transpiler.RegisterSelector("m.UQuasar", m.UQuasar)
	transpiler.RegisterSelector("m.UVela", m.UVela)
	transpiler.RegisterSelector("m.UCorvus", m.UCorvus)
	transpiler.RegisterSelector("m.UCrawler", m.UCrawler)
	transpiler.RegisterSelector("m.UAtrax", m.UAtrax)
	transpiler.RegisterSelector("m.USpiroct", m.USpiroct)
	transpiler.RegisterSelector("m.UArkyid", m.UArkyid)
	transpiler.RegisterSelector("m.UToxopid", m.UToxopid)
	transpiler.RegisterSelector("m.UFlare", m.UFlare)
	transpiler.RegisterSelector("m.UHorizon", m.UHorizon)
	transpiler.RegisterSelector("m.UZenith", m.UZenith)
	transpiler.RegisterSelector("m.UAntumbra", m.UAntumbra)
	transpiler.RegisterSelector("m.UEclipse", m.UEclipse)
	transpiler.RegisterSelector("m.UMono", m.UMono)
	transpiler.RegisterSelector("m.UPoly", m.UPoly)
	transpiler.RegisterSelector("m.UMega", m.UMega)
	transpiler.RegisterSelector("m.UQuad", m.UQuad)
	transpiler.RegisterSelector("m.UOct", m.UOct)
	transpiler.RegisterSelector("m.URisso", m.URisso)
	transpiler.RegisterSelector("m.UMinke", m.UMinke)
	transpiler.RegisterSelector("m.UBryde", m.UBryde)
	transpiler.RegisterSelector("m.USei", m.USei)
	transpiler.RegisterSelector("m.UOmura", m.UOmura)
	transpiler.RegisterSelector("m.UPlayerAlpha", m.UPlayerAlpha)
	transpiler.RegisterSelector("m.UPlayerBeta", m.UPlayerBeta)
	transpiler.RegisterSelector("m.UPlayerGamma", m.UPlayerGamma)
}

func registerNativeVariables() {
	transpiler.RegisterSelector("m.This", "@this")
	transpiler.RegisterSelector("m.ThisX", "@thisx")
	transpiler.RegisterSelector("m.ThisXf", "@thisx")
	transpiler.RegisterSelector("m.ThisY", "@thisy")
	transpiler.RegisterSelector("m.ThisYf", "@thisy")
	transpiler.RegisterSelector("m.Ipt", "@ipt")
	transpiler.RegisterSelector("m.Counter", "@counter")
	transpiler.RegisterSelector("m.Links", "@links")
	transpiler.RegisterSelector("m.CurUnit", "@unit")
	transpiler.RegisterSelector("m.Time", "@time")
	transpiler.RegisterSelector("m.Tick", "@tick")
	transpiler.RegisterSelector("m.MapW", "@mapw")
	transpiler.RegisterSelector("m.MapH", "@maph")
}

func registerConstants() {
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

	transpiler.RegisterSelector("m.ctrlProcessor", strconv.Itoa(m.CtrlProcessor))
	transpiler.RegisterSelector("m.ctrlFormation", strconv.Itoa(m.CtrlFormation))
	transpiler.RegisterSelector("m.ctrlPlayer", strconv.Itoa(m.CtrlPlayer))
	transpiler.RegisterSelector("m.ctrlSelf", strconv.Itoa(m.CtrlSelf))
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

const (
	TRUE  = "true"
	FALSE = "false"
)
