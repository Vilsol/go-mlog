package transpiler

import (
	"errors"
)

func init() {
	RegisterFuncTranslation("print", Translator{
		Count: func(args []Resolvable, vars []Resolvable) int {
			return len(args)
		},
		Translate: func(args []Resolvable, _ []Resolvable) ([]MLOGStatement, error) {
			if len(args) == 0 {
				return nil, errors.New("print with 0 arguments")
			}

			results := make([]MLOGStatement, len(args))
			for i, arg := range args {
				results[i] = &MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "print"},
							&Value{Value: arg.GetValue()},
						},
					},
				}
			}

			return results, nil
		},
	})
	RegisterFuncTranslation("println", Translator{
		Count: func(args []Resolvable, vars []Resolvable) int {
			return len(args) + 1
		},
		Translate: func(args []Resolvable, _ []Resolvable) ([]MLOGStatement, error) {
			if len(args) == 0 {
				return nil, errors.New("println with 0 arguments")
			}

			results := make([]MLOGStatement, len(args)+1)
			for i, arg := range args {
				results[i] = &MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "print"},
							&Value{Value: arg.GetValue()},
						},
					},
				}
			}

			results[len(results)-1] = &MLOG{
				Statement: [][]Resolvable{
					{
						&Value{Value: "print"},
						&Value{Value: `"\n"`},
					},
				},
			}

			return results, nil
		},
	})
	RegisterInlineTranslation("float64", func(args []Resolvable) (Resolvable, error) {
		return args[0], nil
	})
}
