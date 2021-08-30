package impl

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
)

func init() {
	transpiler.RegisterFuncTranslation("m.UnitBind", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ubind"},
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitRadar", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
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
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
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
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
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
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 4,
		// Cleanup arguments once fixed in game
		Translate: genBasicFuncTranslation([]string{"ulocate", "spawn", "core", "true", "@copper"}, 0, 4),
	})
	transpiler.RegisterFuncTranslation("m.UnitLocateDamaged", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 4,
		// Cleanup arguments once fixed in game
		Translate: genBasicFuncTranslation([]string{"ulocate", "damaged", "core", "true", "@copper"}, 0, 4),
	})
}
