package x

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterFuncTranslation("x.Sleep", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 2
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			dVar := &transpiler.DynamicVariable{}

			jump := &transpiler.MLOGJump{
				MLOG: transpiler.MLOG{},
				Condition: []transpiler.Resolvable{
					&transpiler.Value{Value: "lessThan"},
					&transpiler.Value{Value: "@time"},
					dVar,
				},
			}

			jump.JumpTarget = &transpiler.StatementJumpTarget{
				Statement: jump,
				After:     false,
			}

			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "op"},
							&transpiler.Value{Value: "add"},
							dVar,
							&transpiler.Value{Value: "@time"},
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
				jump,
			}, nil
		},
	})
}

// Sleep for n milliseconds
func Sleep(millis int) {
}
