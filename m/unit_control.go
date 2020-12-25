package m

import "github.com/Vilsol/go-mlog/transpiler"

func init() {
	transpiler.RegisterFuncTranslation("m.UnitStop", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "stop"},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitMove", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "move"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitApproach", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "approach"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitBoost", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "boost"},
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitPathfind", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "pathfind"},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitTarget", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "target"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitTargetP", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "targetp"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitItemDrop", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "itemDrop"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitItemTake", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "itemTake"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitPayloadDrop", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "payDrop"},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitPayloadTake", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "payTake"},
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitMine", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "mine"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitFlag", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "flag"},
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitBuild", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "build"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[3].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitGetBlock", transpiler.Translator{
		Count:     1,
		Variables: 2,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "getBlock"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							vars[0],
							vars[1],
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.UnitWithin", transpiler.Translator{
		Count:     1,
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "ucontrol"},
							&transpiler.Value{Value: "within"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							vars[0],
						},
					},
				},
			}, nil
		},
	})
}

// TODO Docs
func UnitStop() {
}

// TODO Docs
func UnitMove(x float64, y float64) {
}

// TODO Docs
func UnitApproach(x float64, y float64, radius float64) {
}

// TODO Docs
func UnitBoost(enable int) {
}

// TODO Docs
func UnitPathfind() {
}

// TODO Docs
func UnitTarget(x float64, y float64, shoot int) {
}

// TODO Docs
func UnitTargetP(target int, shoot int) {
}

// TODO Docs
func UnitItemDrop(to Building, amount int) {
}

// TODO Docs
func UnitItemTake(from Building, item string, amount int) {
}

// TODO Docs
func UnitPayloadDrop() {
}

// TODO Docs
func UnitPayloadTake(takeUnits int) {
}

// TODO Docs
func UnitMine(x float64, y float64) {
}

// TODO Docs
func UnitFlag(flag float64) {
}

// TODO Docs
func UnitBuild(x float64, y float64, block string, rotation int, config int) {
}

// TODO Docs
func UnitGetBlock(x float64, y float64) (blockType string, building Building) {
	return "", nil
}

// TODO Docs
func UnitWithin(x float64, y float64, radius float64) int {
	return 0
}
