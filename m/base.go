package m

// Read a float64 value from memory at specified position
func Read(memory string, position int) int {
	return 0
}

// Write a float64 value to memory at specified position
//
// For integer equivalent use WriteInt
func Write(value int, memory string, position int) {
}

// Write an integer value to memory at specified position
//
// For float64 equivalent use Write
func WriteInt(value int, memory string, position int) {
}

// Flush all printed statements to the provided message block
func PrintFlush(targetMessage string) {
}

// Get the linked tile at the specified address
func GetLink(address int) Link {
	return nil
}

// Retrieve a list of units that match specified conditions
//
// Conditions are combined using an `and` operation
func Radar(from Ranged, target1 RadarTarget, target2 RadarTarget, target3 RadarTarget, sortOrder bool, sort RadarSort) Unit {
	return nil
}

// Extract information indicated by sense from the provided block
func Sensor(block HealthC, sense string) float64 {
	return 0
}
