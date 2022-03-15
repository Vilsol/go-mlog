package m

var (
	// A Building Object that represents the processor itself.
	// You can use this with sensor to find various properties about the processor.
	This = Building(nil)

	// The x coordinate of the processor.
	ThisX = 0

	// The y coordinate of the processor.
	ThisY = 0

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
	Tick = 0

	// Width of the map, in tiles.
	MapW = 0

	// Height of the map, in tiles.
	MapH = 0
)
