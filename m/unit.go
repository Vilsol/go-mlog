package m

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
)

func init() {
	transpiler.RegisterFuncTranslation("m.UnitBind", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ubind"},
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitRadar", transpiler.Translator{
		Count:     1,
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "uradar"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
							&transpiler.Value{Value: "turret1"}, // Remove once fixed in game
							&transpiler.Value{Value: args[3].GetValue()},
							vars[0],
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitLocateOre", transpiler.Translator{
		Count:     1,
		Variables: 3,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ulocate"},
							&transpiler.Value{Value: "ore"},
							&transpiler.Value{Value: "core"}, // Remove once fixed in game
							&transpiler.Value{Value: "true"}, // Remove once fixed in game
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
							vars[0],
							vars[1],
							vars[2],
							&transpiler.Value{Value: "null"}, // Remove once fixed in game
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitLocateBuilding", transpiler.Translator{
		Count:     1,
		Variables: 4,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ulocate"},
							&transpiler.Value{Value: "building"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: "@copper"}, // Remove once fixed in game
							vars[0],
							vars[1],
							vars[2],
							vars[3],
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitLocateSpawn", transpiler.Translator{
		Count:     1,
		Variables: 4,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ulocate"},
							&transpiler.Value{Value: "spawn"},
							&transpiler.Value{Value: "core"},    // Remove once fixed in game
							&transpiler.Value{Value: "true"},    // Remove once fixed in game
							&transpiler.Value{Value: "@copper"}, // Remove once fixed in game
							vars[0],
							vars[1],
							vars[2],
							vars[3],
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitLocateDamaged", transpiler.Translator{
		Count:     1,
		Variables: 4,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ulocate"},
							&transpiler.Value{Value: "damaged"},
							&transpiler.Value{Value: "core"},    // Remove once fixed in game
							&transpiler.Value{Value: "true"},    // Remove once fixed in game
							&transpiler.Value{Value: "@copper"}, // Remove once fixed in game
							vars[0],
							vars[1],
							vars[2],
							vars[3],
						},
					},
				},
			}, nil
		},
	})
}

// Load the next cached unit of the provided type into memory
//
// Will loop over once it reaches the end of the cache
func UnitBind(unitType string) {
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
func UnitLocateBuilding(buildingType BlockFlag, enemy bool) (x int, y int, found bool, building Building) {
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
func UnitLocateDamaged() (x int, y int, found bool, building Building) {
	return 0, 0, false, nil
}
