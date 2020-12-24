package m

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterFuncTranslation("m.UnitBind", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ubind"},
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitRadar", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "radar"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
							&transpiler.Value{Value: "turret1"}, // Remove once fixed in game
							&transpiler.Value{Value: args[3].GetValue()},
							&transpiler.Value{Value: transpiler.FunctionReturnVariable},
						},
					},
				},
			}
		},
	})
	// TODO UnitLocateOre
	// ulocate ore storage A @copper B C D E

	// TODO UnitLocateBuilding
	// ulocate building storage A @copper B C D E

	// TODO UnitLocateSpawn
	// ulocate spawn core A @copper B C D E

	// TODO UnitLocateDamaged
	// ulocate damaged core A @copper B C D E
}

// TODO Docs
func UnitBind(unitType string) {
}

// TODO Docs
func UnitRadar(target1 RadarTarget, target2 RadarTarget, target3 RadarTarget, sortOrder int, sort RadarSort) Unit {
	return nil
}

// TODO Docs
func UnitLocateOre(ore string) (x int, y int, found int) {
	return 0, 0, 0
}

// TODO Docs
func UnitLocateBuilding(buildingType BlockFlag, enemy int) (x int, y int, found int, building Building) {
	return 0, 0, 0, nil
}

// TODO Docs
func UnitLocateSpawn() (x int, y int, found int, building Building) {
	return 0, 0, 0, nil
}

// TODO Docs
func UnitLocateDamaged() (x int, y int, found int, building Building) {
	return 0, 0, 0, nil
}
