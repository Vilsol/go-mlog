package impl

import (
	"errors"
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
)

func init() {
	transpiler.RegisterFuncTranslation("m.Read", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			memoryName := strings.Trim(args[0].GetValue(), "\"")

			// TODO Remove hardcode
			if memoryName == "bank1" {
				return nil, errors.New("can't read/write to memory cell that is used for the stack: bank1")
			}

			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "read"},
							vars[0],
							&transpiler.Value{Value: memoryName},
							&transpiler.Value{Value: args[1].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	write := transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			memoryName := strings.Trim(args[1].GetValue(), "\"")

			if memoryName == "bank1" {
				return nil, errors.New("can't read/write to memory cell that is used for the stack: bank1")
			}

			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "write"},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: memoryName},
							&transpiler.Value{Value: args[2].GetValue()},
						},
					},
				},
			}, nil
		},
	}
	transpiler.RegisterFuncTranslation("m.Write", write)
	transpiler.RegisterFuncTranslation("m.WriteInt", write)
	transpiler.RegisterFuncTranslation("m.PrintFlush", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "printflush"},
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.GetLink", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "getlink"},
							vars[0],
							&transpiler.Value{Value: args[0].GetValue()},
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.Radar", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "radar"},
							&transpiler.Value{Value: args[1].GetValue()},
							&transpiler.Value{Value: args[2].GetValue()},
							&transpiler.Value{Value: args[3].GetValue()},
							&transpiler.Value{Value: args[5].GetValue()},
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[4].GetValue()},
							vars[0],
						},
					},
				},
			}, nil
		},
	})
	transpiler.RegisterFuncTranslation("m.Sensor", transpiler.Translator{
		Count: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) int {
			return 1
		},
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "sensor"},
							vars[0],
							&transpiler.Value{Value: strings.Trim(args[0].GetValue(), "\"")},
							&transpiler.Value{Value: strings.Trim(args[1].GetValue(), "\"")},
						},
					},
				},
			}, nil
		},
	})
}
