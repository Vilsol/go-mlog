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
		return basicLitToMLOG(ctx, ident, castExpr)
	case *ast.Ident:
		return identToMLOG(ctx, ident, castExpr)
	case *ast.BinaryExpr:
		return binaryExprToMLOG(ctx, ident, castExpr)
	case *ast.CallExpr:
		return callExprToMLOG(ctx, castExpr, ident)
	case *ast.UnaryExpr:
		return unaryExprToMLOG(ctx, ident, castExpr)
	case *ast.ParenExpr:
		return expressionToMLOG(ctx, ident, castExpr.X)
	case *ast.SelectorExpr:
		mlog, _, err := selectorExprToMLOG(ctx, ident[0], castExpr)
		return mlog, err
	}

	return nil, Err(ctx, fmt.Sprintf("unsupported expression type: %T", expr))
}

func exprToResolvable(ctx context.Context, expr ast.Expr) ([]Resolvable, []MLOGStatement, error) {
	switch castUnary := expr.(type) {
	case *ast.BasicLit:
		return []Resolvable{&Value{Value: castUnary.Value}}, nil, nil
	case *ast.Ident:
		if castUnary.Name == "true" || castUnary.Name == "false" {
			return []Resolvable{&Value{Value: castUnary.Name}}, nil, nil
		} else {
			return []Resolvable{&NormalVariable{Name: castUnary.Name}}, nil, nil
		}
	case ast.Expr:
		dVars, err := getSuggestedDynamicVariableCount(ctx, castUnary)

		if err != nil {
			return nil, nil, err
		}

		exprInstructions, err := expressionToMLOG(ctx, dVars, castUnary)
		if err != nil {
			return nil, nil, err
		}

		return dVars, exprInstructions, nil
	}

	return nil, nil, Err(ctx, fmt.Sprintf("unknown resolvable expression type: %T", expr))
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
					Comment: "Assign value to variable",
					Statement: [][]Resolvable{
						{
							&Value{Value: "set"},
							ident,
							&Value{Value: selector},
						},
					},
					SourcePos: ctx.Value(contextStatement).(ast.Node),
				},
			}, "", nil
		}
	}

	return nil, "", Err(ctx, fmt.Sprintf("unknown selector: %s", name))
}

func callExprToMLOG(ctx context.Context, callExpr *ast.CallExpr, ident []Resolvable) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)
	var subSelector Resolvable

	var funcName, exprName, selName string
	switch funType := callExpr.Fun.(type) {
	case *ast.Ident:
		funcName = funType.Name
		break
	case *ast.SelectorExpr:
		switch xType := funType.X.(type) {
		case *ast.Ident:
			exprName = xType.Name
			selName = funType.Sel.Name
			funcName = exprName + "." + selName
			break
		case *ast.SelectorExpr:
			_, str, err := selectorExprToMLOG(ctx, nil, xType)

			if err != nil {
				return nil, err
			}

			subSelector = &Value{Value: str}

			// TODO Recursive check
			exprName = xType.Sel.Name
			selName = funType.Sel.Name
			funcName = exprName + "." + selName
			break
		}
		break
	default:
		return nil, Err(ctx, fmt.Sprintf("unknown call expression: %T", callExpr.Fun))
	}

	if inline, ok := inlineTranslations[funcName]; ok {
		resolvables, statements, err := argumentsToResolvables(ctx, callExpr.Args)

		if err != nil {
			return nil, err
		}

		results = append(results, statements...)

		resolvable, err := inline(resolvables)

		if err != nil {
			return nil, err
		}

		results = append(results, &MLOG{
			Comment: "Assign value to variable",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					ident[0],
					resolvable,
				},
			},
			SourcePos: ctx.Value(contextStatement).(ast.Node),
		})
	} else if translatedFunc, ok := funcTranslations[funcName]; ok {
		args, instructions, err := argumentsToResolvables(ctx, callExpr.Args)
		if err != nil {
			return nil, err
		}
		results = append(results, instructions...)
		results = append(results, &MLOGFunc{
			Function:  translatedFunc,
			Arguments: args,
			Variables: ident,
			SourcePos: callExpr,
		})
	} else if translatedFunc, ok := funcTranslations[selName]; ok {
		if subSelector == nil {
			subSelector = &NormalVariable{Name: exprName}
		}

		results = append(results, &MLOGFunc{
			Function: translatedFunc,
			Arguments: []Resolvable{
				subSelector,
			},
			Variables: ident,
			SourcePos: callExpr,
		})
	} else {
		results = append(results, &MLOGCustomFunction{
			Arguments:    callExpr.Args,
			Variables:    ident,
			FunctionName: funcName,
			SourcePos:    callExpr,
		})
	}

	return results, nil
}

