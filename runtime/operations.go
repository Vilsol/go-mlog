package runtime

import (
	"errors"
	"github.com/rs/zerolog/log"
	"math"
	"math/rand"
)

func init() {
	RegisterOperation("jump", func(args []string) (OperationExecutor, error) {
		if len(args) < 2 {
			return nil, errors.New("expecting at least 2 arguments")
		}

		if args[1] != "always" && len(args) != 4 {
			return nil, errors.New("expecting exactly 4 arguments")
		}

		if args[0][0] == '"' {
			return nil, errors.New("jump target must not be a string")
		}

		return func(ctx *ExecutionContext) {
			log.Debug().Str("op", "jump").Strs("args", args).Msg("executing")

			execute := args[1] == "always"

			switch args[1] {
			case "always":
				break
			case "equal":
				execute = ctx.ResolveStr(args[2]) == ctx.ResolveStr(args[3])
				break
			case "notEqual":
				execute = ctx.ResolveStr(args[2]) != ctx.ResolveStr(args[3])
				break
			case "lessThan":
				execute = ctx.ResolveFloat(args[2]) < ctx.ResolveFloat(args[3])
				break
			case "lessThanEq":
				execute = ctx.ResolveFloat(args[2]) <= ctx.ResolveFloat(args[3])
				break
			case "greaterThan":
				execute = ctx.ResolveFloat(args[2]) > ctx.ResolveFloat(args[3])
				break
			case "greaterThanEq":
				execute = ctx.ResolveFloat(args[2]) >= ctx.ResolveFloat(args[3])
				break
			case "strictEqual":
				execute = ctx.Resolve(args[2]) == ctx.Resolve(args[3])
				break
			default:
				panic("unknown operation: " + args[1])
			}

			if execute {
				ctx.Set("@counter", ctx.ResolveInt(args[0]))
			}
		}, nil
	})

	RegisterOperation("set", func(args []string) (OperationExecutor, error) {
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
	})

	RegisterOperation("op", func(args []string) (OperationExecutor, error) {
		if len(args) < 3 {
			return nil, errors.New("expecting at least 3 arguments")
		}

		if args[0][0] == '"' {
			return nil, errors.New("target variable must not be a string")
		}

		return func(ctx *ExecutionContext) {
			log.Debug().Str("op", "op").Strs("args", args).Msg("executing")

			switch args[0] {
			case "add":
				ctx.Set(args[1], ctx.ResolveFloat(args[2])+ctx.ResolveFloat(args[3]))
				break
			case "sub":
				ctx.Set(args[1], ctx.ResolveFloat(args[2])-ctx.ResolveFloat(args[3]))
				break
			case "mul":
				ctx.Set(args[1], ctx.ResolveFloat(args[2])*ctx.ResolveFloat(args[3]))
				break
			case "div":
				ctx.Set(args[1], ctx.ResolveFloat(args[2])/ctx.ResolveFloat(args[3]))
				break
			case "idiv":
				ctx.Set(args[1], math.Floor(ctx.ResolveFloat(args[2])/ctx.ResolveFloat(args[3])))
				break
			case "mod":
				ctx.Set(args[1], math.Mod(ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
				break
			case "pow":
				ctx.Set(args[1], math.Pow(ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
				break
			case "equal":
				// Special Mindustry logic
				// Math.abs(a - b) < 0.000001 ? 1 : 0, (a, b) -> Structs.eq(a, b) ? 1 : 0)
				if ctx.IsNumber(args[2]) && ctx.IsNumber(args[3]) {
					ctx.Set(args[1], math.Abs(ctx.ResolveFloat(args[2])-ctx.ResolveFloat(args[3])) < 0.000001)
				} else {
					ctx.Set(args[1], ctx.Resolve(args[2]) == ctx.Resolve(args[3]))
				}
				break
			case "notEqual":
				// Special Mindustry logic
				// Math.abs(a - b) < 0.000001 ? 1 : 0, (a, b) -> Structs.eq(a, b) ? 1 : 0)
				if ctx.IsNumber(args[2]) && ctx.IsNumber(args[3]) {
					ctx.Set(args[1], !(math.Abs(ctx.ResolveFloat(args[2])-ctx.ResolveFloat(args[3])) < 0.000001))
				} else {
					ctx.Set(args[1], ctx.Resolve(args[2]) != ctx.Resolve(args[3]))
				}
				break
			case "land":
				ctx.Set(args[1], ctx.ResolveFloat(args[2]) != 0 && ctx.ResolveFloat(args[3]) != 0)
				break
			case "lessThan":
				ctx.Set(args[1], ctx.ResolveFloat(args[2]) < ctx.ResolveFloat(args[3]))
				break
			case "lessThanEq":
				ctx.Set(args[1], ctx.ResolveFloat(args[2]) <= ctx.ResolveFloat(args[3]))
				break
			case "greaterThan":
				ctx.Set(args[1], ctx.ResolveFloat(args[2]) > ctx.ResolveFloat(args[3]))
				break
			case "greaterThanEq":
				ctx.Set(args[1], ctx.ResolveFloat(args[2]) >= ctx.ResolveFloat(args[3]))
				break
			case "strictEqual":
				ctx.Set(args[1], ctx.Resolve(args[2]) == ctx.Resolve(args[3]))
				break
			case "shl":
				ctx.Set(args[1], ctx.ResolveInt(args[2])<<ctx.ResolveInt(args[3]))
				break
			case "shr":
				ctx.Set(args[1], ctx.ResolveInt(args[2])>>ctx.ResolveInt(args[3]))
				break
			case "or":
				ctx.Set(args[1], ctx.ResolveInt(args[2])|ctx.ResolveInt(args[3]))
				break
			case "and":
				ctx.Set(args[1], ctx.ResolveInt(args[2])&ctx.ResolveInt(args[3]))
				break
			case "xor":
				ctx.Set(args[1], ctx.ResolveInt(args[2])^ctx.ResolveInt(args[3]))
				break
			case "not":
				ctx.Set(args[1], ^ctx.ResolveInt(args[2]))
				break
			case "max":
				ctx.Set(args[1], math.Max(ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
				break
			case "min":
				ctx.Set(args[1], math.Min(ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
				break
			case "angle":
				x := ctx.ResolveFloat(args[2])
				y := ctx.ResolveFloat(args[3])

				ang := math.Atan2(x, y) * (float64(180) / math.Pi)

				if ang < 0 {
					ang += 360
				}

				ctx.Set(args[1], ang)
				break
			case "len":
				x := ctx.ResolveFloat(args[2])
				y := ctx.ResolveFloat(args[3])
				ctx.Set(args[1], math.Sqrt(x*x+y*y))
				break
			case "noise":
				ctx.Set(args[1], raw2d(0, ctx.ResolveFloat(args[2]), ctx.ResolveFloat(args[3])))
				break
			case "abs":
				ctx.Set(args[1], math.Abs(ctx.ResolveFloat(args[2])))
				break
			case "log":
				ctx.Set(args[1], math.Log(ctx.ResolveFloat(args[2])))
				break
			case "log10":
				ctx.Set(args[1], math.Log10(ctx.ResolveFloat(args[2])))
				break
			case "sin":
				ctx.Set(args[1], math.Sin(ctx.ResolveFloat(args[2])*0.017453292519943295))
				break
			case "cos":
				ctx.Set(args[1], math.Cos(ctx.ResolveFloat(args[2])*0.017453292519943295))
				break
			case "tan":
				ctx.Set(args[1], math.Tan(ctx.ResolveFloat(args[2])*0.017453292519943295))
				break
			case "floor":
				ctx.Set(args[1], math.Floor(ctx.ResolveFloat(args[2])))
				break
			case "ceil":
				ctx.Set(args[1], math.Ceil(ctx.ResolveFloat(args[2])))
				break
			case "sqrt":
				ctx.Set(args[1], math.Sqrt(ctx.ResolveFloat(args[2])))
				break
			case "rand":
				ctx.Set(args[1], rand.Float64()*ctx.ResolveFloat(args[2]))
				break
			default:
				panic("unknown operation: " + args[0])
			}
		}, nil
	})

	RegisterOperation("print", func(args []string) (OperationExecutor, error) {
		if len(args) != 1 {
			return nil, errors.New("expecting exactly 1 argument")
		}

		return func(ctx *ExecutionContext) {
			log.Debug().Str("op", "print").Strs("args", args).Msg("executing")

			ctx.PrintBuffer.WriteString(ctx.ResolveStr(args[0]))
		}, nil
	})

	RegisterOperation("printflush", func(args []string) (OperationExecutor, error) {
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
	})

	RegisterOperation("end", func(args []string) (OperationExecutor, error) {
		if len(args) != 0 {
			return nil, errors.New("requires no arguments")
		}

		return func(ctx *ExecutionContext) {
			log.Debug().Str("op", "end").Strs("args", args).Msg("executing")

			// TODO Jump to 0th instruction
			ctx.Set("@counter", int64(1001))
		}, nil
	})

	RegisterOperation("drawflush", func(args []string) (OperationExecutor, error) {
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
	})

	RegisterOperation("draw", func(args []string) (OperationExecutor, error) {
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
		case string(DrawActionStroke):
			if len(args) < 2 {
				return nil, errors.New("expecting at least 2 arguments")
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
	})
}
