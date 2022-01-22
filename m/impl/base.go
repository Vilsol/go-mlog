package impl

import (
	"errors"
	"github.com/Vilsol/go-mlog/decompiler"
	"github.com/Vilsol/go-mlog/transpiler"
	"go/ast"
	"strings"
)

func init() {
	initBaseTranspiler()
	initBaseDecompiler()
}

func initBaseTranspiler() {
	transpiler.RegisterFuncTranslation("m.Read", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			memoryName := strings.Trim(args[0].GetValue(), "\"")

			// TODO Remove hardcode
			if memoryName == "bank1" {
				return nil, errors.New("can't read/write to memory cell that is used for the stack: bank1")
			}

			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "read"},
							vars[0],
							&transpiler.Value{Value: memoryName},
							&transpiler.Value{Value: args[1].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	write := transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			memoryName := strings.Trim(args[1].GetValue(), "\"")

			if memoryName == "bank1" {
				return nil, errors.New("can't read/write to memory cell that is used for the stack: bank1")
			}

			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "write"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: memoryName},
							&transpiler.Value{Value: args[2].GetValue()},
						},
					},
				},
			}, nil
		},
	}
	transpiler.RegisterFuncTranslation("m.Write", write)
	transpiler.RegisterFuncTranslation("m.WriteInt", write)
	transpiler.RegisterFuncTranslation("m.PrintFlush", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "printflush"},
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.GetLink", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "getlink"},
							vars[0],
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.Radar", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "radar"},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[3].GetValue()},
							&transpiler.Value{Value: args[5].GetValue()},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
							vars[0],
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.Sensor", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "sensor"},
							vars[0],
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
							&transpiler.Value{Value: strings.Trim(args[1].GetValue(), "\"")},
						},
					},
				},
			}, nil
		},
	})
}

func initBaseDecompiler() {
	decompiler.RegisterFuncTranslation("read", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) != 3 {
				return nil, nil, errors.New("expecting 3 arguments")
			}

			tok, err := global.AssignOrDefine(args[0], "int")
			if err != nil {
				return nil, nil, err
			}

			return []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(args[0]),
					},
					Tok: tok,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("Read"),
							},
							Args: []ast.Expr{
								ast.NewIdent("\"" + args[1] + "\""),
								ast.NewIdent(args[2]),
							},
						},
					},
				},
			}, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
	decompiler.RegisterFuncTranslation("write", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) != 3 {
				return nil, nil, errors.New("expecting 3 arguments")
			}

			calledFunc := "Write"

			variableType, err := global.GetType(args[0])
			if err != nil {
				return nil, nil, err
			}

			if variableType == "int" {
				calledFunc = "WriteInt"
			}

			return []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("m"),
							Sel: ast.NewIdent(calledFunc),
						},
						Args: []ast.Expr{
							ast.NewIdent(args[0]),
							ast.NewIdent("\"" + args[1] + "\""),
							ast.NewIdent(args[2]),
						},
					},
				},
			}, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
	decompiler.RegisterFuncTranslation("printflush", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) != 1 {
				return nil, nil, errors.New("expecting 1 argument")
			}

			return []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("m"),
							Sel: ast.NewIdent("PrintFlush"),
						},
						Args: []ast.Expr{
							ast.NewIdent("\"" + args[0] + "\""),
						},
					},
				},
			}, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
	decompiler.RegisterFuncTranslation("getlink", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) != 2 {
				return nil, nil, errors.New("expecting 2 arguments")
			}

			tok, err := global.AssignOrDefine(args[0], "github.com/Vilsol/go-mlog/m.Link")
			if err != nil {
				return nil, nil, err
			}

			return []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(args[0]),
					},
					Tok: tok,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("GetLink"),
							},
							Args: []ast.Expr{
								ast.NewIdent(args[1]),
							},
						},
					},
				},
			}, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
	decompiler.RegisterFuncTranslation("radar", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) != 7 {
				return nil, nil, errors.New("expecting 7 arguments")
			}

			target1 := "m.RT" + strings.ToUpper(args[0][:1]) + args[0][1:]
			target2 := "m.RT" + strings.ToUpper(args[1][:1]) + args[1][1:]
			target3 := "m.RT" + strings.ToUpper(args[2][:1]) + args[2][1:]
			sort := "m.RS" + strings.ToUpper(args[3][:1]) + args[3][1:]

			source, err := global.Resolve(args[5], "github.com/Vilsol/go-mlog/m.Ranged")
			if err != nil {
				return nil, nil, err
			}

			tok, err := global.AssignOrDefine(args[6], "github.com/Vilsol/go-mlog/m.Unit")
			if err != nil {
				return nil, nil, err
			}

			return []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(args[6]),
					},
					Tok: tok,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("Radar"),
							},
							Args: []ast.Expr{
								source,
								ast.NewIdent(target1),
								ast.NewIdent(target2),
								ast.NewIdent(target3),
								ast.NewIdent(args[5]),
								ast.NewIdent(sort),
							},
						},
					},
				},
			}, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
	decompiler.RegisterFuncTranslation("sensor", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) != 3 {
				return nil, nil, errors.New("expecting 3 arguments")
			}

			tok, err := global.AssignOrDefine(args[0], "float64")
			if err != nil {
				return nil, nil, err
			}

			source, err := global.Resolve(args[1], "github.com/Vilsol/go-mlog/m.HealthC")
			if err != nil {
				return nil, nil, err
			}

			return []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(args[0]),
					},
					Tok: tok,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("Sensor"),
							},
							Args: []ast.Expr{
								source,
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("m"),
										Sel: ast.NewIdent("Const"),
									},
									Args: []ast.Expr{
										ast.NewIdent("\"" + args[2] + "\""),
									},
								},
							},
						},
					},
				},
			}, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
}
