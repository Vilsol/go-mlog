package transpiler

func init() {
	RegisterFuncTranslation("print", Translator{
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			if len(args) == 0 {
				panic("print with 0 arguments")
			}
			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "print"},
							&Value{Value: args[0].GetValue()},
						},
					},
				},
			}
		},
	})
	RegisterFuncTranslation("println", Translator{
		Count: 2,
		Translate: func(args []Resolvable) []MLOGStatement {
			if len(args) == 0 {
				panic("println with 0 arguments")
			}
			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "print"},
							&Value{Value: args[0].GetValue()},
						},
					},
				},
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "print"},
							&Value{Value: `"\n"`},
						},
					},
				},
			}
		},
	})
}
