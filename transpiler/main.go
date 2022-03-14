package transpiler

import (
	"context"
	"fmt"
	"github.com/Vilsol/go-mlog/checker"
	"github.com/olekukonko/tablewriter"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"strconv"
	"strings"
)

const stackVariable = `@stack`
const FunctionReturnVariable = `@return`

const FunctionArgumentPrefix = `@funcArg_`
const FunctionTrampolinePrefix = `@funcTramp_`

const mainFuncName = `main`

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
			if castDecl.Name.Name == mainFuncName {
				mainFunc = castDecl
			}
		case *ast.GenDecl:
			if castDecl.Tok.String() == "var" {
				return "", Err(context.WithValue(ctx, contextDecl, decl), "global scope may only contain constants not variables")
			} else if castDecl.Tok.String() == "const" {
				constants = append(constants, castDecl)
			}
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

	ctx = context.WithValue(ctx, contextGlobal, global)

	for _, decl := range f.Decls {
		switch castDecl := decl.(type) {
		case *ast.FuncDecl:
			if castDecl.Name.Name == mainFuncName {
				continue
			}

			fnCtx := context.WithValue(ctx, contextFunction, castDecl)
			statements := make([]MLOGStatement, 0)

			prevArgs := 0
			for i, param := range castDecl.Type.Params.List {
				if paramTypeIdent, ok := param.Type.(*ast.Ident); ok {
					if options.Stacked != "" {
						if paramTypeIdent.Name != "int" && paramTypeIdent.Name != "float64" {
							return "", Err(fnCtx, "function parameters may only be integers or floating point numbers in stack mode")
						}
					} else {
						if paramTypeIdent.Name != "int" && paramTypeIdent.Name != "float64" && paramTypeIdent.Name != "string" {
							return "", Err(fnCtx, "function parameters may only be integers, floating point numbers or strings")
						}
					}
				} else {
					return "", Err(fnCtx, "function parameters may only be basic types")
				}

				position := len(castDecl.Type.Params.List) - i

				dVar := &DynamicVariable{}

				if options.Stacked != "" {
					statements = append(statements, &MLOG{
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
					})

					for _, name := range param.Names {
						resolvable := &NormalVariable{Name: name.Name}
						fnCtx = addVariablesToContext(fnCtx, []*VarReference{{
							Name:     name.Name,
							Identity: resolvable,
						}})

						statements = append(statements, &MLOG{
							Comment: "Read parameter into variable",
							Statement: [][]Resolvable{
								{
									&Value{Value: "read"},
									resolvable,
									&Value{Value: options.Stacked},
									dVar,
								},
							},
						})
					}
				} else {
					for j, name := range param.Names {
						resolvable := &NormalVariable{Name: name.Name}
						fnCtx = addVariablesToContext(fnCtx, []*VarReference{{
							Name:     name.Name,
							Identity: resolvable,
						}})

						statements = append(statements, &MLOG{
							Comment: "Read parameter into variable",
							Statement: [][]Resolvable{
								{
									&Value{Value: "set"},
									resolvable,
									&Value{Value: FunctionArgumentPrefix + castDecl.Name.Name + "_" + strconv.Itoa(prevArgs+j)},
								},
							},
						})
					}
				}

				prevArgs += len(param.Names)
			}

			innerStatements, references, err := statementToMLOG(fnCtx, castDecl.Body)
			if err != nil {
				return "", err
			}

			if len(innerStatements) == 0 {
				continue
			}

			statements = append(statements, innerStatements...)

			lastStatement := statements[len(statements)-1]
			if _, ok := lastStatement.(*MLOGTrampolineBack); !ok {
				statements = append(statements, &MLOGTrampolineBack{
					Stacked:  options.Stacked,
					Function: castDecl.Name.Name,
				})
			}

			global.Functions = append(global.Functions, &Function{
				Name:                 castDecl.Name.Name,
				Declaration:          castDecl,
				Statements:           statements,
				ArgumentCount:        len(castDecl.Type.Params.List),
				ScopeVariableCounter: make(map[string]int),
				RootVariables:        references,
			})
		}
	}

	mainStatements, references, err := statementToMLOG(context.WithValue(ctx, contextFunction, mainFunc), mainFunc.Body)

	if err != nil {
		return "", err
	}

	if len(mainStatements) == 0 {
		return "", Err(ctx, "empty main function")
	}

	mainStatements = append(mainStatements, &MLOGTrampolineBack{
		Stacked:  ctx.Value(contextOptions).(Options).Stacked,
		Function: mainFuncName,
	})

	global.Functions = append(global.Functions, &Function{
		Name:                 mainFuncName,
		Called:               true,
		Declaration:          mainFunc,
		Statements:           mainStatements,
		ArgumentCount:        len(mainFunc.Type.Params.List),
		ScopeVariableCounter: make(map[string]int),
		RootVariables:        references,
	})

	var startup []MLOGStatement
	if options.Stacked != "" {
		startup = append(startup, &MLOG{
			Comment: "Reset Stack",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					&Value{Value: stackVariable},
					&Value{Value: "0"},
				},
			},
		})
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
				case *ast.Ident:
					value = valueType.Name
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
					Comment:   "Set global variable",
					SourcePos: spec,
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
			FunctionName: mainFuncName,
		},
	})

	if options.NoStartup {
		startup = make([]MLOGStatement, 0)
	}

	for _, statement := range startup {
		if err := statement.PreProcess(context.WithValue(ctx, contextFunction, mainFunc), global, nil); err != nil {
			return "", err
		}
	}

	for _, fn := range global.Functions {
		for _, statement := range fn.Statements {
			fnCtx := context.WithValue(ctx, contextFunction, fn.Declaration)
			fnCtx = addVariablesToContext(fnCtx, fn.RootVariables)
			if err := statement.PreProcess(fnCtx, global, fn); err != nil {
				return "", err
			}
		}
	}

	position := len(startup)
	for _, fn := range global.Functions {
		if !fn.Called {
			continue
		}

		for _, statement := range fn.Statements {
			if statement.Size() > 0 {
				position += statement.SetPosition(position)
			}
		}
	}

	for _, statement := range startup {
		if err := statement.PostProcess(context.WithValue(ctx, contextFunction, mainFunc), global, nil); err != nil {
			return "", err
		}
	}

	for _, fn := range global.Functions {
		if !fn.Called {
			continue
		}

		for _, statement := range fn.Statements {
			if err := statement.PostProcess(context.WithValue(ctx, contextFunction, fn.Declaration), global, fn); err != nil {
				return "", err
			}
		}
	}

	// Type check at the end to prioritize custom errors
	conf := types.Config{Importer: checker.New()}

	_, err = conf.Check("cmd/hello", fileSet, []*ast.File{f}, nil)
	if err != nil {
		if tError, ok := err.(types.Error); ok {
			return "", Err(context.WithValue(ctx, typeError, tError), tError.Msg)
		}

		return "", Err(ctx, err.Error())
	}

	var tableString *strings.Builder
	var table *tablewriter.Table
	if options.Comments || options.Numbers || options.Source {
		tableString = &strings.Builder{}
		table = tablewriter.NewWriter(tableString)
		table.SetBorder(false)
		table.SetAutoWrapText(false)
		table.SetCenterSeparator("#")
		table.SetColumnSeparator("#")
		table.SetHeaderLine(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetNoWhiteSpace(true)
		table.SetTablePadding("\t")
	}

	outputData := ""

	lineNumber := 0
	for _, statement := range startup {
		statements := statement.ToMLOG()
		mlogLines := MLOGToString(context.WithValue(ctx, contextFunction, mainFunc), statements, statement, lineNumber, input)
		if table != nil {
			table.AppendBulk(mlogLines)
		} else {
			for _, line := range mlogLines {
				outputData += line[0] + "\n"
			}
		}
		if statement.Size() > 0 {
			lineNumber += len(statements)
		}
	}

	for _, fn := range global.Functions {
		if !fn.Called {
			continue
		}

		if options.Comments && table != nil {
			table.Append([]string{"#"})
			table.Append([]string{"# Function: " + fn.Name + " #"})
			table.Append([]string{"#"})
		}

		for _, statement := range fn.Statements {
			statements := statement.ToMLOG()
			mlogLines := MLOGToString(context.WithValue(ctx, contextFunction, fn.Declaration), statements, statement, lineNumber, input)
			if table != nil {
				table.AppendBulk(mlogLines)
			} else {
				for _, line := range mlogLines {
					outputData += line[0] + "\n"
				}
			}
			if statement.Size() > 0 {
				lineNumber += len(statements)
			}
		}
	}

	if table != nil && tableString != nil {
		table.Render()
		return tableString.String(), nil
	}

	return outputData, nil
}
