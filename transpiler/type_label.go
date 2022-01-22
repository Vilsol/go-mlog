package transpiler

type MLOGLabel struct {
	MLOG
	Name string
}

func (m *MLOGLabel) ToMLOG() [][]Resolvable {
	return [][]Resolvable{
		{
			&Value{Value: m.Name + ":"},
		},
	}
}

func (m *MLOGLabel) Size() int {
	return 0
}

func (m *MLOGLabel) GetComment(int) string {
	return "Add a label"
}