func argumentsToResolvables(ctx context.Context, args []ast.Expr) ([]Resolvable, []MLOGStatement, error) {
	result := make([]Resolvable, 0)
	instructions := make([]MLOGStatement, 0)

	for _, arg := range args {
		switch argType := arg.(type) {
		case *ast.SelectorExpr:
			_, str, err := selectorExprToMLOG(ctx, nil, argType)
			if err != nil {
				return nil, nil, err
			}
			result = append(result, &Value{Value: str})
			break
		default:
			if callExpr, ok := argType.(*ast.CallExpr); ok {
				var funcName string
				switch funType := callExpr.Fun.(type) {
				case *ast.Ident:
					funcName = funType.Name
					break
				case *ast.SelectorExpr:
					switch xType := funType.X.(type) {
					case *ast.Ident:
						funcName = xType.Name + "." + funType.Sel.Name
						break
					case *ast.SelectorExpr:
						funcName = xType.Sel.Name + "." + funType.Sel.Name
						break
					}

					break
				default:
					return nil, nil, Err(ctx, fmt.Sprintf("unknown call expression: %T", callExpr.Fun))
				}

				if inline, ok := inlineTranslations[funcName]; ok {
					resolvables, statements, err := argumentsToResolvables(ctx, callExpr.Args)

					if err != nil {
						return nil, nil, err
					}

					instructions = append(instructions, statements...)

					resolvable, err := inline(resolvables)

					if err != nil {
						return nil, nil, err
					}

					result = append(result, resolvable)
					break
				}
			}

			res, leftExprInstructions, err := exprToResolvable(ctx, arg)
			if err != nil {
				return nil, nil, err
			}
			instructions = append(instructions, leftExprInstructions...)
			result = append(result, res...)
			break
		}
	}

	return result, instructions, nil
}

func unaryExprToMLOG(ctx context.Context, ident []Resolvable, expr *ast.UnaryExpr) ([]MLOGStatement, error) {
	if _, ok := regularOperators[expr.Op]; ok {
		instructions := make([]MLOGStatement, 0)

		x, exprInstructions, err := exprToResolvable(ctx, expr.X)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, exprInstructions...)

		if len(x) != 1 {
			return nil, Err(ctx, "unknown error")
		}

		var statement []Resolvable
		switch expr.Op {
		case token.NOT:
			statement = []Resolvable{
				&Value{Value: "op"},
				&Value{Value: regularOperators[token.NOT]},
				ident[0],
				x[0],
				&Value{Value: "-1"},
			}
			break
		case token.SUB:
			statement = []Resolvable{
				&Value{Value: "op"},
				&Value{Value: regularOperators[token.MUL]},
				ident[0],
				x[0],
				&Value{Value: "-1"},
			}
			break
		}

		if statement == nil {
			return nil, Err(ctx, fmt.Sprintf("unsupported unary operation: %s", expr.Op.String()))
		}

		return append(instructions, &MLOG{
			Comment:   "Execute unary operation",
			Statement: [][]Resolvable{statement},
		}), nil
	}

	return nil, Err(ctx, fmt.Sprintf("operator statement cannot use this operation: %s", expr.Op.String()))
}

func binaryExprToMLOG(ctx context.Context, ident []Resolvable, expr *ast.BinaryExpr) ([]MLOGStatement, error) {
	if opTranslated, ok := regularOperators[expr.Op]; ok {
		instructions := make([]MLOGStatement, 0)

		leftSide, leftExprInstructions, err := exprToResolvable(ctx, expr.X)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, leftExprInstructions...)

		if len(leftSide) != 1 {
			return nil, Err(ctx, "unknown error")
		}

		rightSide, rightExprInstructions, err := exprToResolvable(ctx, expr.Y)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, rightExprInstructions...)

		if len(rightSide) != 1 {
			return nil, Err(ctx, "unknown error")
		}

		return append(instructions, &MLOG{
			Comment: "Execute operation",
			Statement: [][]Resolvable{
				{
					&Value{Value: "op"},
					&Value{Value: opTranslated},
					ident[0],
					leftSide[0],
					rightSide[0],
				},
			},
			SourcePos: expr,
		}), nil
	}

	return nil, Err(ctx, fmt.Sprintf("operator statement cannot use this operation: %s", expr.Op.String()))
}

func identToMLOG(ctx context.Context, ident []Resolvable, expr *ast.Ident) ([]MLOGStatement, error) {
	if len(ident) < 1 {
		return nil, Err(ctx, "assignment identity not provided")
	}

	if expr.Name == "true" || expr.Name == "false" {
		return []MLOGStatement{&MLOG{
			Comment: "Assign value to variable",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					ident[0],
					&Value{Value: expr.Name},
				},
			},
			SourcePos: ctx.Value(contextStatement).(ast.Node),
		}}, nil
	}

	return []MLOGStatement{&MLOG{
		Comment: "Assign variable to another",
		Statement: [][]Resolvable{
			{
				&Value{Value: "set"},
				ident[0],
				&NormalVariable{Name: expr.Name},
			},
		},
		SourcePos: ctx.Value(contextStatement).(ast.Node),
	}}, nil
}

func basicLitToMLOG(ctx context.Context, ident []Resolvable, expr *ast.BasicLit) ([]MLOGStatement, error) {
	value := expr.Value
	if expr.Kind == token.CHAR {
		value = "\"" + strings.Trim(value, "'") + "\""
	}

	return []MLOGStatement{&MLOG{
		Comment: "Assign value to variable",
		Statement: [][]Resolvable{
			{
				&Value{Value: "set"},
				ident[0],
				&Value{Value: value},
			},
		},
		SourcePos: ctx.Value(contextStatement).(ast.Node),
	}}, nil
}
