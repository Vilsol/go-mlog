package transpiler

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func expressionToMLOG(ctx context.Context, ident []Resolvable, expr ast.Expr) ([]MLOGStatement, error) {
	switch castExpr := expr.(type) {
	case *ast.BasicLit:
		value := castExpr.Value
		if castExpr.Kind == token.CHAR {
			value = "\"" + strings.Trim(value, "'") + "\""
		}

		return []MLOGStatement{&MLOG{
			Comment: "Set the variable to the value",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					ident[0],
					&Value{Value: value},
				},
			},
		}}, nil
	case *ast.Ident:
		if castExpr.Name == "true" || castExpr.Name == "false" {
			return []MLOGStatement{&MLOG{
				Comment: "Set the variable to the value",
				Statement: [][]Resolvable{
					{
						&Value{Value: "set"},
						ident[0],
						&Value{Value: castExpr.Name},
					},
				},
			}}, nil
		}

		return []MLOGStatement{&MLOG{
			Comment: "Set the variable to the value of other variable",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					ident[0],
					&NormalVariable{Name: castExpr.Name},
				},
			},
		}}, nil
	case *ast.BinaryExpr:
		if opTranslated, ok := regularOperators[castExpr.Op]; ok {
			instructions := make([]MLOGStatement, 0)
			var leftSide Resolvable
			var rightSide Resolvable

			// TODO Convert to switch
			if basicLit, ok := castExpr.X.(*ast.BasicLit); ok {
				leftSide = &Value{Value: basicLit.Value}
			} else if leftIdent, ok := castExpr.X.(*ast.Ident); ok {
				leftSide = &NormalVariable{Name: leftIdent.Name}
			} else if leftExpr, ok := castExpr.X.(ast.Expr); ok {
				dVar := &DynamicVariable{}

				exprInstructions, err := expressionToMLOG(ctx, []Resolvable{dVar}, leftExpr)
				if err != nil {
					return nil, err
				}

				instructions = append(instructions, exprInstructions...)
				leftSide = dVar
			} else {
				return nil, Err(ctx, fmt.Sprintf("unknown left side expression type: %T", castExpr.X))
			}

			// TODO Convert to switch
			if basicLit, ok := castExpr.Y.(*ast.BasicLit); ok {
				rightSide = &Value{Value: basicLit.Value}
			} else if rightIdent, ok := castExpr.Y.(*ast.Ident); ok {
				rightSide = &NormalVariable{Name: rightIdent.Name}
			} else if rightExpr, ok := castExpr.Y.(ast.Expr); ok {
				dVar := &DynamicVariable{}

				exprInstructions, err := expressionToMLOG(ctx, []Resolvable{dVar}, rightExpr)
				if err != nil {
					return nil, err
				}

				instructions = append(instructions, exprInstructions...)
				rightSide = dVar
			} else {
				return nil, Err(ctx, fmt.Sprintf("unknown right side expression type: %T", castExpr.Y))
			}

			return append(instructions, &MLOG{
				Comment: "Execute operation",
				Statement: [][]Resolvable{
					{
						&Value{Value: "op"},
						&Value{Value: opTranslated},
						ident[0],
						leftSide,
						rightSide,
					},
				},
			}), nil
		} else {
			return nil, Err(ctx, fmt.Sprintf("operator statement cannot use this operation: %s", castExpr.Op.String()))
		}
	case *ast.CallExpr:
		callInstructions, err := callExprToMLOG(ctx, castExpr, ident)
		if err != nil {
			return nil, err
		}
		return callInstructions, err
	case *ast.UnaryExpr:
		if _, ok := regularOperators[castExpr.Op]; ok {
			instructions := make([]MLOGStatement, 0)

			var x Resolvable
			// TODO Convert to switch
			if basicLit, ok := castExpr.X.(*ast.BasicLit); ok {
				x = &Value{Value: basicLit.Value}
			} else if leftIdent, ok := castExpr.X.(*ast.Ident); ok {
				x = &NormalVariable{Name: leftIdent.Name}
			} else if leftExpr, ok := castExpr.X.(ast.Expr); ok {
				dVar := &DynamicVariable{}

				exprInstructions, err := expressionToMLOG(ctx, []Resolvable{dVar}, leftExpr)
				if err != nil {
					return nil, err
				}

				instructions = append(instructions, exprInstructions...)
				x = dVar
			} else {
				return nil, Err(ctx, fmt.Sprintf("unknown unary expression type: %T", castExpr.X))
			}

			var statement []Resolvable
			switch castExpr.Op {
			case token.NOT:
				statement = []Resolvable{
					&Value{Value: "op"},
					&Value{Value: regularOperators[token.NOT]},
					ident[0],
					x,
					&Value{Value: "-1"},
				}
				break
			case token.SUB:
				statement = []Resolvable{
					&Value{Value: "op"},
					&Value{Value: regularOperators[token.MUL]},
					ident[0],
					x,
					&Value{Value: "-1"},
				}
				break
			}

			if statement == nil {
				return nil, Err(ctx, fmt.Sprintf("unsupported unary operation: %s", castExpr.Op.String()))
			}

			return append(instructions, &MLOG{
				Comment:   "Execute unary operation",
				Statement: [][]Resolvable{statement},
			}), nil
		} else {
			return nil, Err(ctx, fmt.Sprintf("operator statement cannot use this operation: %s", castExpr.Op.String()))
		}
	case *ast.ParenExpr:
		instructions, err := expressionToMLOG(ctx, ident, castExpr.X)
		if err != nil {
			return nil, err
		}
		return instructions, nil
	case *ast.SelectorExpr:
		mlog, _, err := selectorExprToMLOG(ctx, ident[0], castExpr)
		return mlog, err
	default:
		return nil, Err(ctx, fmt.Sprintf("unsupported expression type: %T", expr))
	}
}

