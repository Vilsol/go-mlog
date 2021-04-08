package m

// Stop all actions including shooting
func UnitStop() {
}

// Move to the provided absolute position on the map
func UnitMove(x float64, y float64) {
}

// Approach a circular radius around the provided point
//
// Will stop moving once it is the provided radius away from the point
func UnitApproach(x float64, y float64, radius float64) {
}

// Enable/Disable boosting for mechs
func UnitBoost(enable bool) {
}

// Make the unit follow standard AI
//
// Find enemy cores, guard spawns, obey command centers
func UnitPathfind() {
}

// Like ControlShoot but for units
//
// Shoot with the cached unit at the target absolute position
//
// If shoot parameter is false, it will cease firing
//
// Will not shoot outside of the units range!
func UnitTarget(x float64, y float64, shoot bool) {
}

// Like ControlShootP but for units
//
// Shoot with the cached unit at the predicted position of target unit
//
// If shoot parameter is false, it will cease firing
func UnitTargetP(target HealthC, shoot bool) {
}

// Drops items into the provided building
//
// Will not drop more than provided amount
func UnitItemDrop(to Building, amount int) {
}

// Takes the provided item type from the provided building
//
// Will not take more than provided amount
func UnitItemTake(from Building, item string, amount int) {
}

// Drops the current payload
//
// Will only drop blocks if there is an empty space
func UnitPayloadDrop() {
}

// Pick up payload from underneath the unit
//
// If takeUnits is true, will also pick up units
func UnitPayloadTake(takeUnits bool) {
}

// Mine the ore at the specified absolute position
//
// Will not do anything if there is no minable ore or it is already being mined
func UnitMine(x float64, y float64) {
}

// Set the units flag
//
// Shown as a number when hovering over a unit
func UnitFlag(flag float64) {
}

// Build a block at the specified absolute position
func UnitBuild(x float64, y float64, block string, rotation int, config int) {
}

// Retrieve the building and its type at the specified absolute position
func UnitGetBlock(x float64, y float64) (blockType string, building Building) {
	return "", nil
}

// Checks whether there is a unit within the specified radius around the provided absolute position
func UnitWithin(x float64, y float64, radius float64) bool {
	return false
}
