package m

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterFuncTranslation("m.Floor", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "op"},
							&transpiler.Value{Value: "floor"},
							vars[0],
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.Random", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "op"},
							&transpiler.Value{Value: "rand"},
							vars[0],
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
}

// TODO Operations

// Floor the provided floating point number and convert to integer
func Floor(number float64) int {
	return 0
}

// Generate a random floating point number between 0 (inclusive) and max (exclusive)
func Random(max float64) float64 {
	return 0
}
