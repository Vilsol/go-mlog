package impl

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterFuncTranslation("m.ControlEnabled", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"control", "enabled"}, 2, 0),
	})
	transpiler.RegisterFuncTranslation("m.ControlShoot", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"control", "shoot"}, 4, 0),
	})
	transpiler.RegisterFuncTranslation("m.ControlShootP", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"control", "shootp"}, 3, 0),
	})
	transpiler.RegisterFuncTranslation("m.ControlConfigure", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"control", "configure"}, 2, 0),
	})
}
