package m

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterFuncTranslation("m.Floor", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "op"},
							&transpiler.Value{Value: "floor"},
							&transpiler.Value{Value: transpiler.FunctionReturnVariable},
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.Random", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "op"},
							&transpiler.Value{Value: "rand"},
							&transpiler.Value{Value: transpiler.FunctionReturnVariable},
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}
		},
	})
}

// TODO Operations

// TODO Docs
func Floor(number float64) int {
	return 0
}

// TODO Docs
func Random(max float64) float64 {
	return 0
}
