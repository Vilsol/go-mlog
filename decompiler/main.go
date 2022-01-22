package decompiler

import (
	"bytes"
	"errors"
	"github.com/Vilsol/go-mlog/mlog"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"strconv"
)

const LabelPrefix = "jumpTo"

func MLOGToGolangFile(fileName string) (string, error) {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return MLOGToGolangBytes(file)
}

func MLOGToGolangBytes(input []byte) (string, error) {
	return MLOGToGolang(string(input))
}

func MLOGToGolang(input string) (string, error) {
	lines, _ := mlog.Tokenize(input)

	global := &Global{
		Lines:       lines,
		Labels:      make(map[string]*mlog.MLOGLine),
		MappedLines: make(map[int]*mlog.MLOGLine),
		Variables:   make(map[string]string),
	}

	allJumpTargets := make(map[int]bool)
	for _, line := range lines {
		// Detect "set @counter" and early exit
		// TODO Benchmark
		if len(line.Instruction) >= 2 && line.Instruction[0] == "set" && line.Instruction[1] == "@counter" {
			return "", errors.New("decompiler does not support programs that set @counter variable")
		}

		tempLine := line
		global.MappedLines[line.SourceLine] = &tempLine

		if line.Label != "" {
			global.Labels[line.Label] = &tempLine
		}

		if len(line.Instruction) > 0 {
			translator, ok := funcTranslations[line.Instruction[0]]
			if !ok {
				return "", errors.New("unknown statement: " + line.Instruction[0])
			}

			if translator.Preprocess != nil {
				jumpTargets, err := translator.Preprocess(line.Instruction[1:])
				if err != nil {
					return "", err
				}

				for _, target := range jumpTargets {
					allJumpTargets[target] = true
				}
			}
		}
	}

	allImports := make(map[string]bool)
	statements := make([]ast.Stmt, 0)
	for _, line := range lines {
		// TODO Comments
		if len(line.Instruction) > 0 {
			translator := funcTranslations[line.Instruction[0]]

			statement, imports, err := translator.Translate(line.Instruction[1:], global)
			if err != nil {
				return "", err
			}

			for _, s := range imports {
				allImports[s] = true
			}

			labelName := line.Label
			if labelName == "" {
				if _, ok := allJumpTargets[line.SourceLine]; ok {
					labelName = LabelPrefix + strconv.Itoa(line.SourceLine)
				}
			}

			if labelName != "" && len(statement) > 0 {
				statements = append(statements, &ast.LabeledStmt{
					Label: ast.NewIdent(labelName),
					Stmt:  statement[0],
				})
				statements = append(statements, statement[1:]...)
			} else {
				statements = append(statements, statement...)
			}
		}
	}

	importSpecs := make([]ast.Spec, 0)
	for s := range allImports {
		importSpecs = append(importSpecs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Value: "\"" + s + "\"",
			},
		})
	}

	mainFile := &ast.File{
		Name: ast.NewIdent("main"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: importSpecs,
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("main"),
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{
					List: statements,
				},
			},
		},
	}

	var buf bytes.Buffer
	fileSet := token.NewFileSet()
	if err := printer.Fprint(&buf, fileSet, mainFile); err != nil {
		return "", err
	}

	return buf.String(), nil
}
