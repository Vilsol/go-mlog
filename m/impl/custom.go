package impl

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
)

func init() {
	// TODO Optimize
	transpiler.RegisterFuncTranslation("m.Const", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "set"},
							vars[0],
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
						},
					},
				},
			}, nil
		},
	})
}
