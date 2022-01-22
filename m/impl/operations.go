package impl

import (
	"errors"
	"github.com/Vilsol/go-mlog/decompiler"
	"github.com/Vilsol/go-mlog/transpiler"
	"go/ast"
	"go/token"
	"strconv"
)

func init() {
	initOperationTranspiler()
	initOperationDecompiler()
}

func initOperationTranspiler() {
	transpiler.RegisterFuncTranslation("m.Floor", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "op"},
							&transpiler.Value{Value: "floor"},
							vars[0],
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.Random", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "op"},
							&transpiler.Value{Value: "rand"},
							vars[0],
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.Log10", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "op"},
							&transpiler.Value{Value: "log10"},
							vars[0],
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.Ceil", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "op"},
							&transpiler.Value{Value: "ceil"},
							vars[0],
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	//op idiv result a b
}

type opDecompiler struct {
	ReturnType string
	ArgCount   int
	Token      token.Token
	Custom     func(args []string, global *decompiler.Global) ast.Expr
}

var opReturnTypes = map[string]opDecompiler{
	"add": {
		ReturnType: "float64",
		ArgCount:   4,
		Token:      token.ADD,
	},
	"sub": {
		ReturnType: "float64",
		ArgCount:   4,
		Token:      token.SUB,
	},
	"mul": {
		ReturnType: "float64",
		ArgCount:   4,
		Token:      token.MUL,
	},
	"div": {
		ReturnType: "float64",
		ArgCount:   4,
		Token:      token.QUO,
	},
	"idiv": {
		ReturnType: "int",
		ArgCount:   4,
		Custom: func(args []string, global *decompiler.Global) ast.Expr {
			return &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("m"),
					Sel: ast.NewIdent("IntDiv"),
				},
				Args: []ast.Expr{
					ast.NewIdent(args[2]),
					ast.NewIdent(args[3]),
				},
			}
		},
	},
	"mod": {
		ReturnType: "int",
		ArgCount:   4,
		Token:      token.REM,
	},
	"pow": {
		ReturnType: "float64",
		ArgCount:   4,
	},
	"equal": {
		ReturnType: "bool",
		ArgCount:   4,
		Token:      token.EQL,
	},
	"notEqual": {
		ReturnType: "bool",
		ArgCount:   4,
		Token:      token.NEQ,
	},
	"land": {
		ReturnType: "bool",
		ArgCount:   4,
		Token:      token.LAND,
	},
	"lessThan": {
		ReturnType: "bool",
		ArgCount:   4,
		Token:      token.LSS,
	},
	"lessThanEq": {
		ReturnType: "bool",
		ArgCount:   4,
		Token:      token.LEQ,
	},
	"greaterThan": {
		ReturnType: "bool",
		ArgCount:   4,
		Token:      token.GTR,
	},
	"greaterThanEq": {
		ReturnType: "bool",
		ArgCount:   4,
		Token:      token.GEQ,
	},
	"strictEqual": {
		ReturnType: "bool",
		ArgCount:   4,
		Token:      token.EQL,
	},
	"shl": {
		ReturnType: "int",
		ArgCount:   4,
		Token:      token.SHL,
	},
	"shr": {
		ReturnType: "int",
		ArgCount:   4,
		Token:      token.SHR,
	},
	"or": {
		ReturnType: "int",
		ArgCount:   4,
		Token:      token.OR,
	},
	"and": {
		ReturnType: "int",
		ArgCount:   4,
		Token:      token.AND,
	},
	"xor": {
		ReturnType: "int",
		ArgCount:   4,
		Token:      token.XOR,
	},
	"not": {
		ReturnType: "int",
		ArgCount:   3,
		Custom: func(args []string, global *decompiler.Global) ast.Expr {
			return &ast.BinaryExpr{
				X:  ast.NewIdent(args[2]),
				Op: token.XOR,
				Y:  ast.NewIdent("0xFFFFFFFFFFFFFFFF"),
			}
		},
	},
	"max": {
		ReturnType: "float64",
		ArgCount:   4,
	},
	"min": {
		ReturnType: "float64",
		ArgCount:   4,
	},
	"angle": {
		ReturnType: "float64",
		ArgCount:   4,
	},
	"len": {
		ReturnType: "float64",
		ArgCount:   4,
	},
	"noise": {
		ReturnType: "float64",
		ArgCount:   4,
	},
	"abs": {
		ReturnType: "float64",
		ArgCount:   3,
	},
	"log": {
		ReturnType: "float64",
		ArgCount:   3,
	},
	"log10": {
		ReturnType: "float64",
		ArgCount:   3,
		Custom: func(args []string, global *decompiler.Global) ast.Expr {
			return &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("m"),
					Sel: ast.NewIdent("Log10"),
				},
				Args: []ast.Expr{
					ast.NewIdent(args[2]),
				},
			}
		},
	},
	"floor": {
		ReturnType: "int",
		ArgCount:   3,
		Custom: func(args []string, global *decompiler.Global) ast.Expr {
			return &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("m"),
					Sel: ast.NewIdent("Floor"),
				},
				Args: []ast.Expr{
					ast.NewIdent(args[2]),
				},
			}
		},
	},
	"ceil": {
		ReturnType: "int",
		ArgCount:   3,
		Custom: func(args []string, global *decompiler.Global) ast.Expr {
			return &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("m"),
					Sel: ast.NewIdent("Ceil"),
				},
				Args: []ast.Expr{
					ast.NewIdent(args[2]),
				},
			}
		},
	},
	"sqrt": {
		ReturnType: "float64",
		ArgCount:   3,
	},
	"rand": {
		ReturnType: "float64",
		ArgCount:   3,
		Custom: func(args []string, global *decompiler.Global) ast.Expr {
			return &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("m"),
					Sel: ast.NewIdent("Random"),
				},
				Args: []ast.Expr{
					ast.NewIdent(args[2]),
				},
			}
		},
	},
	"sin": {
		ReturnType: "float64",
		ArgCount:   3,
	},
	"cos": {
		ReturnType: "float64",
		ArgCount:   3,
	},
	"tan": {
		ReturnType: "float64",
		ArgCount:   3,
	},
	"asin": {
		ReturnType: "float64",
		ArgCount:   3,
	},
	"acos": {
		ReturnType: "float64",
		ArgCount:   3,
	},
	"atan": {
		ReturnType: "float64",
		ArgCount:   3,
	},
}

