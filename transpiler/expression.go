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

func exprToResolvable(ctx context.Context, expr ast.Expr) (Resolvable, []MLOGStatement, error) {
	switch castUnary := expr.(type) {
	case *ast.BasicLit:
		return &Value{Value: castUnary.Value}, nil, nil
	case *ast.Ident:
		if castUnary.Name == "true" || castUnary.Name == "false" {
			return &Value{Value: castUnary.Name}, nil, nil
		} else {
			return &NormalVariable{Name: castUnary.Name}, nil, nil
		}
	case ast.Expr:
		dVar := &DynamicVariable{}

		exprInstructions, err := expressionToMLOG(ctx, []Resolvable{dVar}, castUnary)
		if err != nil {
			return nil, nil, err
		}

		return dVar, exprInstructions, nil
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
	switch funType := callExpr.Fun.(type) {
	case *ast.Ident:
		funcName = funType.Name
		break
	case *ast.SelectorExpr:
		funcName = funType.X.(*ast.Ident).Name + "." + funType.Sel.Name
		break
	default:
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
			SourcePos: callExpr.Pos(),
		})
	} else {
		results = append(results, &MLOGCustomFunction{
			Arguments:    callExpr.Args,
			Variables:    ident,
			FunctionName: funcName,
			SourcePos:    callExpr.Pos(),
		})
	}

	return results, nil
}

func argumentsToResolvables(ctx context.Context, args []ast.Expr) ([]Resolvable, []MLOGStatement, error) {
	result := make([]Resolvable, len(args))
	instructions := make([]MLOGStatement, 0)

	for i, arg := range args {

		switch argType := arg.(type) {
		case *ast.SelectorExpr:
			_, str, err := selectorExprToMLOG(ctx, nil, argType)
			if err != nil {
				return nil, nil, err
			}
			result[i] = &Value{Value: str}
			break
		default:
			res, leftExprInstructions, err := exprToResolvable(ctx, arg)
			if err != nil {
				return nil, nil, err
			}
			instructions = append(instructions, leftExprInstructions...)
			result[i] = res
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

		var statement []Resolvable
		switch expr.Op {
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

		rightSide, rightExprInstructions, err := exprToResolvable(ctx, expr.Y)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, rightExprInstructions...)

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
	}

	return nil, Err(ctx, fmt.Sprintf("operator statement cannot use this operation: %s", expr.Op.String()))
}

func identToMLOG(_ context.Context, ident []Resolvable, expr *ast.Ident) ([]MLOGStatement, error) {
	if expr.Name == "true" || expr.Name == "false" {
		return []MLOGStatement{&MLOG{
			Comment: "Set the variable to the value",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					ident[0],
					&Value{Value: expr.Name},
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
				&NormalVariable{Name: expr.Name},
			},
		},
	}}, nil
}

func basicLitToMLOG(_ context.Context, ident []Resolvable, expr *ast.BasicLit) ([]MLOGStatement, error) {
	value := expr.Value
	if expr.Kind == token.CHAR {
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
}
