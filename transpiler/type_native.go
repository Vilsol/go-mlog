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

func (m *MLOGStackWriter) GetComment() string {
	return "Update Stack Pointer"
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
