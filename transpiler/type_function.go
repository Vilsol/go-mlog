package transpiler

import (
	"context"
	"go/ast"
	"go/token"
	"strconv"
)

type MLOGCustomFunction struct {
	Position     int
	Arguments    []ast.Expr
	Variables    []Resolvable
	Unresolved   []MLOGStatement
	FunctionName string
	Comments     map[int]string
	SourcePos    token.Pos
}

func (m *MLOGCustomFunction) ToMLOG() [][]Resolvable {
	results := make([][]Resolvable, 0)
	m.Comments = make(map[int]string)
	for _, statement := range m.Unresolved {
		lines := statement.ToMLOG()
		results = append(results, lines...)
		for i := statement.GetPosition(); i < statement.GetPosition()+len(lines); i++ {
			m.Comments[i] = statement.GetComment(i)
		}
	}
	return results
}

func (m *MLOGCustomFunction) GetPosition() int {
	return m.Position
}

func (m *MLOGCustomFunction) Size() int {
	size := 0
	for _, statement := range m.Unresolved {
		size += statement.Size()
	}
	return size
}

func (m *MLOGCustomFunction) SetPosition(position int) int {
	m.Position = position
	size := 0
	for _, statement := range m.Unresolved {
		size += statement.SetPosition(size + m.Position)
	}
	return size
}

func (m *MLOGCustomFunction) PreProcess(ctx context.Context, global *Global, function *Function) error {
	if len(m.Unresolved) > 0 {
		return nil
	}

	stacked := ctx.Value(contextOptions).(Options).Stacked

	for i, arg := range m.Arguments {
		value, argInstructions, err := exprToResolvable(ctx, arg)
		if err != nil {
			return err
		}
		m.Unresolved = append(m.Unresolved, argInstructions...)

		if stacked != "" {
			m.Unresolved = append(m.Unresolved, &MLOGStackWriter{
				Action: "add",
			})

			m.Unresolved = append(m.Unresolved, &MLOG{
				Comment: "Write argument to memory",
				Statement: [][]Resolvable{
					{
						&Value{Value: "write"},
						value,
						&Value{Value: stacked},
						&Value{Value: stackVariable},
					},
				},
			})
		} else {
			argNum := strconv.Itoa(i)
			m.Unresolved = append(m.Unresolved, &MLOG{
				Comment: "Set " + m.FunctionName + " argument: " + argNum,
				Statement: [][]Resolvable{
					{
						&Value{Value: "set"},
						&Value{Value: FunctionArgumentPrefix + m.FunctionName + "_" + argNum},
						value,
					},
				},
			})
		}
	}

	if stacked != "" {
		m.Unresolved = append(m.Unresolved, &MLOGStackWriter{
			Action: "add",
		})
	}

	m.Unresolved = append(m.Unresolved, &MLOGTrampoline{
		Variable: stackVariable,
		Extra:    2,
		Stacked:  stacked,
		Function: m.FunctionName,
	})

	m.Unresolved = append(m.Unresolved, &MLOGJump{
		MLOG: MLOG{
			Comment: "Jump to function: " + m.FunctionName,
		},
		Condition: []Resolvable{
			&Value{Value: "always"},
		},
		JumpTarget: &FunctionJumpTarget{
			FunctionName: m.FunctionName,
			SourcePos:    m.SourcePos,
		},
	})

	if stacked != "" {
		m.Unresolved = append(m.Unresolved, &MLOGStackWriter{
			Action: "sub",
			Extra:  len(m.Arguments),
		})
	}

	if len(m.Variables) > 0 {
		m.Unresolved = append(m.Unresolved, &MLOG{
			Comment: "Set variable to returned value",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					m.Variables[0],
					&Value{Value: FunctionReturnVariable},
				},
			},
		})
	}

	for _, statement := range m.Unresolved {
		if err := statement.PreProcess(ctx, global, function); err != nil {
			return err
		}
	}

	return nil
}
func (m *MLOGCustomFunction) PostProcess(ctx context.Context, global *Global, function *Function) error {
	for i, statement := range m.Unresolved {
		statement.SetPosition(m.Position + i)
		if err := statement.PostProcess(ctx, global, function); err != nil {
			return err
		}
	}
	return nil
}

func (m *MLOGCustomFunction) GetComment(pos int) string {
	return m.Comments[pos]
}

func (m *MLOGCustomFunction) SetSourcePos(pos token.Pos) {
	m.SourcePos = pos
}

func (m *MLOGCustomFunction) GetSourcePos() token.Pos {
	return m.SourcePos
}
