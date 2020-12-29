package transpiler

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strconv"
)

const stackVariable = `@stack`
const FunctionReturnVariable = `@return`

const StackCellName = `bank1`
const debugCellName = `cell2`
const debugCount = 2

func GolangToMLOGFile(fileName string, options Options) (string, error) {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return GolangToMLOGBytes(file, options)
}

func GolangToMLOGBytes(input []byte, options Options) (string, error) {
	return GolangToMLOG(string(input), options)
}

func GolangToMLOG(input string, options Options) (string, error) {
	ctx := context.WithValue(context.Background(), contextOptions, options)

	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "foo", input, 0)

	if err != nil {
		return "", err
	}

	if f.Name.Name != "main" {
		return "", Err(ctx, "package must be main")
	}

	for _, imp := range f.Imports {
		if _, ok := validImports[imp.Path.Value]; !ok {
			return "", Err(context.WithValue(ctx, contextSpec, imp), "unregistered import used: "+imp.Path.Value)
		}
	}

	constants := make([]*ast.GenDecl, 0)
	var mainFunc *ast.FuncDecl
	for _, decl := range f.Decls {
		switch castDecl := decl.(type) {
		case *ast.FuncDecl:
			if castDecl.Name.Name == "main" {
				mainFunc = castDecl
			}
			break
		case *ast.GenDecl:
			if castDecl.Tok.String() == "var" {
				return "", Err(context.WithValue(ctx, contextDecl, decl), "global scope may only contain constants not variables")
			} else if castDecl.Tok.String() == "const" {
				constants = append(constants, castDecl)
			}
			break
		case *ast.BadDecl:
			return "", Err(ctx, "syntax error in input file")
		}
	}

	if mainFunc == nil {
		return "", Err(ctx, "file does not contain a main function")
	}

	global := &Global{
		Functions: make([]*Function, 0),
	}

	for _, decl := range f.Decls {
		switch castDecl := decl.(type) {
		case *ast.FuncDecl:
			if castDecl.Name.Name == "main" {
				continue
			}
			fnCtx := context.WithValue(ctx, contextFunction, castDecl)
			statements, err := statementToMLOG(fnCtx, castDecl.Body)
			if err != nil {
				return "", err
			}

			for i, param := range castDecl.Type.Params.List {
				if paramTypeIdent, ok := param.Type.(*ast.Ident); ok {
					if paramTypeIdent.Name != "int" && paramTypeIdent.Name != "float64" {
						return "", Err(fnCtx, "function parameters may only be integers or floating point numbers")
					}
				} else {
					return "", Err(fnCtx, "function parameters may only be integers or floating point numbers")
				}

				position := len(castDecl.Type.Params.List) - i

				dVar := &DynamicVariable{}

				for _, name := range param.Names {
					statements = append([]MLOGStatement{&MLOG{
						Comment: "Read parameter into variable",
						Statement: [][]Resolvable{
							{
								&Value{Value: "read"},
								&NormalVariable{Name: name.Name},
								&Value{Value: StackCellName},
								dVar,
							},
						},
					}}, statements...)
				}

				statements = append([]MLOGStatement{
					&MLOG{
						Comment: "Calculate address of parameter",
						Statement: [][]Resolvable{
							{
								&Value{Value: "op"},
								&Value{Value: "sub"},
								dVar,
								&Value{Value: stackVariable},
								&Value{Value: strconv.Itoa(position)},
							},
						},
					},
				}, statements...)
			}

			lastStatement := statements[len(statements)-1]
			if _, ok := lastStatement.(*MLOGTrampolineBack); !ok {
				statements = append(statements, &MLOGTrampolineBack{})
			}

			global.Functions = append(global.Functions, &Function{
				Name:          castDecl.Name.Name,
				Declaration:   castDecl,
				Statements:    statements,
				ArgumentCount: len(castDecl.Type.Params.List),
			})
			break
		}
	}

	mainStatements, err := statementToMLOG(context.WithValue(ctx, contextFunction, mainFunc), mainFunc.Body)

	if err != nil {
		return "", err
	}

	if len(mainStatements) == 0 {
		return "", Err(ctx, "empty main function")
	}

	global.Functions = append(global.Functions, &Function{
		Name:          mainFunc.Name.Name,
		Declaration:   mainFunc,
		Statements:    mainStatements,
		ArgumentCount: len(mainFunc.Type.Params.List),
	})

	startup := []MLOGStatement{
		&MLOG{
			Comment: "Reset Stack",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					&Value{Value: stackVariable},
					&Value{Value: "0"},
				},
			},
		},
	}

	global.Constants = make(map[string]bool)
	constantPos := 0
	for _, constant := range constants {
		for _, spec := range constant.Specs {
			// Constants can only be ValueSpec
			valueSpec := spec.(*ast.ValueSpec)
			for i, name := range valueSpec.Names {
				var value string
				switch valueType := valueSpec.Values[i].(type) {
				case *ast.BasicLit:
					value = valueType.Value
					break
				case *ast.Ident:
					value = valueType.Name
					break
				default:
					return "", Err(context.WithValue(ctx, contextSpec, spec), fmt.Sprintf("unknown constant type: %T", valueSpec.Values[i]))
				}

				startup = append(startup, &MLOG{
					Position: constantPos,
					Statement: [][]Resolvable{
						{
							&Value{Value: "set"},
							&Value{Value: name.Name},
							&Value{Value: value},
						},
					},
					Comment: "Set global variable",
				})
				constantPos += 1

				global.Constants[name.Name] = true
			}
		}
	}

	startup = append(startup, &MLOGJump{
		MLOG: MLOG{
			Comment:  "Jump to start of main",
			Position: len(startup),
		},
		Condition: []Resolvable{
			&Value{Value: "always"},
		},
		JumpTarget: &FunctionJumpTarget{
			FunctionName: "main",
		},
	})

	if options.NoStartup {
		startup = make([]MLOGStatement, 0)
	}

	debugWriter := []MLOGAble{
		&MLOG{
			Comment: "Debug",
			Statement: [][]Resolvable{
				{
					&Value{Value: "write"},
					&Value{Value: "@counter"},
					&Value{Value: debugCellName},
					&Value{Value: "0"},
				},
			},
		},
		&MLOG{
			Comment: "Debug",
			Statement: [][]Resolvable{
				{
					&Value{Value: "write"},
					&Value{Value: stackVariable},
					&Value{Value: debugCellName},
					&Value{Value: "1"},
				},
			},
		},
	}

	if len(debugWriter) != debugCount {
		panic("debugWriter count != debugCount")
	}

	position := len(startup)
	for _, fn := range global.Functions {
		for _, statement := range fn.Statements {
			if options.Debug {
				position += debugCount
			}

			position += statement.SetPosition(position)
		}
	}

	for _, statement := range startup {
		if err := statement.PostProcess(context.WithValue(ctx, contextFunction, mainFunc), global, nil); err != nil {
			return "", err
		}
	}

	for _, fn := range global.Functions {
		for _, statement := range fn.Statements {
			if err := statement.PostProcess(context.WithValue(ctx, contextFunction, fn.Declaration), global, fn); err != nil {
				return "", err
			}
		}
	}

	result := ""
	lineNumber := 0
	for _, statement := range startup {
		statements := statement.ToMLOG()
		result += MLOGToString(context.WithValue(ctx, contextFunction, mainFunc), statements, statement, lineNumber)
		lineNumber += len(statements)
	}

	for _, fn := range global.Functions {
		if options.Comments {
			if result != "" {
				result += "\n"
			}

			result += "     // Function: " + fn.Name + " //\n"
		}

		for _, statement := range fn.Statements {
			if options.Debug {
				for _, debugStatement := range debugWriter {
					deb := debugStatement.ToMLOG()
					result += MLOGToString(context.WithValue(ctx, contextFunction, fn.Declaration), deb, debugStatement, lineNumber)
					lineNumber += len(deb)
				}
			}

			statements := statement.ToMLOG()
			result += MLOGToString(context.WithValue(ctx, contextFunction, fn.Declaration), statements, statement, lineNumber)
			lineNumber += len(statements)
		}
	}

	return result, nil
}
