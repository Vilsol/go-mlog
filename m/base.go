package m

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

	transpiler.RegisterFuncTranslation("GetHealth", CreateSensorFuncTranslation("@health"))
	transpiler.RegisterFuncTranslation("GetName", CreateSensorFuncTranslation("@name"))
	transpiler.RegisterFuncTranslation("GetX", CreateSensorFuncTranslation("@x"))
	transpiler.RegisterFuncTranslation("GetY", CreateSensorFuncTranslation("@y"))
}

// Read a float64 value from memory at specified position
func Read(memory string, position int) int {
	return 0
}

// Write a float64 value to memory at specified position
//
// For integer equivalent use WriteInt
func Write(value int, memory string, position int) {
}

// Write an integer value to memory at specified position
//
// For float64 equivalent use Write
func WriteInt(value int, memory string, position int) {
}

// Flush all printed statements to the provided message block
func PrintFlush(targetMessage string) {
}

// Get the linked tile at the specified address
func GetLink(address int) Link {
	return nil
}

// Retrieve a list of units that match specified conditions
//
// Conditions are combined using an `and` operation
func Radar(from Building, target1 RadarTarget, target2 RadarTarget, target3 RadarTarget, sortOrder bool, sort RadarSort) Unit {
	return nil
}

// Extract information indicated by sense from the provided block
func Sensor(block HealthC, sense string) float64 {
	return 0
}

func CreateSensorFuncTranslation(attribute string) transpiler.Translator {
	return transpiler.Translator{
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
							&transpiler.Value{Value: attribute},
						},
					},
				},
			}, nil
		},
	}
}
