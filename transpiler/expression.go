package transpiler

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go/ast"
	"go/token"
)

func expressionToMLOG(ident Resolvable, expr ast.Expr, options Options) ([]MLOGStatement, error) {
	switch expr.(type) {
	case *ast.BasicLit:
		basicExpr := expr.(*ast.BasicLit)
		if basicExpr.Kind != token.INT && basicExpr.Kind != token.FLOAT {
			return nil, errors.New("only integers and floating point numbers can be assigned to variables")
		}
		return []MLOGStatement{&MLOG{
			Comment: "Set the variable to the value",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					ident,
					&Value{Value: basicExpr.Value},
				},
			},
		}}, nil
	case *ast.Ident:
		identExpr := expr.(*ast.Ident)
		return []MLOGStatement{&MLOG{
			Comment: "Set the variable to the value of other variable",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					ident,
					&NormalVariable{Name: identExpr.Name},
				},
			},
		}}, nil
	case *ast.BinaryExpr:
		binaryExpr := expr.(*ast.BinaryExpr)

		if opTranslated, ok := regularOperators[binaryExpr.Op]; ok {
			instructions := make([]MLOGStatement, 0)
			var leftSide Resolvable
			var rightSide Resolvable

			if basicLit, ok := binaryExpr.X.(*ast.BasicLit); ok {
				leftSide = &Value{Value: basicLit.Value}
			} else if leftIdent, ok := binaryExpr.X.(*ast.Ident); ok {
				leftSide = &NormalVariable{Name: leftIdent.Name}
			} else if leftExpr, ok := binaryExpr.X.(ast.Expr); ok {
				dVar := &DynamicVariable{}

				exprInstructions, err := expressionToMLOG(dVar, leftExpr, options)
				if err != nil {
					return nil, err
				}

				instructions = append(instructions, exprInstructions...)
				leftSide = dVar
			} else {
				return nil, errors.New(fmt.Sprintf("unknown left side expression type: %T", binaryExpr.X))
			}

			if basicLit, ok := binaryExpr.Y.(*ast.BasicLit); ok {
				rightSide = &Value{Value: basicLit.Value}
			} else if rightIdent, ok := binaryExpr.Y.(*ast.Ident); ok {
				rightSide = &NormalVariable{Name: rightIdent.Name}
			} else if rightExpr, ok := binaryExpr.Y.(ast.Expr); ok {
				dVar := &DynamicVariable{}

				exprInstructions, err := expressionToMLOG(dVar, rightExpr, options)
				if err != nil {
					return nil, err
				}

				instructions = append(instructions, exprInstructions...)
				rightSide = dVar
			} else {
				return nil, errors.New(fmt.Sprintf("unknown right side expression type: %T", binaryExpr.Y))
			}

			return append(instructions, &MLOG{
				Comment: "Execute operation",
				Statement: [][]Resolvable{
					{
						&Value{Value: "op"},
						&Value{Value: opTranslated},
						ident,
						leftSide,
						rightSide,
					},
				},
			}), nil
		} else {
			return nil, errors.New(fmt.Sprintf("operator statement cannot use this operation: %s", binaryExpr.Op.String()))
		}
	case *ast.CallExpr:
		callInstructions, err := callExprToMLOG(expr.(*ast.CallExpr), options)
		if err != nil {
			return nil, err
		}

		callInstructions = append(callInstructions, &MLOG{
			Comment: "Set the variable to the value",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					ident,
					&Value{Value: functionReturnVariable},
				},
			},
		})

		return callInstructions, err
	case *ast.UnaryExpr:
		unaryExpr := expr.(*ast.UnaryExpr)

		if _, ok := regularOperators[unaryExpr.Op]; ok {
			instructions := make([]MLOGStatement, 0)

			var x Resolvable
			if basicLit, ok := unaryExpr.X.(*ast.BasicLit); ok {
				x = &Value{Value: basicLit.Value}
			} else if leftIdent, ok := unaryExpr.X.(*ast.Ident); ok {
				x = &NormalVariable{Name: leftIdent.Name}
			} else if leftExpr, ok := unaryExpr.X.(ast.Expr); ok {
				dVar := &DynamicVariable{}

				exprInstructions, err := expressionToMLOG(dVar, leftExpr, options)
				if err != nil {
					return nil, err
				}

				instructions = append(instructions, exprInstructions...)
				x = dVar
			} else {
				return nil, errors.New(fmt.Sprintf("unknown unary expression type: %T", unaryExpr.X))
			}

			var statement []Resolvable
			switch unaryExpr.Op {
			case token.NOT:
				statement = []Resolvable{
					&Value{Value: "op"},
					&Value{Value: regularOperators[token.NOT]},
					ident,
					x,
					&Value{Value: "-1"},
				}
				break
			case token.SUB:
				statement = []Resolvable{
					&Value{Value: "op"},
					&Value{Value: regularOperators[token.MUL]},
					ident,
					x,
					&Value{Value: "-1"},
				}
				break
			}

			if statement == nil {
				return nil, errors.New(fmt.Sprintf("unsupported unary operation: %s", unaryExpr.Op.String()))
			}

			return append(instructions, &MLOG{
				Comment:   "Execute unary operation",
				Statement: [][]Resolvable{statement},
			}), nil
		} else {
			return nil, errors.New(fmt.Sprintf("operator statement cannot use this operation: %s", unaryExpr.Op.String()))
		}
	case *ast.ParenExpr:
		parenExpr := expr.(*ast.ParenExpr)
		instructions, err := expressionToMLOG(ident, parenExpr.X, options)
		if err != nil {
			return nil, err
		}
		return instructions, nil
	default:
		spew.Dump(expr)
		return nil, errors.New(fmt.Sprintf("unsupported expression type: %T", expr))
	}
}

