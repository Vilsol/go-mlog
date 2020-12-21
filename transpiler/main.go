package transpiler

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strconv"
)

const stackVariable = `@stack`
const functionReturnVariable = `@return`

const stackCellName = `bank1`
const debugCellName = `cell2`
const debugCount = 2

// TODO Change to List, Support Math
const validImport = `"github.com/Vilsol/go-mlog/m"`

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
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "foo", input, 0)

	if err != nil {
		return "", err
	}

	if f.Name.Name != "main" {
		return "", errors.New("package must be main")
	}

	for _, imp := range f.Imports {
		if imp.Path.Value != validImport {
			return "", errors.New("you may not use any external imports")
		}
	}

	constants := make([]*ast.GenDecl, 0)
	var mainFunc *ast.FuncDecl
	for _, decl := range f.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			funcDecl := decl.(*ast.FuncDecl)
			if funcDecl.Name.Name == "main" {
				mainFunc = funcDecl
			}
			break
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl.Tok.String() == "var" {
				return "", errors.New("global scope may only contain constants not variables")
			} else if genDecl.Tok.String() == "const" {
				constants = append(constants, genDecl)
			}
			break
		case *ast.BadDecl:
			return "", errors.New("syntax error in input file")
		}
	}

	if mainFunc == nil {
		return "", errors.New("file does not contain a main function")
	}

	global := &Global{
		Functions: make([]*Function, 0),
	}

	for _, decl := range f.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			funcDecl := decl.(*ast.FuncDecl)
			if funcDecl.Name.Name == "main" {
				continue
			}
			statements, err := statementToMLOG(funcDecl.Body, options)
			if err != nil {
				return "", err
			}

			for i, param := range funcDecl.Type.Params.List {
				position := len(funcDecl.Type.Params.List) - i

				dVar := &DynamicVariable{}

				// TODO Support multiple names
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
					&MLOG{
						Comment: "Read parameter into variable",
						Statement: [][]Resolvable{
							{
								&Value{Value: "read"},
								&NormalVariable{Name: param.Names[0].Name},
								&Value{Value: stackCellName},
								dVar,
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
				Name:          funcDecl.Name.Name,
				Statements:    statements,
				ArgumentCount: len(funcDecl.Type.Params.List),
			})
			break
		}
	}

	mainStatements, err := statementToMLOG(mainFunc.Body, options)
	if err != nil {
		return "", err
	}
	global.Functions = append(global.Functions, &Function{
		Name:          mainFunc.Name.Name,
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
				if basicLit, ok := valueSpec.Values[i].(*ast.BasicLit); ok {
					value = basicLit.Value
				} else if ident, ok := valueSpec.Values[i].(*ast.Ident); ok {
					value = ident.Name
				} else {
					return "", errors.New(fmt.Sprintf("unknown constant type: %T", valueSpec.Values[i]))
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
		if err := statement.PostProcess(global, nil); err != nil {
			return "", err
		}
	}

	for _, fn := range global.Functions {
		for _, statement := range fn.Statements {
			if err := statement.PostProcess(global, fn); err != nil {
				return "", err
			}
		}
	}

	result := ""
	lineNumber := 0
	for _, statement := range startup {
		statements := statement.ToMLOG()
		result += MLOGToString(statements, statement, lineNumber, options)
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
					result += MLOGToString(deb, debugStatement, lineNumber, options)
					lineNumber += len(deb)
				}
			}

			statements := statement.ToMLOG()
			result += MLOGToString(statements, statement, lineNumber, options)
			lineNumber += len(statements)
		}
	}

	return result, nil
}
