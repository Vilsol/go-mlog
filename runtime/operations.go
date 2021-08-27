package runtime

import (
	"errors"
	"github.com/rs/zerolog/log"
	"math"
	"math/rand"
	"time"
)

func init() {
	RegisterOperation("jump", jumpOperation)
	RegisterOperation("end", endOperation)
	RegisterOperation("wait", waitOperation)

	RegisterOperation("set", setOperation)
	RegisterOperation("op", opOperation)

	RegisterOperation("print", printOperation)
	RegisterOperation("printflush", printFlushOperation)

	RegisterOperation("draw", drawOperation)
	RegisterOperation("drawflush", drawFlushOperation)

	RegisterOperation("read", readOperation)
	RegisterOperation("write", writeOperation)

	// TODO Get Link
	// TODO Control
	// TODO Radar
	// TODO Sensor

	// TODO Lookup

	// TODO Unit Bind
	// TODO Unit Control
	// TODO Unit Radar
	// TODO Unit Locate
}

func jumpOperation(args []string) (OperationExecutor, error) {
	if len(args) < 2 {
		return nil, errors.New("expecting at least 2 arguments")
	}

	if args[1] != "always" && len(args) != 4 {
		return nil, errors.New("expecting exactly 4 arguments")
	}

	if args[0][0] == '"' {
		return nil, errors.New("jump target must not be a string")
	}

	logLine := log.Debug().Str("op", "jump").Strs("args", args)

	common := func(execute bool, ctx *ExecutionContext) {
		if execute {
			ctx.Set("@counter", ctx.ResolveInt(args[0]))
		}
	}

	switch args[1] {
	case "always":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			common(true, ctx)
		}, nil
	case "equal":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			common(ctx.ResolveStr(args[2]) == ctx.ResolveStr(args[3]), ctx)
		}, nil
	case "notEqual":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			common(ctx.ResolveStr(args[2]) != ctx.ResolveStr(args[3]), ctx)
		}, nil
	case "lessThan":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			common(ctx.ResolveFloat(args[2]) < ctx.ResolveFloat(args[3]), ctx)
		}, nil
	case "lessThanEq":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			common(ctx.ResolveFloat(args[2]) <= ctx.ResolveFloat(args[3]), ctx)
		}, nil
	case "greaterThan":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			common(ctx.ResolveFloat(args[2]) > ctx.ResolveFloat(args[3]), ctx)
		}, nil
	case "greaterThanEq":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			common(ctx.ResolveFloat(args[2]) >= ctx.ResolveFloat(args[3]), ctx)
		}, nil
	case "strictEqual":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			common(ctx.Resolve(args[2]) == ctx.Resolve(args[3]), ctx)
		}, nil
	}

	return nil, errors.New("unknown operation: " + args[1])
}

func setOperation(args []string) (OperationExecutor, error) {
	if len(args) != 2 {
		return nil, errors.New("expecting exactly 2 arguments")
	}

	if args[0][0] == '"' {
		return nil, errors.New("target variable must not be a string")
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "set").Strs("args", args).Msg("executing")

		ctx.Set(args[0], ctx.Resolve(args[1]))
	}, nil
}

