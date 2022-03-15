package m

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

type RadarSort = string

const (
	RSDistance  = RadarSort("distance")
	RSHealth    = RadarSort("health")
	RSShield    = RadarSort("shield")
	RSArmor     = RadarSort("armor")
	RSMaxHealth = RadarSort("maxHealth")
)

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

type UnitType = string

const (
	UDagger   = UnitType("@dagger")
	UMace     = UnitType("@mace")
	UFortress = UnitType("@fortress")
	UScepter  = UnitType("@scepter")
	UReign    = UnitType("@reign")

	UNova   = UnitType("@nova")
	UPulsar = UnitType("@pulsar")
	UQuasar = UnitType("@quasar")
	UVela   = UnitType("@vela")
	UCorvus = UnitType("@corvus")

	UCrawler = UnitType("@crawler")
	UAtrax   = UnitType("@atrax")
	USpiroct = UnitType("@spiroct")
	UArkyid  = UnitType("@arkyid")
	UToxopid = UnitType("@toxopid")

	UFlare    = UnitType("@flare")
	UHorizon  = UnitType("@horizon")
	UZenith   = UnitType("@zenith")
	UAntumbra = UnitType("@antumbra")
	UEclipse  = UnitType("@eclipse")

	UMono = UnitType("@mono")
	UPoly = UnitType("@poly")
	UMega = UnitType("@mega")
	UQuad = UnitType("@quad")
	UOct  = UnitType("@oct")

	URisso = UnitType("@risso")
	UMinke = UnitType("@minke")
	UBryde = UnitType("@bryde")
	USei   = UnitType("@sei")
	UOmura = UnitType("@omura")

	UPlayerAlpha = UnitType("@alpha")
	UPlayerBeta  = UnitType("@beta")
	UPlayerGamma = UnitType("@gamma")
)

type ItemType = string

const (
	ITCopper        = ItemType("@copper")
	ITLead          = ItemType("@lead")
	ITMetaGlass     = ItemType("@metaglass")
	ITGraphite      = ItemType("@graphite")
	ITSand          = ItemType("@sand")
	ITCoal          = ItemType("@coal")
	ITTitanium      = ItemType("@titanium")
	ITThorium       = ItemType("@thorium")
	ITScrap         = ItemType("@scrap")
	ITSilicon       = ItemType("@silicon")
	ITPlastanium    = ItemType("@plastanium")
	ITPhaseFabric   = ItemType("@phase-fab")
	ITSurgeAlloy    = ItemType("@surge-alloy")
	ITSporePod      = ItemType("@spore-pod")
	ITBlastCompound = ItemType("@blast-compound")
	ITPyratite      = ItemType("@pyratite")
	// ITNone Used for checks where no item is found
	ITNone = ItemType("")
)

type FluidType = string

const (
	FWater     = FluidType("@water")
	FOil       = FluidType("@oil")
	FSlag      = FluidType("@slag")
	FCryofluid = FluidType("@cryofluid")
)

type Controller = int

const (
	CtrlProcessor = Controller(1)
	//TODO: doc
	CtrlFormation = Controller(2)
	CtrlPlayer    = Controller(3)
	// CtrlSelf The unit is idle or doing its own stuff -> not controlled
	CtrlSelf = Controller(0)
)
