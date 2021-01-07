package transpiler

import (
	"context"
	"fmt"
	"go/token"
	"strconv"
)

type MLOG struct {
	Statement [][]Resolvable
	Position  int
	Comment   string
	SourcePos token.Pos
}

func (m *MLOG) ToMLOG() [][]Resolvable {
	return m.Statement
}

func (m *MLOG) PreProcess(ctx context.Context, global *Global, function *Function) error {
	for _, resolvables := range m.Statement {
		for _, resolvable := range resolvables {
			if err := resolvable.PreProcess(ctx, global, function); err != nil {
				return err
			}
		}
	}
	return nil
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

func (m *MLOG) GetComment(int) string {
	return m.Comment
}

func (m *MLOG) SetSourcePos(pos token.Pos) {
	m.SourcePos = pos
}

func (m *MLOG) GetSourcePos() token.Pos {
	return m.SourcePos
}

type MLOGFunc struct {
	Position   int
	Function   Translator
	Arguments  []Resolvable
	Variables  []Resolvable
	Unresolved []MLOGStatement
	SourcePos  token.Pos
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
	return m.Function.Count(m.Arguments, m.Variables)
}

func (m *MLOGFunc) SetPosition(position int) int {
	m.Position = position
	return m.Function.Count(m.Arguments, m.Variables)
}

func (m *MLOGFunc) PreProcess(ctx context.Context, global *Global, function *Function) error {
	if len(m.Variables) != m.Function.Variables {
		return Err(ctx, fmt.Sprintf("function requires %d variables, provided: %d", m.Function.Variables, len(m.Variables)))
	}

	for _, argument := range m.Arguments {
		if err := argument.PreProcess(ctx, global, function); err != nil {
			return err
		}
	}

	var err error
	m.Unresolved, err = m.Function.Translate(m.Arguments, m.Variables)
	if err != nil {
		return err
	}

	for _, statement := range m.Unresolved {
		if err := statement.PreProcess(ctx, global, function); err != nil {
			return err
		}
	}
	return nil
}

func (m *MLOGFunc) PostProcess(ctx context.Context, global *Global, function *Function) error {
	for _, argument := range m.Arguments {
		if err := argument.PostProcess(ctx, global, function); err != nil {
			return err
		}
	}

	for i, statement := range m.Unresolved {
		statement.SetPosition(m.Position + i)
		if err := statement.PostProcess(ctx, global, function); err != nil {
			return err
		}
	}
	return nil
}

func (m *MLOGFunc) GetComment(int) string {
	return "Call to native function"
}

func (m *MLOGFunc) SetSourcePos(pos token.Pos) {
	m.SourcePos = pos
}

func (m *MLOGFunc) GetSourcePos() token.Pos {
	return m.SourcePos
}

type MLOGTrampoline struct {
	MLOG
	Variable string
	Extra    int
	Stacked  string
	Function string
}

func (m *MLOGTrampoline) ToMLOG() [][]Resolvable {
	if m.Stacked != "" {
		return [][]Resolvable{
			{
				&Value{Value: "write"},
				&Value{Value: strconv.Itoa(m.Position + m.Extra)},
				&Value{Value: m.Stacked},
				&Value{Value: m.Variable},
			},
		}
	}

	return [][]Resolvable{
		{
			&Value{Value: "set"},
			&Value{Value: FunctionTrampolinePrefix + m.Function},
			&Value{Value: strconv.Itoa(m.Position + m.Extra)},
		},
	}
}

func (m *MLOGTrampoline) GetComment(int) string {
	return "Set Trampoline Address"
}
