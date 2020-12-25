package transpiler

import "errors"

func init() {
	RegisterFuncTranslation("print", Translator{
		Count: 1,
		Translate: func(args []Resolvable, _ []Resolvable) ([]MLOGStatement, error) {
			if len(args) == 0 {
				return nil, errors.New("print with 0 arguments")
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
			}, nil
		},
	})
	RegisterFuncTranslation("println", Translator{
		Count: 2,
		Translate: func(args []Resolvable, _ []Resolvable) ([]MLOGStatement, error) {
			if len(args) == 0 {
				return nil, errors.New("println with 0 arguments")
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
			}, nil
		},
	})
}
