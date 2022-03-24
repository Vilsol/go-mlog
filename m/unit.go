package m

// Load the next cached unit of the provided type into memory
//
// Will loop over once it reaches the end of the cache
func UnitBind(unitType UnitType) {
}

// Like Radar but originates from the cached units
//
// Retrieve a list of units that match specified conditions
//
// Conditions are combined using an `and` operation
func UnitRadar(target1 RadarTarget, target2 RadarTarget, target3 RadarTarget, sortOrder bool, sort RadarSort) Unit {
	return nil
}

// Locate a block of the provided ore type
//
// Also locates blocks outside the range of the unit
func UnitLocateOre(ore string) (x int, y int, found bool) {
	return 0, 0, false
}

// Locate a building of the provided type
//
// If enemy is true, derelict blocks cannot be located
//
// Also locates blocks outside the range of the unit
func UnitLocateBuilding(buildingType BlockFlag, enemy bool) (x int, y int, found bool, building UnspecifiedBuilding) {
	return 0, 0, false, nil
}

// Locate the enemy spawn
//
// Also locates blocks outside the range of the unit
func UnitLocateSpawn() (x int, y int, found bool, building Building) {
	return 0, 0, false, nil
}

// Locate a damaged building
//
// Also locates blocks outside the range of the unit
func UnitLocateDamaged() (x int, y int, found bool, building UnspecifiedBuilding) {
	return 0, 0, false, nil
}
