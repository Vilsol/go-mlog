package m

import (
	"errors"
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
)

func init() {
	transpiler.RegisterFuncTranslation("m.Read", transpiler.Translator{
		Count:     1,
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			memoryName := strings.Trim(args[0].GetValue(), "\"")

			if memoryName == transpiler.StackCellName {
				return nil, errors.New("can't read/write to memory cell that is used for the stack: " + transpiler.StackCellName)
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
	transpiler.RegisterFuncTranslation("m.Write", transpiler.Translator{
		Count: 1,
		Translate: func(args []transpiler.Resolvable, _ []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			memoryName := strings.Trim(args[1].GetValue(), "\"")

			if memoryName == transpiler.StackCellName {
				return nil, errors.New("can't read/write to memory cell that is used for the stack: " + transpiler.StackCellName)
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
	})
	transpiler.RegisterFuncTranslation("m.PrintFlush", transpiler.Translator{
		Count: 1,
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
		Count:     1,
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
		Count:     1,
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
		Count:     1,
		Variables: 1,
		Translate: func(args []transpiler.Resolvable, vars []transpiler.Resolvable) ([]transpiler.MLOGStatement, error) {
			return []transpiler.MLOGStatement{
				&transpiler.MLOG{
					Statement: [][]transpiler.Resolvable{
						{
							&transpiler.Value{Value: "sensor"},
							vars[0],
							&transpiler.Value{Value: args[0].GetValue()},
							&transpiler.Value{Value: args[1].GetValue()},
						},
					},
				},
			}, nil
		},
	})
}

// TODO Docs
func Read(memory string, position int) int {
	return 0
}

// TODO Docs
func Write(value int, memory string, position int) {
}

// TODO Docs
func PrintFlush(targetMessage string) {
}

// TODO Docs
func GetLink(address int) Link {
	return nil
}

// TODO Docs
func Radar(from string, target1 RadarTarget, target2 RadarTarget, target3 RadarTarget, sortOrder int, sort RadarSort) Unit {
	return nil
}

// TODO Docs
func Sensor(block string, sense string) int {
	return 0
}
