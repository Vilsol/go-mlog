package transpiler

import (
	"context"
	"fmt"
	"go/ast"
)

func getSuggestedDynamicVariableCount(ctx context.Context, expr ast.Expr) ([]Resolvable, error) {
	if callExpr, ok := expr.(*ast.CallExpr); ok {
		count, err := getFunctionReturnCount(ctx, callExpr)

		if err != nil {
			return nil, err
		}

		resolvables := make([]Resolvable, count)
		for i := 0; i < count; i++ {
			resolvables[i] = &DynamicVariable{}
		}

		return resolvables, nil
	}

	return []Resolvable{&DynamicVariable{}}, nil
}

func getFunctionReturnCount(ctx context.Context, callExpr *ast.CallExpr) (int, error) {
	var funcName, exprName, selName string
	switch funType := callExpr.Fun.(type) {
	case *ast.Ident:
		funcName = funType.Name
		break
	case *ast.SelectorExpr:
		exprName = funType.X.(*ast.Ident).Name
		selName = funType.Sel.Name
		funcName = exprName + "." + selName
		break
	default:
		return 0, Err(ctx, fmt.Sprintf("unknown call expression: %T", callExpr.Fun))
	}

	if translatedFunc, ok := funcTranslations[funcName]; ok {
		return translatedFunc.Variables, nil
	} else if translatedFunc, ok := funcTranslations[selName]; ok {
		return translatedFunc.Variables, nil
	} else {
		global := ctx.Value(contextGlobal).(*Global)
		for _, function := range global.Functions {
			if function.Name == funcName {
				return len(function.Declaration.Type.Results.List), nil
			}
		}
		return 0, nil
	}
}
