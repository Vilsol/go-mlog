package decompiler

import (
	"errors"
	"go/ast"
	"go/token"
	"strconv"
)

func init() {
	RegisterFuncTranslation("print", Translator{
		Translate: func(args []string, global *Global) ([]ast.Stmt, []string, error) {
			if len(args) != 1 {
				return nil, nil, errors.New("expecting 1 argument")
			}

			return []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: ast.NewIdent("print"),
						Args: []ast.Expr{
							ast.NewIdent(args[0]),
						},
					},
				},
			}, nil, nil
		},
	})

	RegisterFuncTranslation("jump", Translator{
		Preprocess: func(args []string) ([]int, error) {
			if len(args) != 1 {
				return nil, errors.New("expecting 1 argument")
			}

			if jumpTarget, err := strconv.ParseInt(args[0], 10, 64); err == nil {
				return []int{int(jumpTarget)}, nil
			}

			return nil, nil
		},
		Translate: func(args []string, global *Global) ([]ast.Stmt, []string, error) {
			if len(args) != 1 {
				return nil, nil, errors.New("expecting 1 argument")
			}

			line, ok := global.Labels[args[0]]
			var labelName string
			if ok {
				labelName = line.Label
			} else {
				if jumpTarget, err := strconv.ParseInt(args[0], 10, 64); err == nil {
					jumpTargetInt := int(jumpTarget)
					targetLine, ok := global.MappedLines[jumpTargetInt]
					if ok && targetLine.Label != "" {
						labelName = targetLine.Label
					} else {
						labelName = LabelPrefix + strconv.Itoa(jumpTargetInt)
					}
				} else {
					return nil, nil, errors.New("unknown jump target: " + args[0])
				}
			}

			return []ast.Stmt{
				&ast.BranchStmt{
					Tok:   token.GOTO,
					Label: ast.NewIdent(labelName),
				},
			}, nil, nil
		},
	})

	RegisterFuncTranslation("set", Translator{
		Translate: func(args []string, global *Global) ([]ast.Stmt, []string, error) {
			if len(args) != 2 {
				return nil, nil, errors.New("expecting 2 arguments")
			}

			expectedType := "string"
			if _, err := strconv.ParseInt(args[1], 10, 64); err == nil {
				expectedType = "int"
			} else if _, err := strconv.ParseFloat(args[1], 64); err == nil {
				expectedType = "float64"
			} else if storedType, ok := global.Variables[args[1]]; ok {
				expectedType = storedType
			}

			tok, err := global.AssignOrDefine(args[0], expectedType)
			if err != nil {
				return nil, nil, err
			}

			return []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(args[0]),
					},
					Tok: tok,
					Rhs: []ast.Expr{
						ast.NewIdent(args[1]),
					},
				},
			}, []string{"github.com/Vilsol/go-mlog/m"}, nil
		},
	})
}
