package transpiler

import (
	"errors"
	"fmt"
	"strconv"
)

type Global struct {
	Functions []*Function
	Constants map[string]bool
}

type Function struct {
	Name            string
	Statements      []MLOGStatement
	ArgumentCount   int
	VariableCounter int
}

type MLOGAble interface {
	ToMLOG() [][]Resolvable
	GetComment() string
}

type Processable interface {
	PostProcess(*Global, *Function) error
}

type WithPosition interface {
	GetPosition() int
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

func MLOGToString(statements [][]Resolvable, statement MLOGAble, lineNumber int, options Options) string {
	result := ""
	for _, line := range statements {
		resultLine := ""
		for _, t := range line {
			if resultLine != "" {
				resultLine += " "
			}
			resultLine += t.GetValue()
		}

		if options.Numbers {
			result += fmt.Sprintf("%3d: ", lineNumber)
		}

		if options.Comments {
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

func (m *MLOG) PostProcess(global *Global, function *Function) error {
	for _, resolvables := range m.Statement {
		for _, resolvable := range resolvables {
			if err := resolvable.PostProcess(global, function); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MLOG) GetPosition() int {
	return m.Position
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

func (m *MLOGFunc) SetPosition(position int) int {
	m.Position = position
	return m.Function.Count
}

func (m *MLOGFunc) PostProcess(global *Global, function *Function) error {
	if len(m.Variables) != m.Function.Variables {
		return errors.New(fmt.Sprintf("function requires %d variables, provided: %d", m.Function.Variables, len(m.Variables)))
	}

	for _, argument := range m.Arguments {
		if err := argument.PostProcess(global, function); err != nil {
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
		if err := statement.PostProcess(global, function); err != nil {
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

func (m *MLOGJump) PostProcess(global *Global, function *Function) error {
	for _, resolvable := range m.Condition {
		if err := resolvable.PostProcess(global, function); err != nil {
			return err
		}
	}
	return m.JumpTarget.PostProcess(global, function)
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

func (m *FunctionJumpTarget) PostProcess(global *Global, _ *Function) error {
	for _, fn := range global.Functions {
		if fn.Name == m.FunctionName {
			m.Statement = fn.Statements[0]
			return nil
		}
	}
	return errors.New("unknown function: " + m.FunctionName)
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

func (m *Value) PostProcess(*Global, *Function) error {
	return nil
}

func (m *Value) String() string {
	return m.Value
}

type NormalVariable struct {
	Name           string
	CalculatedName string
}

func (m *NormalVariable) PostProcess(global *Global, function *Function) error {
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

func (m *DynamicVariable) PostProcess(global *Global, function *Function) error {
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
		return m.Statement.GetPosition() + 1
	}
	return m.Statement.GetPosition()
}

func (m *StatementJumpTarget) PostProcess(*Global, *Function) error {
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