func initOperationDecompiler() {
	decompiler.RegisterFuncTranslation("op", decompiler.Translator{
		Translate: func(args []string, global *decompiler.Global) ([]ast.Stmt, []string, error) {
			if len(args) < 2 {
				return nil, nil, errors.New("expecting at least 2 arguments")
			}

			dec, ok := opReturnTypes[args[0]]
			if !ok {
				return nil, nil, errors.New("unknown operation: " + args[0])
			}

			if len(args) < dec.ArgCount {
				return nil, nil, errors.New("expecting at least " + strconv.Itoa(dec.ArgCount) + " arguments")
			}

			tok, err := global.AssignOrDefine(args[1], dec.ReturnType)
			if err != nil {
				return nil, nil, err
			}

			var expr ast.Expr
			if dec.Custom != nil {
				expr = dec.Custom(args, global)
			} else if dec.Token != 0 {
				expr = &ast.BinaryExpr{
					X:  ast.NewIdent(args[2]),
					Op: dec.Token,
					Y:  ast.NewIdent(args[3]),
				}
			} else {
				return nil, nil, errors.New("unsupported operation: " + args[0])
			}

			return []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(args[1]),
					},
					Tok: tok,
					Rhs: []ast.Expr{
						expr,
					},
				},
			}, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
}
