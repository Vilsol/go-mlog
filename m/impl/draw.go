package impl

import (
	"errors"
	"github.com/Vilsol/go-mlog/decompiler"
	"github.com/Vilsol/go-mlog/transpiler"
	"go/ast"
	"strings"
)

func init() {
	initDrawTranspiler()
	initDrawDecompiler()
}

func initDrawTranspiler() {
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

func initDrawDecompiler() {
	decompiler.RegisterFuncTranslation("draw", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) < 2 {
				return nil, nil, errors.New("expecting at least 2 arguments")
			}

			var result []ast.Stmt

			switch args[0] {
			case "clear":
				if len(args) < 4 {
					return nil, nil, errors.New("expecting at least 4 arguments")
				}

				for i := 1; i < 4; i++ {
					if err := global.Assert(args[i], "int"); err != nil {
						return nil, nil, err
					}
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("DrawClear"),
							},
							Args: []ast.Expr{
								ast.NewIdent(args[1]),
								ast.NewIdent(args[2]),
								ast.NewIdent(args[3]),
							},
						},
					},
				}
			case "color":
				if len(args) < 5 {
					return nil, nil, errors.New("expecting at least 5 arguments")
				}

				for i := 1; i < 5; i++ {
					if err := global.Assert(args[i], "int"); err != nil {
						return nil, nil, err
					}
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("DrawColor"),
							},
							Args: []ast.Expr{
								ast.NewIdent(args[1]),
								ast.NewIdent(args[2]),
								ast.NewIdent(args[3]),
								ast.NewIdent(args[4]),
							},
						},
					},
				}
			case "stroke":
				if err := global.Assert(args[1], "int"); err != nil {
					return nil, nil, err
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("DrawStroke"),
							},
							Args: []ast.Expr{
								ast.NewIdent(args[1]),
							},
						},
					},
				}
			case "line":
				if len(args) < 5 {
					return nil, nil, errors.New("expecting at least 5 arguments")
				}

				for i := 1; i < 5; i++ {
					if err := global.Assert(args[i], "int"); err != nil {
						return nil, nil, err
					}
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("DrawLine"),
							},
							Args: []ast.Expr{
								ast.NewIdent(args[1]),
								ast.NewIdent(args[2]),
								ast.NewIdent(args[3]),
								ast.NewIdent(args[4]),
							},
						},
					},
				}
			case "rect":
				fallthrough
			case "lineRect":
				if len(args) < 5 {
					return nil, nil, errors.New("expecting at least 5 arguments")
				}

				for i := 1; i < 5; i++ {
					if err := global.Assert(args[i], "int"); err != nil {
						return nil, nil, err
					}
				}

				function := "DrawRect"
				if args[0] == "linePoly" {
					function = "DrawLineRect"
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent(function),
							},
							Args: []ast.Expr{
								ast.NewIdent(args[1]),
								ast.NewIdent(args[2]),
								ast.NewIdent(args[3]),
								ast.NewIdent(args[4]),
							},
						},
					},
				}
			case "poly":
				fallthrough
			case "linePoly":
				if len(args) < 6 {
					return nil, nil, errors.New("expecting at least 6 arguments")
				}

				for i := 1; i < 4; i++ {
					if err := global.Assert(args[i], "int"); err != nil {
						return nil, nil, err
					}
				}

				for i := 4; i < 6; i++ {
					if err := global.Assert(args[i], "float64"); err != nil {
						return nil, nil, err
					}
				}

				function := "DrawPoly"
				if args[0] == "linePoly" {
					function = "DrawLinePoly"
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent(function),
							},
							Args: []ast.Expr{
								ast.NewIdent(args[1]),
								ast.NewIdent(args[2]),
								ast.NewIdent(args[3]),
								ast.NewIdent(args[4]),
								ast.NewIdent(args[5]),
							},
						},
					},
				}
			case "triangle":
				if len(args) < 7 {
					return nil, nil, errors.New("expecting at least 7 arguments")
				}

				for i := 1; i < 7; i++ {
					if err := global.Assert(args[i], "int"); err != nil {
						return nil, nil, err
					}
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("DrawTriangle"),
							},
							Args: []ast.Expr{
								ast.NewIdent(args[1]),
								ast.NewIdent(args[2]),
								ast.NewIdent(args[3]),
								ast.NewIdent(args[4]),
								ast.NewIdent(args[5]),
								ast.NewIdent(args[6]),
							},
						},
					},
				}
			case "image":
				if len(args) < 6 {
					return nil, nil, errors.New("expecting at least 6 arguments")
				}

				for i := 1; i < 3; i++ {
					if err := global.Assert(args[i], "int"); err != nil {
						return nil, nil, err
					}
				}

				image, err := global.Resolve(args[3], "string")
				if err != nil {
					return nil, nil, err
				}

				for i := 4; i < 6; i++ {
					if err := global.Assert(args[i], "float64"); err != nil {
						return nil, nil, err
					}
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("DrawImage"),
							},
							Args: []ast.Expr{
								ast.NewIdent(args[1]),
								ast.NewIdent(args[2]),
								image,
								ast.NewIdent(args[4]),
								ast.NewIdent(args[5]),
							},
						},
					},
				}
			default:
				return nil, nil, errors.New("unknown control command: " + args[0])
			}

			return result, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
	decompiler.RegisterFuncTranslation("drawflush", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) != 1 {
				return nil, nil, errors.New("expecting 1 argument")
			}

			return []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("m"),
							Sel: ast.NewIdent("DrawFlush"),
						},
						Args: []ast.Expr{
							ast.NewIdent("\"" + args[0] + "\""),
						},
					},
				},
			}, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
}
