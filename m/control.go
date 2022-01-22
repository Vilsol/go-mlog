package m

// Enable/Disable a block e.g. conveyor, door, switch
func ControlEnabled(target Building, enabled bool) {
}

// Shoot with the provided turret at the target absolute position
//
// If shoot parameter is false, it will cease firing
func ControlShoot(turret Building, x int, y int, shoot bool) {
}

// Smart version of ControlShoot
//
// Shoot with the provided turret at the predicted position of target unit
//
// If shoot parameter is false, it will cease firing
func ControlShootP(turret Building, target HealthC, shoot bool) {
}

// Set the configuration of the target building
func ControlConfigure(target Building, configuration int) {
}
