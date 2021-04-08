package impl

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
)

func init() {
	transpiler.RegisterInlineTranslation("m.Const", func(args []transpiler.Resolvable) (transpiler.Resolvable, error) {
		return &transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")}, nil
	})
	transpiler.RegisterInlineTranslation("m.B", func(args []transpiler.Resolvable) (transpiler.Resolvable, error) {
		return &transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")}, nil
	})
}
