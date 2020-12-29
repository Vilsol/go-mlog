package transpiler

import (
	"go/token"
	"strconv"
)

type MLOGBranch struct {
	MLOG
	Block *ContextBlock
	Token token.Token
}

func (m *MLOGBranch) ToMLOG() [][]Resolvable {
	lastStatement := m.Block.Statements[len(m.Block.Statements)-1]
	if m.Token != token.CONTINUE && m.Block.Extra != nil && len(m.Block.Extra) > 0 {
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

func (m *MLOGBranch) Size() int {
	return 1
}

func (m *MLOGBranch) GetComment() string {
	return "Branch: " + m.Token.String()
}