func selectorExprToMLOG(ctx context.Context, ident Resolvable, selectorExpr *ast.SelectorExpr) ([]MLOGStatement, string, error) {
	if _, ok := selectorExpr.X.(*ast.Ident); !ok {
		return nil, "", Err(ctx, fmt.Sprintf("unsupported selector type: %T", selectorExpr.X))
	}

	name := selectorExpr.X.(*ast.Ident).Name + "." + selectorExpr.Sel.Name
	if selector, ok := selectors[name]; ok {
		if ident == nil {
			return nil, selector, nil
		} else {
			return []MLOGStatement{
				&MLOG{
					Comment: "Set the variable to the value",
					Statement: [][]Resolvable{
						{
							&Value{Value: "set"},
							ident,
							&Value{Value: selector},
						},
					},
				},
			}, "", nil
		}
	}

	return nil, "", Err(ctx, fmt.Sprintf("unknown selector: %s", name))
}

func callExprToMLOG(ctx context.Context, callExpr *ast.CallExpr, ident []Resolvable) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)

	var funcName string
	// TODO Convert to switch
	if identity, ok := callExpr.Fun.(*ast.Ident); ok {
		funcName = identity.Name
	} else if selector, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		funcName = selector.X.(*ast.Ident).Name + "." + selector.Sel.Name
	} else {
		return nil, Err(ctx, fmt.Sprintf("unknown call expression: %T", callExpr.Fun))
	}

	if translatedFunc, ok := funcTranslations[funcName]; ok {
		args, instructions, err := argumentsToResolvables(ctx, callExpr.Args)
		if err != nil {
			return nil, err
		}
		results = append(results, instructions...)
		results = append(results, &MLOGFunc{
			Function:  translatedFunc,
			Arguments: args,
			Variables: ident,
		})
	} else {
		for _, arg := range callExpr.Args {
			results = append(results, &MLOGStackWriter{
				Action: "add",
			})

			var value Resolvable
			// TODO Convert to switch
			if basicLit, ok := arg.(*ast.BasicLit); ok {
				value = &Value{Value: basicLit.Value}
			} else if ident, ok := arg.(*ast.Ident); ok {
				value = &NormalVariable{Name: ident.Name}
			} else if argExpr, ok := arg.(ast.Expr); ok {
				dVar := &DynamicVariable{}

				instructions, err := expressionToMLOG(ctx, []Resolvable{dVar}, argExpr)
				if err != nil {
					return nil, err
				}

				results = append(results, instructions...)
				value = dVar
			} else {
				return nil, Err(ctx, fmt.Sprintf("unknown argument type: %T", arg))
			}

			results = append(results, &MLOG{
				Comment: "Write argument to memory",
				Statement: [][]Resolvable{
					{
						&Value{Value: "write"},
						value,
						&Value{Value: StackCellName},
						&Value{Value: stackVariable},
					},
				},
			})
		}

		results = append(results, &MLOGStackWriter{
			Action: "add",
		})

		extra := 2
		if ctx.Value(contextOptions).(Options).Debug {
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

		if len(ident) > 0 {
			results = append(results, &MLOG{
				Comment: "Set the variable to the value",
				Statement: [][]Resolvable{
					{
						&Value{Value: "set"},
						ident[0],
						&Value{Value: FunctionReturnVariable},
					},
				},
			})
		}
	}

	return results, nil
}

func argumentsToResolvables(ctx context.Context, args []ast.Expr) ([]Resolvable, []MLOGStatement, error) {
	result := make([]Resolvable, len(args))
	instructions := make([]MLOGStatement, 0)

	for i, arg := range args {
		// TODO Convert to switch
		if basicExpr, ok := arg.(*ast.BasicLit); ok {
			result[i] = &Value{Value: basicExpr.Value}
		} else if identExpr, ok := arg.(*ast.Ident); ok {
			if identExpr.Name == "true" || identExpr.Name == "false" {
				result[i] = &Value{Value: identExpr.Name}
			} else {
				result[i] = &NormalVariable{Name: identExpr.Name}
			}
		} else if selectorExpr, ok := arg.(*ast.SelectorExpr); ok {
			_, str, err := selectorExprToMLOG(ctx, nil, selectorExpr)
			if err != nil {
				return nil, nil, err
			}
			result[i] = &Value{Value: str}
		} else if expr, ok := arg.(ast.Expr); ok {
			dVar := &DynamicVariable{}

			exprInstructions, err := expressionToMLOG(ctx, []Resolvable{dVar}, expr)
			if err != nil {
				return nil, nil, err
			}

			instructions = append(instructions, exprInstructions...)

			result[i] = dVar
		} else {
			return nil, nil, Err(ctx, fmt.Sprintf("unknown argument type received: %T", arg))
		}
	}

	return result, instructions, nil
}
