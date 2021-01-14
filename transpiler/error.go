package transpiler

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

type ContextualError struct {
	error
	Context context.Context
	Pos     *token.Pos
}

func (e ContextualError) Error() string {
	if e.Pos != nil {
		return fmt.Sprintf("error at %d: %s", *e.Pos, e.error.Error())
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

func ErrPos(ctx context.Context, pos token.Pos, err string) ContextualError {
	return ContextualError{
		error:   errors.New(err),
		Context: ctx,
		Pos:     &pos,
	}
}
