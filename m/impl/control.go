package impl

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterFuncTranslation("m.ControlEnabled", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "control"},
							&transpiler.Value{Value: "enabled"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.ControlShoot", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "control"},
							&transpiler.Value{Value: "shoot"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[3].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.ControlShootP", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "control"},
							&transpiler.Value{Value: "shootp"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.ControlConfigure", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "control"},
							&transpiler.Value{Value: "configure"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
						},
					},
				},
			}, nil
		},
	})
}
