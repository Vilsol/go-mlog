package impl

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterFuncTranslation("m.UnitStop", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "stop"}, 0, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitMove", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "move"}, 2, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitApproach", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "approach"}, 3, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitBoost", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "boost"}, 1, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitPathfind", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "pathfind"}, 0, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitTarget", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "target"}, 3, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitTargetP", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "targetp"}, 2, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitItemDrop", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "itemDrop"}, 2, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitItemTake", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "itemTake"}, 3, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitPayloadDrop", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "payDrop"}, 0, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitPayloadTake", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "payTake"}, 1, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitMine", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "mine"}, 2, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitFlag", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "flag"}, 1, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitBuild", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: genBasicFuncTranslation([]string{"ucontrol", "build"}, 5, 0),
	})
	transpiler.RegisterFuncTranslation("m.UnitGetBlock", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 2,
		Translate: genBasicFuncTranslation([]string{"ucontrol", "getBlock"}, 2, 2),
	})
	transpiler.RegisterFuncTranslation("m.UnitWithin", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: genBasicFuncTranslation([]string{"ucontrol", "within"}, 3, 1),
	})
}
