package transpiler

import (
	"context"
	"strconv"
)

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

func (m *MLOGStackWriter) GetComment(int) string {
	return "Update Stack Pointer"
}

type MLOGTrampolineBack struct {
	MLOG
	Stacked  string
	Function string
}

func (m *MLOGTrampolineBack) PreProcess(ctx context.Context, global *Global, function *Function) error {
	if m.Stacked != "" {
		m.Statement = [][]Resolvable{
			{
				&Value{Value: "read"},
				&Value{Value: "@counter"},
				&Value{Value: m.Stacked},
				&Value{Value: stackVariable},
			},
		}
	} else {
		m.Statement = [][]Resolvable{
			{
				&Value{Value: "set"},
				&Value{Value: "@counter"},
				&Value{Value: FunctionTrampolinePrefix + m.Function},
			},
		}
	}

	return m.MLOG.PreProcess(ctx, global, function)
}

func (m *MLOGTrampolineBack) GetComment(int) string {
	return "Trampoline back"
}
