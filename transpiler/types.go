package transpiler

import (
	"context"
	"fmt"
	"go/ast"
	"strconv"
)

type Global struct {
	Functions []*Function
	Constants map[string]bool
}

type Function struct {
	Name            string
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

func MLOGToString(ctx context.Context, statements [][]Resolvable, statement MLOGAble, lineNumber int) string {
	result := ""
	for _, line := range statements {
		resultLine := ""
		for _, t := range line {
			if resultLine != "" {
				resultLine += " "
			}
			resultLine += t.GetValue()
		}

		if ctx.Value(contextOptions).(Options).Numbers {
			result += fmt.Sprintf("%3d: ", lineNumber)
		}

		if ctx.Value(contextOptions).(Options).Comments {
			result += fmt.Sprintf("%-45s", resultLine)
			result += " // " + statement.GetComment()
		} else {
			result += resultLine
		}

		result += "\n"
		lineNumber++
	}
	return result
}

type MLOG struct {
	Statement [][]Resolvable
	Position  int
	Comment   string
}

func (m *MLOG) ToMLOG() [][]Resolvable {
	return m.Statement
}

func (m *MLOG) PostProcess(ctx context.Context, global *Global, function *Function) error {
	for _, resolvables := range m.Statement {
		for _, resolvable := range resolvables {
			if err := resolvable.PostProcess(ctx, global, function); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MLOG) GetPosition() int {
	return m.Position
}

func (m *MLOG) Size() int {
	if len(m.Statement) == 0 {
		panic("statement without instructions")
	}

	return len(m.Statement)
}

func (m *MLOG) SetPosition(position int) int {
	m.Position = position
	return 1
}

func (m *MLOG) GetComment() string {
	return m.Comment
}

type MLOGFunc struct {
	Position   int
	Function   Translator
	Arguments  []Resolvable
	Variables  []Resolvable
	Unresolved []MLOGStatement
}

func (m *MLOGFunc) ToMLOG() [][]Resolvable {
	results := make([][]Resolvable, 0)
	for _, statement := range m.Unresolved {
		results = append(results, statement.ToMLOG()...)
	}
	return results
}

func (m *MLOGFunc) GetPosition() int {
	return m.Position
}

func (m *MLOGFunc) Size() int {
	return m.Function.Count
}

func (m *MLOGFunc) SetPosition(position int) int {
	m.Position = position
	return m.Function.Count
}

func (m *MLOGFunc) PostProcess(ctx context.Context, global *Global, function *Function) error {
	if len(m.Variables) != m.Function.Variables {
		return Err(ctx, fmt.Sprintf("function requires %d variables, provided: %d", m.Function.Variables, len(m.Variables)))
	}

	for _, argument := range m.Arguments {
		if err := argument.PostProcess(ctx, global, function); err != nil {
			return err
		}
	}

	var err error
	m.Unresolved, err = m.Function.Translate(m.Arguments, m.Variables)
	if err != nil {
		return err
	}

	for i, statement := range m.Unresolved {
		statement.SetPosition(m.Position + i)
		if err := statement.PostProcess(ctx, global, function); err != nil {
			return err
		}
	}
	return nil
}

func (m *MLOGFunc) GetComment() string {
	return "Call to native function"
}

type JumpTarget interface {
	Processable
	WithPosition
}

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

func (m *MLOGJump) PostProcess(ctx context.Context, global *Global, function *Function) error {
	for _, resolvable := range m.Condition {
		if err := resolvable.PostProcess(ctx, global, function); err != nil {
			return err
		}
	}
	return m.JumpTarget.PostProcess(ctx, global, function)
}

func (m *MLOGJump) GetComment() string {
	if m.Comment == "" {
		return "Jump to target"
	}
	return m.Comment
}

type FunctionJumpTarget struct {
	Statement    WithPosition
	FunctionName string
}

func (m *FunctionJumpTarget) GetPosition() int {
	return m.Statement.GetPosition()
}

func (m *FunctionJumpTarget) Size() int {
	return 1
}

func (m *FunctionJumpTarget) PostProcess(ctx context.Context, global *Global, _ *Function) error {
	for _, fn := range global.Functions {
		if fn.Name == m.FunctionName {
			m.Statement = fn.Statements[0]
			return nil
		}
	}
	return Err(ctx, "unknown function: "+m.FunctionName)
}

type Resolvable interface {
	Processable
	GetValue() string
}

type Value struct {
	Value string
}

func (m *Value) GetValue() string {
	return m.Value
}

func (m *Value) PostProcess(context.Context, *Global, *Function) error {
	return nil
}

func (m *Value) String() string {
	return m.Value
}

type NormalVariable struct {
	Name           string
	CalculatedName string
}

func (m *NormalVariable) PostProcess(ctx context.Context, global *Global, function *Function) error {
	if m.CalculatedName == "" {
		if _, ok := global.Constants[m.Name]; ok {
			m.CalculatedName = m.Name
		} else {
			m.CalculatedName = "_" + function.Name + "_" + m.Name
		}
	}
	return nil
}

func (m *NormalVariable) GetValue() string {
	if m.CalculatedName == "" {
		panic("PostProcess not called on NormalVariable (" + m.Name + ")")
	}
	return m.CalculatedName
}

type DynamicVariable struct {
	Name string
}

func (m *DynamicVariable) PostProcess(ctx context.Context, global *Global, function *Function) error {
	if m.Name == "" {
		suffix := function.VariableCounter
		function.VariableCounter += 1
		m.Name = "_" + function.Name + "_" + strconv.Itoa(suffix)
	}
	return nil
}

func (m *DynamicVariable) GetValue() string {
	if m.Name == "" {
		panic("PostProcess not called on DynamicVariable")
	}
	return m.Name
}

type MLOGTrampoline struct {
	MLOG
	Variable string
	Extra    int
}

func (m *MLOGTrampoline) ToMLOG() [][]Resolvable {
	return [][]Resolvable{
		{
			&Value{Value: "write"},
			&Value{Value: strconv.Itoa(m.Position + m.Extra)},
			&Value{Value: StackCellName},
			&Value{Value: m.Variable},
		},
	}
}

func (m *MLOGTrampoline) GetComment() string {
	return "Set Trampoline Address"
}

type MLOGStackWriter struct {
	MLOG
	Action string
	Extra  int
}

func (m *MLOGStackWriter) ToMLOG() [][]Resolvable {
	return [][]Resolvable{
		{
			&Value{Value: "op"},
			&Value{Value: m.Action},
			&Value{Value: stackVariable},
			&Value{Value: stackVariable},
			&Value{Value: strconv.Itoa(1 + m.Extra)},
		},
	}
}

func (m *MLOGStackWriter) GetComment() string {
	return "Update Stack Pointer"
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

func (m *StatementJumpTarget) PostProcess(context.Context, *Global, *Function) error {
	return nil
}

type MLOGTrampolineBack struct {
	MLOG
}

func (m *MLOGTrampolineBack) ToMLOG() [][]Resolvable {
	return [][]Resolvable{
		{
			&Value{Value: "read"},
			&Value{Value: "@counter"},
			&Value{Value: StackCellName},
			&Value{Value: stackVariable},
		},
	}
}

func (m *MLOGTrampolineBack) GetComment() string {
	return "Trampoline back"
}

type MLOGBreak struct {
	MLOG
	Block *ContextBlock
}

func (m *MLOGBreak) ToMLOG() [][]Resolvable {
	lastStatement := m.Block.Statements[len(m.Block.Statements)-1]
	if m.Block.Extra != nil && len(m.Block.Extra) > 0 {
		lastStatement = m.Block.Extra[len(m.Block.Extra)-1]
	}
	return [][]Resolvable{
		{
			&Value{Value: "jump"},
			&Value{Value: strconv.Itoa(lastStatement.GetPosition() + lastStatement.Size())},
			&Value{Value: "always"},
		},
	}
}

func (m *MLOGBreak) Size() int {
	return 1
}

func (m *MLOGBreak) GetComment() string {
	return "Break"
}

type MLOGContinue struct {
	MLOG
	Block *ContextBlock
}

func (m *MLOGContinue) ToMLOG() [][]Resolvable {
	lastStatement := m.Block.Statements[len(m.Block.Statements)-1]
	return [][]Resolvable{
		{
			&Value{Value: "jump"},
			&Value{Value: strconv.Itoa(lastStatement.GetPosition() + lastStatement.Size())},
			&Value{Value: "always"},
		},
	}
}

func (m *MLOGContinue) Size() int {
	return 1
}

func (m *MLOGContinue) GetComment() string {
	return "Continue"
}

type MLOGFallthrough struct {
	MLOG
	Block *ContextBlock
}

func (m *MLOGFallthrough) ToMLOG() [][]Resolvable {
	lastStatement := m.Block.Statements[len(m.Block.Statements)-1]
	if m.Block.Extra != nil && len(m.Block.Extra) > 0 {
		lastStatement = m.Block.Extra[len(m.Block.Extra)-1]
	}
	return [][]Resolvable{
		{
			&Value{Value: "jump"},
			&Value{Value: strconv.Itoa(lastStatement.GetPosition() + lastStatement.Size())},
			&Value{Value: "always"},
		},
	}
}

func (m *MLOGFallthrough) Size() int {
	return 1
}

func (m *MLOGFallthrough) GetComment() string {
	return "Fallthrough"
}
