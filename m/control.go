package m

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterFuncTranslation("m.ControlEnabled", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
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
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.ControlShoot", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
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
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.ControlShootP", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
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
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.ControlConfigure", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
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
			}
		},
	})
}

// TODO Docs
func ControlEnabled(target string, enabled int) {
}

// TODO Docs
func ControlShoot(target string, x int, y int, shoot int) {
}

// TODO Docs
func ControlShootP(target string, unit int, shoot int) {
}

// TODO Docs
func ControlConfigure(target string, configuration int) {
}
