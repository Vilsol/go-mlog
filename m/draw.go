package m

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
)

func init() {
	transpiler.RegisterFuncTranslation("m.DrawClear", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "clear"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
						},
					},
				},
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.DrawColor", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "color"},
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
	transpiler.RegisterFuncTranslation("m.DrawStroke", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "stroke"},
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.DrawLine", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "line"},
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
	transpiler.RegisterFuncTranslation("m.DrawRect", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "rect"},
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
	transpiler.RegisterFuncTranslation("m.DrawLineRect", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "lineRect"},
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
	transpiler.RegisterFuncTranslation("m.DrawPoly", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "poly"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[3].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
						},
					},
				},
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.DrawLinePoly", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "linePoly"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[3].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
						},
					},
				},
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.DrawTriangle", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "triangle"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[3].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
							&transpiler.Value{Value: args[5].GetValue()},
						},
					},
				},
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.DrawImage", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "draw"},
							&transpiler.Value{Value: "image"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[3].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
						},
					},
				},
			}
		},
	})
	transpiler.RegisterFuncTranslation("m.DrawFlush", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable) []transpiler.MLOGStatement {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "drawflush"},
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
						},
					},
				},
			}
		},
	})
}

// TODO Docs
func DrawClear(r int, g int, b int) {
}

// TODO Docs
func DrawColor(r int, g int, b int, a int) {
}

// TODO Docs
func DrawStroke(width int) {
}

// TODO Docs
func DrawLine(x1 int, y1 int, x2 int, y2 int) {
}

// TODO Docs
func DrawRect(x int, y int, width int, height int) {
}

// TODO Docs
func DrawLineRect(x int, y int, width int, height int) {
}

// TODO Docs
func DrawPoly(x int, y int, sides int, radius float32, rotation float32) {
}

// TODO Docs
func DrawLinePoly(x int, y int, sides int, radius float32, rotation float32) {
}

// TODO Docs
func DrawTriangle(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) {
}

// TODO Docs
func DrawImage(x int, y int, image string, size float32, rotation float32) {
}

// TODO Docs
func DrawFlush(targetDisplay string) {
}
