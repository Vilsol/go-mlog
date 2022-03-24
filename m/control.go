package m

// Enable/Disable a block e.g. conveyor, door, switch
func ControlEnabled(target Building, enabled bool) {
}

// Shoot with the provided turret at the target absolute position
//
// If shoot parameter is false, it will cease firing
func ControlShoot[A integer, B integer](turret HasAmmo, x A, y B, shoot bool) {
}

// Smart version of ControlShoot
//
// Shoot with the provided turret at the predicted position of target unit
//
// If shoot parameter is false, it will cease firing
func ControlShootP(turret HasAmmo, target Shootable, shoot bool) {
}

// Set the configuration of the target building
func ControlConfigure[A integer](target HasConfig, configuration A) {
}
