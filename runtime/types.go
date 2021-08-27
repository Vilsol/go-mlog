package runtime

import (
	"strings"
)

type ExecutionContext struct {
	Variables   map[string]*Variable
	PrintBuffer strings.Builder
	DrawBuffer  []DrawStatement
	Objects     map[string]interface{}
	Metrics     map[int64]*Metrics
}

type Metrics struct {
	Executions uint64
}

type Variable struct {
	Value    interface{}
	Constant bool
}

type MLOGLine struct {
	Instruction []string
	Comment     string
	SourceLine  int
}

type Operation struct {
	Line     MLOGLine
	Executor OperationExecutor
}

type OperationExecutor func(ctx *ExecutionContext)

type OperationSetup func(args []string) (OperationExecutor, error)

type Message interface {
	PrintFlush(buffer string)
}

type DrawAction string

const (
	DrawActionClear    = DrawAction("clear")
	DrawActionColor    = DrawAction("color")
	DrawActionStroke   = DrawAction("stroke")
	DrawActionLine     = DrawAction("line")
	DrawActionRect     = DrawAction("rect")
	DrawActionLineRect = DrawAction("lineRect")
	DrawActionPoly     = DrawAction("poly")
	DrawActionLinePoly = DrawAction("linePoly")
	DrawActionTriangle = DrawAction("triangle")
	DrawActionImage    = DrawAction("image")
)

type DrawStatement struct {
	Action    DrawAction
	Arguments []interface{}
}

type Display interface {
	DrawFlush(buffer []DrawStatement)
}

type Memory interface {
	Write(value float64, position int64)
	Read(position int64) float64
}

type PostExecute interface {
	PostExecute()
}