func opOperation(args []string) (OperationExecutor, error) {
	if len(args) < 2 {
		return nil, errors.New("expecting at least 2 arguments")
	}

	if args[0][0] == '"' {
		return nil, errors.New("target variable must not be a string")
	}

	logLine := log.Debug().Str("op", "op").Strs("args", args)

	switch args[0] {
	case "add":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveFloat(args[2])+ctx.ResolveFloat(args[3]))
		}, nil
	case "sub":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveFloat(args[2])-ctx.ResolveFloat(args[3]))
		}, nil
	case "mul":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveFloat(args[2])*ctx.ResolveFloat(args[3]))
		}, nil
	case "div":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveFloat(args[2])/ctx.ResolveFloat(args[3]))
		}, nil
	case "idiv":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Floor(ctx.ResolveFloat(args[2])/ctx.ResolveFloat(args[3])))
		}, nil
	case "mod":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Mod(ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
		}, nil
	case "pow":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Pow(ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
		}, nil
	case "equal":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			// Special Mindustry logic
			// Math.abs(a - b) < 0.000001 ? 1 : 0, (a, b) -> Structs.eq(a, b) ? 1 : 0)
			if ctx.IsNumber(args[2]) && ctx.IsNumber(args[3]) {
				ctx.Set(args[1], math.Abs(ctx.ResolveFloat(args[2])-ctx.ResolveFloat(args[3])) < 0.000001)
			} else {
				ctx.Set(args[1], ctx.Resolve(args[2]) == ctx.Resolve(args[3]))
			}
		}, nil
	case "notEqual":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			// Special Mindustry logic
			// Math.abs(a - b) < 0.000001 ? 1 : 0, (a, b) -> Structs.eq(a, b) ? 1 : 0)
			if ctx.IsNumber(args[2]) && ctx.IsNumber(args[3]) {
				ctx.Set(args[1], !(math.Abs(ctx.ResolveFloat(args[2])-ctx.ResolveFloat(args[3])) < 0.000001))
			} else {
				ctx.Set(args[1], ctx.Resolve(args[2]) != ctx.Resolve(args[3]))
			}
		}, nil
	case "land":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveFloat(args[2]) != 0 && ctx.ResolveFloat(args[3]) != 0)
		}, nil
	case "lessThan":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveFloat(args[2]) < ctx.ResolveFloat(args[3]))
		}, nil
	case "lessThanEq":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveFloat(args[2]) <= ctx.ResolveFloat(args[3]))
		}, nil
	case "greaterThan":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveFloat(args[2]) > ctx.ResolveFloat(args[3]))
		}, nil
	case "greaterThanEq":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveFloat(args[2]) >= ctx.ResolveFloat(args[3]))
		}, nil
	case "strictEqual":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.Resolve(args[2]) == ctx.Resolve(args[3]))
		}, nil
	case "shl":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveInt(args[2])<<ctx.ResolveInt(args[3]))
		}, nil
	case "shr":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveInt(args[2])>>ctx.ResolveInt(args[3]))
		}, nil
	case "or":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveInt(args[2])|ctx.ResolveInt(args[3]))
		}, nil
	case "and":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveInt(args[2])&ctx.ResolveInt(args[3]))
		}, nil
	case "xor":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ctx.ResolveInt(args[2])^ctx.ResolveInt(args[3]))
		}, nil
	case "not":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], ^ctx.ResolveInt(args[2]))
		}, nil
	case "max":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Max(ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
		}, nil
	case "min":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Min(ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
		}, nil
	case "angle":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			x := ctx.ResolveFloat(args[2])
			y := ctx.ResolveFloat(args[3])

			ang := math.Atan2(x, y) * (float64(180) / math.Pi)

			if ang < 0 {
				ang += 360
			}

			ctx.Set(args[1], ang)
		}, nil
	case "len":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			x := ctx.ResolveFloat(args[2])
			y := ctx.ResolveFloat(args[3])
			ctx.Set(args[1], math.Sqrt(x*x+y*y))
		}, nil
	case "noise":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], raw2d(0, ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
		}, nil
	case "abs":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Abs(ctx.ResolveFloat(args[2])))
		}, nil
	case "log":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Log(ctx.ResolveFloat(args[2])))
		}, nil
	case "log10":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Log10(ctx.ResolveFloat(args[2])))
		}, nil
	case "sin":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Sin(ctx.ResolveFloat(args[2])*0.017453292519943295))
		}, nil
	case "cos":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Cos(ctx.ResolveFloat(args[2])*0.017453292519943295))
		}, nil
	case "tan":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Tan(ctx.ResolveFloat(args[2])*0.017453292519943295))
		}, nil
	case "floor":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Floor(ctx.ResolveFloat(args[2])))
		}, nil
	case "ceil":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Ceil(ctx.ResolveFloat(args[2])))
		}, nil
	case "sqrt":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], math.Sqrt(ctx.ResolveFloat(args[2])))
		}, nil
	case "rand":
		return func(ctx *ExecutionContext) {
			logLine.Msg("executing")
			ctx.Set(args[1], rand.Float64()*ctx.ResolveFloat(args[2]))
		}, nil
	default:
		panic("unknown operation: " + args[0])
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "op").Strs("args", args).Msg("executing")

	}, nil
}

func printOperation(args []string) (OperationExecutor, error) {
	if len(args) != 1 {
		return nil, errors.New("expecting exactly 1 argument")
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "print").Strs("args", args).Msg("executing")

		ctx.PrintBuffer.WriteString(ctx.ResolveStr(args[0]))
	}, nil
}

