package m

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

type HealthC = interface {
	GetHealth() int
	GetName() string
	GetX() float64
	GetY() float64

	GetTotalItems() int
	GetItemCapacity() int
	GetRotation() float64
	GetShootX() float64
	GetShootY() float64
	IsShooting() bool
}

type Ranged = interface{}

type Unit = interface {
	HealthC
	Ranged
}

type Building = interface {
	HealthC
	Ranged

	GetTotalLiquids() float64
	GetLiquidCapaticy() float64

	// In case of unbuffered consumers, this is the percentage (1.0f = 100%) of the demanded power which can be supplied.
	// Blocks will work at a reduced efficiency if this is not equal to 1.0f.
	// In case of buffered consumers, this is storage capacity.
	GetTotalPower() float64

	// In case of unbuffered consumers, this is the 0
	// n case of buffered consumers, this is the maximum storage capacity.
	GetPowerCapaticy() float64

	GetPowerNetStored() float64
	GetPowerNetCapacity() float64
	GetPowerNetIn() float64
	GetPowerNetOut() float64
	GetHeat() float64
	GetEfficiency() float64
	IsEnabled() bool
}

var (
	// A Building Object that represents the processor itself.
	// You can use this with sensor to find various properties about the processor.
	This = Building(nil)

	// The x coordinate of the processor.
	ThisX = 0
	// Convenience constant, same as float64(ThisX)
	ThisXf = float64(ThisX)

	// The y coordinate of the processor.
	ThisY = 0
	// Convenience constant, same as float64(ThisY)
	ThisYf = float64(ThisY)

	// The number of instructions executed per tick (60 ticks/second).
	//
	// Micro Processor -> 2
	// Logic Processor -> 8
	// Hyper Processor -> 25
	Ipt = 0

	// A variable that represents the next line the processor will read code from, equivalent to %IP in x86.
	// It can be changed like any other variable as another way to perform jumps.
	Counter = 0

	// A constant that equals the number of buildings linked to the processor.
	// It is changed by the processor when blocks are linked or unlinked.
	Links = 0

	// A constant that represents the current bound unit.
	// It only changes when the processor unbinds a unit, or binds another one.
	CurUnit = Unit(nil)

	// Represents the current UNIX timestamp in milliseconds.
	Time = 0

	// Represents the amount of ticks (60 ticks/second) since the map began.
	Tick = float64(0)

	// Width of the map, in tiles.
	MapW = 0

	// Height of the map, in tiles.
	MapH = 0
)

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
