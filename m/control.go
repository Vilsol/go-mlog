package m

// Enable/Disable an block e.g. conveyor, door, switch
func ControlEnabled(target string, enabled bool) {
}

// Shoot with the provided turret at the target absolute position
//
// If shoot parameter is false, it will cease firing
func ControlShoot(turret string, x int, y int, shoot bool) {
}

// Smart version of ControlShoot
//
// Shoot with the provided turret at the predicted position of target unit
//
// If shoot parameter is false, it will cease firing
func ControlShootP(turret string, target int, shoot bool) {
}

// Set the configuration of the target building
func ControlConfigure(target string, configuration int) {
}
