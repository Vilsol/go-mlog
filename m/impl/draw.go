package impl

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
)

func init() {
	transpiler.RegisterFuncTranslation("m.DrawClear", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"draw", "clear"}, 3, 0),
	})
	transpiler.RegisterFuncTranslation("m.DrawColor", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"draw", "color"}, 4, 0),
	})
	transpiler.RegisterFuncTranslation("m.DrawStroke", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"draw", "stroke"}, 1, 0),
	})
	transpiler.RegisterFuncTranslation("m.DrawLine", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"draw", "line"}, 4, 0),
	})
	transpiler.RegisterFuncTranslation("m.DrawRect", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"draw", "rect"}, 4, 0),
	})
	transpiler.RegisterFuncTranslation("m.DrawLineRect", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"draw", "lineRect"}, 4, 0),
	})
	transpiler.RegisterFuncTranslation("m.DrawPoly", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"draw", "poly"}, 5, 0),
	})
	transpiler.RegisterFuncTranslation("m.DrawLinePoly", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"draw", "linePoly"}, 5, 0),
	})
	transpiler.RegisterFuncTranslation("m.DrawTriangle", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"draw", "triangle"}, 6, 0),
	})
	transpiler.RegisterFuncTranslation("m.DrawImage", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "image"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: strings.Trim(args[2].GetValue(), "\"")},
							&transpiler.Value{Value: args[3].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.DrawFlush", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "drawflush"},
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
						},
					},
				},
			}, nil
		},
	})
}
