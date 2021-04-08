package transpiler

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/types"
)

type ContextualError struct {
	error
	Context context.Context
	Pos     *ast.Node
}

func (e ContextualError) Error() string {
	if e.Pos != nil {
		return fmt.Sprintf("error at %d-%d: %s", (*e.Pos).Pos(), (*e.Pos).End(), e.error.Error())
	}

	if e.Context != nil {
		if stmt, ok := e.Context.Value(contextStatement).(ast.Stmt); ok {
			return fmt.Sprintf("error at %d: %s", stmt.Pos(), e.error.Error())
		} else if fn, ok := e.Context.Value(contextFunction).(*ast.FuncDecl); ok {
			return fmt.Sprintf("error at %d: %s", fn.Pos(), e.error.Error())
		} else if spec, ok := e.Context.Value(contextSpec).(ast.Spec); ok {
			return fmt.Sprintf("error at %d: %s", spec.Pos(), e.error.Error())
		} else if decl, ok := e.Context.Value(contextDecl).(ast.Decl); ok {
			return fmt.Sprintf("error at %d: %s", decl.Pos(), e.error.Error())
		} else if tError, ok := e.Context.Value(typeError).(types.Error); ok {
			return fmt.Sprintf("error at %d: %s", tError.Pos, e.error.Error())
		}
	}
	return e.error.Error()
}

func Err(ctx context.Context, err string) ContextualError {
	return ContextualError{
		error:   errors.New(err),
		Context: ctx,
	}
}

func ErrPos(ctx context.Context, pos ast.Node, err string) ContextualError {
	return ContextualError{
		error:   errors.New(err),
		Context: ctx,
		Pos:     &pos,
	}
}
