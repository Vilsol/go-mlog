package transpiler

import "strings"

type Translator struct {
	Count     int
	Translate func(args []Resolvable) []MLOGStatement
}

var funcTranslations = map[string]Translator{
	"PrintFlush": {
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "printflush"},
							&Value{Value: strings.Trim(args[0].GetValue(), "\"")},
						},
					},
				},
			}
		},
	},
	"print": {
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
	},
	"println": {
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
	},
	"Write": {
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			memoryName := strings.Trim(args[1].GetValue(), "\"")

			if memoryName == stackCellName {
				panic("can't read/write to memory cell that is used for the stack: " + stackCellName)
			}

			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "write"},
							&Value{Value: args[0].GetValue()},
							&Value{Value: memoryName},
							&Value{Value: args[2].GetValue()},
						},
					},
				},
			}
		},
	},
	"Random": {
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "op"},
							&Value{Value: "rand"},
							&Value{Value: functionReturnVariable},
							&Value{Value: args[0].GetValue()},
						},
					},
				},
			}
		},
	},
	"Floor": {
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "op"},
							&Value{Value: "floor"},
							&Value{Value: functionReturnVariable},
							&Value{Value: args[0].GetValue()},
						},
					},
				},
			}
		},
	},
	"Read": {
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			memoryName := strings.Trim(args[0].GetValue(), "\"")

			if memoryName == stackCellName {
				panic("can't read/write to memory cell that is used for the stack: " + stackCellName)
			}

			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "read"},
							&Value{Value: functionReturnVariable},
							&Value{Value: memoryName},
							&Value{Value: args[1].GetValue()},
						},
					},
				},
			}
		},
	},
	"DrawRect": {
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "draw"},
							&Value{Value: "rect"},
							&Value{Value: args[0].GetValue()},
							&Value{Value: args[1].GetValue()},
							&Value{Value: args[2].GetValue()},
							&Value{Value: args[3].GetValue()},
						},
					},
				},
			}
		},
	},
	"DrawFlush": {
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "drawflush"},
							&Value{Value: strings.Trim(args[0].GetValue(), "\"")},
						},
					},
				},
			}
		},
	},
	"DrawClear": {
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "draw"},
							&Value{Value: "clear"},
							&Value{Value: args[0].GetValue()},
							&Value{Value: args[1].GetValue()},
							&Value{Value: args[2].GetValue()},
						},
					},
				},
			}
		},
	},
	"DrawColor": {
		Count: 1,
		Translate: func(args []Resolvable) []MLOGStatement {
			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "draw"},
							&Value{Value: "color"},
							&Value{Value: args[0].GetValue()},
							&Value{Value: args[1].GetValue()},
							&Value{Value: args[2].GetValue()},
							&Value{Value: args[3].GetValue()},
						},
					},
				},
			}
		},
	},
	"Sleep": {
		Count: 2,
		Translate: func(args []Resolvable) []MLOGStatement {
			dVar := &DynamicVariable{}

			jump := &MLOGJump{
				MLOG: MLOG{},
				Condition: []Resolvable{
					&Value{Value: "lessThan"},
					&Value{Value: "@time"},
					dVar,
				},
			}

			jump.JumpTarget = &StatementJumpTarget{
				Statement: jump,
				After:     false,
			}

			return []MLOGStatement{
				&MLOG{
					Statement: [][]Resolvable{
						{
							&Value{Value: "op"},
							&Value{Value: "add"},
							dVar,
							&Value{Value: "@time"},
							&Value{Value: args[0].GetValue()},
						},
					},
				},
				jump,
			}
		},
	},
}
