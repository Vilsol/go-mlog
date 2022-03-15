package m

// Read a float64 value from memory at specified position
func Read[A integer](memory string, position A) int {
	return 0
}

// Write a value to memory at specified position
func Write[A float, B integer](value A, memory string, position B) {
}

// Flush all printed statements to the provided message block
func PrintFlush(targetMessage string) {
}

// Get the linked tile at the specified address
func GetLink[A integer](address A) UnspecifiedBuilding {
	return nil
}

// Retrieve a list of units that match specified conditions
//
// Conditions are combined using an `and` operation
func Radar(from Building, target1 RadarTarget, target2 RadarTarget, target3 RadarTarget, sortOrder bool, sort RadarSort) Unit {
	return nil
}

// Extract information indicated by sense from the provided block.
//Use this only if the needed information is not available using The getters of the building itself
//main purpose is to use generic "senses" or use things that are not mapped
func Sensor(block UnspecifiedBuilding, sense string) float64 {
	return 0
}

// String equivalent of Sensor
func SensorStr(block UnspecifiedBuilding, sense string) string {
	return ""
}
