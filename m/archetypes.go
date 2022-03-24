package m

type Turret interface {
	Building
	HasAmmo
	HasInventory
	Ranged
	HasLiquid
	HasPower
	Controllable
	// GetRotation the rotation at a scale of 0f..360f. 0 and 360 deg is when the turret looks at the very right side
	GetRotation() int
	// GetShootX the x cord of the targeted tile
	GetShootX() int
	// GetShootY the y cord of the targeted tile
	GetShootY() int
	// GetShootPosition the position of the targeted tile
	GetShootPosition() HasPosition
}

type Building interface {
	HasPosition
	Healthc
	Teamc
	GameElement
	// GetType the type of the block (e.g. "cyclone" "microprocessor") not compatible with BlockFlag
	GetType() string
	IsEnabled() bool
}
type Unit interface {
	HasPosition
	Healthc
	Ranged
	Controllable
	Teamc
	GameElement
	HasInventory
	GetType() UnitType
	// IsDead true if unit is not alive or invalid. PS: Why anuken why "dead" and not alive?
	IsDead() bool
	IsBoosting() bool
	GetFlag() int
}
type MiningUnit interface {
	Unit
	CanMine
}
type CarrierUnit interface {
	Payloadc
	Unit
}
type UnspecifiedBuilding interface {
	Building
	HasPower
	HasAmmo
	HasInventory
	HasLiquid
	HasConfig
	Ranged
}
type GameElement interface{}

type Healthc interface {
	// GetHealth Get the current health value of Unit or Building
	GetHealth() int
	GetMaxHealth() int
}
type HasInventory interface {
	// GetTotalItems items in inventory
	GetTotalItems() int
	// GetFirstItem Type of the first item in the inventory
	GetFirstItem() ItemType
	// GetItemCapacity The maximum amount of items this unit/block can hold
	GetItemCapacity() int
}
type HasLiquid interface {
	// GetTotalLiquids liquids in internal storage
	GetTotalLiquids() int
	// GetLiquidCapacity The maximum amount of liquid this unit/block can hold
	GetLiquidCapacity() int
}
type Shootable interface {
	Healthc
	HasPosition
}
type HasPower interface {
	// GetTotalPower In case of unbuffered consumers, this is the percentage (1.0f = 100%) of the demanded power which can be supplied.
	// Blocks will work at a reduced efficiency if this is not equal to 1.0f.
	// In case of buffered consumers, this is storage capacity.
	GetTotalPower() int
	// GetPowerCapacity The maximum amount of power this unit/block can hold
	// In case of unbuffered consumers, this is the 0
	GetPowerCapacity() int
	// GetPowerNetStored power stored in the connected power network
	GetPowerNetStored() int64
	// GetPowerNetCapacity power the connected power network is able to store
	GetPowerNetCapacity() int64
	// GetPowerNetIn power the connected power network produces
	GetPowerNetIn() int64
	// GetPowerNetOut power the connected power network consumes
	GetPowerNetOut() int64
}
type HasAmmo interface {
	// GetAmmo The amount of ammo (in shots not units) this unit/turret holds
	GetAmmo() int
	// GetAmmoCapacity the amount of shots this turret/unit can store
	GetAmmoCapacity() int
	IsShooting() bool
}
type HasHeat interface {
	// GetHeat The heat this building has accumulated (in nuclear reactor for example)
	GetHeat()
}
type Productive interface {
	// GetEfficiency The effi ciency at which this item produces/consumes, scale of 0..1 efficiencies
	//above 100% (1) are not measurable. A Thermal generator with 300% efficiency
	//and one with 100% will both return 1
	//Note: Only seems to return the "electric efficiency"
	GetEfficiency() float64
}

//TODO: sensor timescale -> no idea

type HasPosition interface {
	GetX() int
	GetY() int
	// GetSize the side length of the block measured in tiles
	GetSize() int
}

type Controllable interface {
	GetControlled() Controller
	// GetController name of the thing that commands this thing. "processor" <playerName> or the unit itself
	//TODO: dont offer this method or find a more typesave way
	GetController() string
}
type Teamc interface {
	// GetTeam The team the thing belongs to
	GetTeam() int
}
type Ranged interface {
	// GetRange To shoot and or detect who knows (TODO: better docs)
	GetRange() float64
}
type CanMine interface {
	GetMineX() int
	GetMineY() int
	IsMining() bool
}
type HasName interface {
	// GetName the exact name of the blockType or unit type. Use GetType preferably
	GetName() string
}
type HasConfig interface {
	// GetConfig The filtered item
	GetConfig() string //mapped to "@configure" not "@config" nobody knows what latter one is
}

type Payloadc interface {
	// GetPayloadType block or unit type this unit is carrying
	GetPayloadType() string
	// GetPayloadCount How many things this unit is carrying (not items)
	GetPayloadCount() int
}