func callExprToMLOG(callExpr *ast.CallExpr, options Options) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)

	var funcName string
	if identity, ok := callExpr.Fun.(*ast.Ident); ok {
		funcName = identity.Name
	} else if selector, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		funcName = selector.Sel.Name
	} else {
		return nil, errors.New(fmt.Sprintf("unknown call expression: %T", callExpr.Fun))
	}

	if translatedFunc, ok := funcTranslations[funcName]; ok {
		args, instructions, err := argumentsToResolvables(callExpr.Args, options)
		if err != nil {
			return nil, err
		}
		results = append(results, instructions...)
		results = append(results, &MLOGFunc{
			Function:  translatedFunc,
			Arguments: args,
		})
	} else {
		for _, arg := range callExpr.Args {
			results = append(results, &MLOGStackWriter{
				Action: "add",
			})

			var value Resolvable
			if basicLit, ok := arg.(*ast.BasicLit); ok {
				value = &Value{Value: basicLit.Value}
			} else if ident, ok := arg.(*ast.Ident); ok {
				value = &NormalVariable{Name: ident.Name}
			} else if binaryExpr, ok := arg.(*ast.BinaryExpr); ok {
				dVar := &DynamicVariable{}

				instructions, err := expressionToMLOG(dVar, binaryExpr, options)
				if err != nil {
					return nil, err
				}

				results = append(results, instructions...)
				value = dVar
			} else {
				return nil, errors.New(fmt.Sprintf("unknown argument type: %T", arg))
			}

			results = append(results, &MLOG{
				Comment: "Write argument to memory",
				Statement: [][]Resolvable{
					{
						&Value{Value: "write"},
						value,
						&Value{Value: stackCellName},
						&Value{Value: stackVariable},
					},
				},
			})
		}

		results = append(results, &MLOGStackWriter{
			Action: "add",
		})

		extra := 2
		if options.Debug {
			extra += debugCount
		}

		results = append(results, &MLOGTrampoline{
			Variable: stackVariable,
			Extra:    extra,
		})

		results = append(results, &MLOGJump{
			MLOG: MLOG{
				Comment: "Jump to function: " + funcName,
			},
			Condition: []Resolvable{
				&Value{Value: "always"},
			},
			JumpTarget: &FunctionJumpTarget{
				FunctionName: funcName,
			},
		})

		results = append(results, &MLOGStackWriter{
			Action: "sub",
			Extra:  len(callExpr.Args),
		})
	}

	return results, nil
}

func argumentsToResolvables(args []ast.Expr, options Options) ([]Resolvable, []MLOGStatement, error) {
	result := make([]Resolvable, len(args))
	instructions := make([]MLOGStatement, 0)

	for i, arg := range args {
		if basicExpr, ok := arg.(*ast.BasicLit); ok {
			result[i] = &Value{Value: basicExpr.Value}
		} else if identExpr, ok := arg.(*ast.Ident); ok {
			result[i] = &NormalVariable{Name: identExpr.Name}
		} else if expr, ok := arg.(ast.Expr); ok {
			dVar := &DynamicVariable{}

			exprInstructions, err := expressionToMLOG(dVar, expr, options)
			if err != nil {
				return nil, nil, err
			}

			instructions = append(instructions, exprInstructions...)

			result[i] = dVar
		} else {
			return nil, nil, errors.New(fmt.Sprintf("only arguments of basic types or variables are supported, received: %T", arg))
		}
	}

	return result, instructions, nil
}
