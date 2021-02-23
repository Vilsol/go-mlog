package transpiler

import (
	"context"
	"go/ast"
	"strconv"
)

type MLOGJump struct {
	MLOG
	Condition  []Resolvable
	JumpTarget JumpTarget
}

func (m *MLOGJump) ToMLOG() [][]Resolvable {
	return [][]Resolvable{
		append([]Resolvable{
			&Value{Value: "jump"},
			&Value{Value: strconv.Itoa(m.JumpTarget.GetPosition())},
		}, m.Condition...),
	}
}

func (m *MLOGJump) Size() int {
	return 1
}

func (m *MLOGJump) PreProcess(ctx context.Context, global *Global, function *Function) error {
	for _, resolvable := range m.Condition {
		if err := resolvable.PreProcess(ctx, global, function); err != nil {
			return err
		}
	}
	return m.JumpTarget.PreProcess(ctx, global, function)
}

func (m *MLOGJump) PostProcess(ctx context.Context, global *Global, function *Function) error {
	for _, resolvable := range m.Condition {
		if err := resolvable.PostProcess(ctx, global, function); err != nil {
			return err
		}
	}
	return m.JumpTarget.PostProcess(ctx, global, function)
}

func (m *MLOGJump) GetComment(int) string {
	if m.Comment == "" {
		return "Jump to target"
	}
	return m.Comment
}

type FunctionJumpTarget struct {
	Statement    WithPosition
	FunctionName string
	SourcePos    ast.Node
}

func (m *FunctionJumpTarget) GetPosition() int {
	return m.Statement.GetPosition()
}

func (m *FunctionJumpTarget) Size() int {
	return 1
}

func (m *FunctionJumpTarget) PreProcess(ctx context.Context, global *Global, _ *Function) error {
	for _, fn := range global.Functions {
		if fn.Name == m.FunctionName {
			fn.Called = true
			m.Statement = fn.Statements[0]
			return nil
		}
	}
	return ErrPos(ctx, m.SourcePos, "unknown function: "+m.FunctionName)
}

func (m *FunctionJumpTarget) PostProcess(context.Context, *Global, *Function) error {
	return nil
}

type StatementJumpTarget struct {
	Statement WithPosition
	After     bool
}

func (m *StatementJumpTarget) GetPosition() int {
	if m.After {
		return m.Statement.GetPosition() + m.Statement.Size()
	}
	return m.Statement.GetPosition()
}

func (m *StatementJumpTarget) Size() int {
	return 1
}

func (m *StatementJumpTarget) PreProcess(context.Context, *Global, *Function) error {
	return nil
}

func (m *StatementJumpTarget) PostProcess(context.Context, *Global, *Function) error {
	return nil
}
