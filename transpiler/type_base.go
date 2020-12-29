package transpiler

import (
	"context"
	"go/ast"
)

type Global struct {
	Functions []*Function
	Constants map[string]bool
}

type Function struct {
	Name            string
	Called          bool
	Declaration     *ast.FuncDecl
	Statements      []MLOGStatement
	ArgumentCount   int
	VariableCounter int
}

type MLOGAble interface {
	ToMLOG() [][]Resolvable
	GetComment() string
}

type Processable interface {
	PreProcess(context.Context, *Global, *Function) error
	PostProcess(context.Context, *Global, *Function) error
}

type WithPosition interface {
	GetPosition() int
	Size() int
}

type MutablePosition interface {
	SetPosition(int) int
}

type MLOGStatement interface {
	MLOGAble
	WithPosition
	MutablePosition
	Processable
}

type JumpTarget interface {
	Processable
	WithPosition
}

type Resolvable interface {
	Processable
	GetValue() string
}