func printFlushOperation(args []string) (OperationExecutor, error) {
	if len(args) != 1 {
		return nil, errors.New("expecting exactly 1 argument")
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "printflush").Strs("args", args).Msg("executing")

		message, err := ctx.Message(args[0])

		if err != nil {
			panic(err)
		}

		buffer := ctx.PrintBuffer.String()
		ctx.PrintBuffer.Reset()

		message.PrintFlush(buffer)
	}, nil
}

func endOperation(args []string) (OperationExecutor, error) {
	if len(args) != 0 {
		return nil, errors.New("requires no arguments")
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "end").Strs("args", args).Msg("executing")

		// TODO Jump to 0th instruction
		ctx.Set("@counter", int64(1001))
	}, nil
}

func drawFlushOperation(args []string) (OperationExecutor, error) {
	if len(args) != 1 {
		return nil, errors.New("expecting exactly 1 argument")
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "drawflush").Strs("args", args).Msg("executing")

		display, err := ctx.Display(args[0])
		if err != nil {
			panic(err)
		}

		buffer := ctx.DrawBuffer
		ctx.DrawBuffer = make([]DrawStatement, 0)

		display.DrawFlush(buffer)
	}, nil
}

func drawOperation(args []string) (OperationExecutor, error) {
	if len(args) < 2 {
		return nil, errors.New("expecting at least 2 arguments")
	}

	switch args[0] {
	case string(DrawActionClear):
		if len(args) < 4 {
			return nil, errors.New("expecting at least 4 arguments")
		}
		break
	case string(DrawActionColor):
		if len(args) < 5 {
			return nil, errors.New("expecting at least 5 arguments")
		}
		break
	case string(DrawActionLine):
		if len(args) < 5 {
			return nil, errors.New("expecting at least 5 arguments")
		}
		break
	case string(DrawActionRect):
		if len(args) < 5 {
			return nil, errors.New("expecting at least 5 arguments")
		}
		break
	case string(DrawActionLineRect):
		if len(args) < 5 {
			return nil, errors.New("expecting at least 5 arguments")
		}
		break
	case string(DrawActionPoly):
		if len(args) < 6 {
			return nil, errors.New("expecting at least 6 arguments")
		}
		break
	case string(DrawActionLinePoly):
		if len(args) < 6 {
			return nil, errors.New("expecting at least 6 arguments")
		}
		break
	case string(DrawActionTriangle):
		if len(args) < 7 {
			return nil, errors.New("expecting at least 7 arguments")
		}
		break
	case string(DrawActionImage):
		if len(args) < 6 {
			return nil, errors.New("expecting at least 6 arguments")
		}
		break
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "draw").Strs("args", args).Msg("executing")

		resolved := make([]interface{}, len(args)-1)

		for i, arg := range args[1:] {
			resolved[i] = ctx.Resolve(arg)
		}

		ctx.DrawBuffer = append(ctx.DrawBuffer, DrawStatement{
			Action:    DrawAction(args[0]),
			Arguments: resolved,
		})
	}, nil
}

func waitOperation(args []string) (OperationExecutor, error) {
	if len(args) != 1 {
		return nil, errors.New("expecting exactly 1 argument")
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "wait").Strs("args", args).Msg("executing")

		time.Sleep(time.Duration(int(float64(time.Second) * ctx.ResolveFloat(args[0]))))
	}, nil
}

func readOperation(args []string) (OperationExecutor, error) {
	if len(args) != 3 {
		return nil, errors.New("expecting exactly 3 arguments")
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "read").Strs("args", args).Msg("executing")

		memory, err := ctx.Memory(args[1])
		if err != nil {
			panic(err)
		}

		ctx.Set(args[0], memory.Read(ctx.ResolveInt(args[2])))
	}, nil
}

func writeOperation(args []string) (OperationExecutor, error) {
	if len(args) != 3 {
		return nil, errors.New("expecting exactly 3 arguments")
	}

	return func(ctx *ExecutionContext) {
		log.Debug().Str("op", "write").Strs("args", args).Msg("executing")

		memory, err := ctx.Memory(args[1])
		if err != nil {
			panic(err)
		}

		memory.Write(ctx.ResolveFloat(args[0]), ctx.ResolveInt(args[2]))
	}, nil
}
