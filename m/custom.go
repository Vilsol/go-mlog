package m

// Process provided string as a game constant
//
// Example m.Const("@copper") (only use this if you cannot use the enums like ItemType or UnitType
// use this if you have dynamic values like m.Const("@"+"fl"+"are") (whyever you want to do this)
func Const(constant string) string {
	return ""
}

// Return a building of the provided name
//
// Example m.B("message1")
func B(name string) UnspecifiedBuilding {
	return nil
}

// GetTurret Get a block that is specificaly a turret alias for B with better typing
func GetTurret(name string) Turret {
	return nil
}
