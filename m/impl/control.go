package impl

import (
	"errors"
	"github.com/Vilsol/go-mlog/decompiler"
	"github.com/Vilsol/go-mlog/transpiler"
	"go/ast"
)

func init() {
	initControlTranspiler()
	initControlDecompiler()
}

func initControlTranspiler() {
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

func initControlDecompiler() {
	decompiler.RegisterFuncTranslation("control", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) < 3 {
				return nil, nil, errors.New("expecting at least 3 arguments")
			}

			source, err := global.Resolve(args[1], "github.com/Vilsol/go-mlog/m.Building")
			if err != nil {
				return nil, nil, err
			}

			var result []ast.Stmt

			switch args[0] {
			case "enabled":
				enabled := args[2]
				if enabled != TRUE && enabled != FALSE && enabled != "1" && enabled != "0" {
					if err := global.Assert(enabled, "bool"); err != nil {
						return nil, nil, err
					}
				} else if enabled == "1" {
					enabled = TRUE
				} else if enabled == "0" {
					enabled = FALSE
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("ControlEnabled"),
							},
							Args: []ast.Expr{
								source,
								ast.NewIdent(enabled),
							},
						},
					},
				}
			case "shoot":
				if len(args) < 5 {
					return nil, nil, errors.New("expecting at least 5 arguments")
				}

				shoot := args[4]
				if shoot != TRUE && shoot != FALSE && shoot != "1" && shoot != "0" {
					if err := global.Assert(shoot, "bool"); err != nil {
						return nil, nil, err
					}
				} else if shoot == "1" {
					shoot = TRUE
				} else if shoot == "0" {
					shoot = FALSE
				}

				if err := global.Assert(args[2], "int"); err != nil {
					return nil, nil, err
				}

				if err := global.Assert(args[3], "int"); err != nil {
					return nil, nil, err
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("ControlShoot"),
							},
							Args: []ast.Expr{
								source,
								ast.NewIdent(args[2]),
								ast.NewIdent(args[3]),
								ast.NewIdent(shoot),
							},
						},
					},
				}
			case "shootp":
				if len(args) < 4 {
					return nil, nil, errors.New("expecting at least 4 arguments")
				}

				shoot := args[3]
				if shoot != "true" && shoot != "false" && shoot != "1" && shoot != "0" {
					if err := global.Assert(shoot, "bool"); err != nil {
						return nil, nil, err
					}
				} else if shoot == "1" {
					shoot = "true"
				} else if shoot == "0" {
					shoot = "false"
				}

				target, err := global.Resolve(args[2], "github.com/Vilsol/go-mlog/m.HealthC")
				if err != nil {
					return nil, nil, err
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("ControlShootP"),
							},
							Args: []ast.Expr{
								source,
								target,
								ast.NewIdent(shoot),
							},
						},
					},
				}
			case "config":
				if err := global.Assert(args[2], "int"); err != nil {
					return nil, nil, err
				}

				result = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("ControlConfigure"),
							},
							Args: []ast.Expr{
								source,
								ast.NewIdent(args[2]),
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
}
