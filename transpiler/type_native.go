package transpiler

import "strconv"

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

func (m *MLOGTrampolineBack) ToMLOG() [][]Resolvable {
	if m.Stacked != "" {
		return [][]Resolvable{
			{
				&Value{Value: "read"},
				&Value{Value: "@counter"},
				&Value{Value: m.Stacked},
				&Value{Value: stackVariable},
			},
		}
	}

	return [][]Resolvable{
		{
			&Value{Value: "set"},
			&Value{Value: "@counter"},
			&Value{Value: FunctionTrampolinePrefix + m.Function},
		},
	}
}

func (m *MLOGTrampolineBack) GetComment(int) string {
	return "Trampoline back"
}
