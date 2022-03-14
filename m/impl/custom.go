package impl

import (
	"github.com/Vilsol/go-mlog/transpiler"
)

func init() {
	transpiler.RegisterInlineTranslation("m.Const", func(args []transpiler.Resolvable) (transpiler.Resolvable, error) {
		return &transpiler.InlineVariable{Value: args[0]}, nil
	})
	transpiler.RegisterInlineTranslation("m.B", func(args []transpiler.Resolvable) (transpiler.Resolvable, error) {
		return &transpiler.InlineVariable{Value: args[0]}, nil
	})
}
