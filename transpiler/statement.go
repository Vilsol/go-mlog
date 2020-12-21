package transpiler

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

func statementToMLOG(statement ast.Stmt, options Options) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)

	switch statement.(type) {
	case *ast.ForStmt:
		forStatement := statement.(*ast.ForStmt)

		// TODO Switch from do while to while do

		if len(forStatement.Body.List) == 0 {
			break
		}

		// TODO Support all statements
		if assignStatement, ok := forStatement.Init.(*ast.AssignStmt); ok {
			assignMlog, err := assignStmtToMLOG(assignStatement, options)
			if err != nil {
				return nil, err
			}
			results = append(results, assignMlog...)
		} else {
			return nil, errors.New("for loop can only have variable assignment initiators")
		}

		var loopStartJump *MLOGJump
		if binaryExpr, ok := forStatement.Cond.(*ast.BinaryExpr); ok {
			if translatedOp, ok := jumpOperators[binaryExpr.Op]; ok {
				var leftSide Resolvable
				var rightSide Resolvable

				if basicLit, ok := binaryExpr.X.(*ast.BasicLit); ok {
					leftSide = &Value{Value: basicLit.Value}
				} else if ident, ok := binaryExpr.X.(*ast.Ident); ok {
					leftSide = &NormalVariable{Name: ident.Name}
				} else {
					return nil, errors.New(fmt.Sprintf("unknown left side expression type: %T", binaryExpr.X))
				}

				if basicLit, ok := binaryExpr.Y.(*ast.BasicLit); ok {
					rightSide = &Value{Value: basicLit.Value}
				} else if ident, ok := binaryExpr.Y.(*ast.Ident); ok {
					rightSide = &NormalVariable{Name: ident.Name}
				} else {
					return nil, errors.New(fmt.Sprintf("unknown right side expression type: %T", binaryExpr.Y))
				}

				loopStartJump = &MLOGJump{
					MLOG: MLOG{
						Comment: "Jump to start of loop",
					},
					Condition: []Resolvable{
						&Value{Value: translatedOp},
						leftSide,
						rightSide,
					},
				}
				results = append(results)
			} else {
				return nil, errors.New(fmt.Sprintf("jump statement cannot use this operation: %T", binaryExpr.Op))
			}
		} else {
			return nil, errors.New("for loop can only have binary conditional expressions")
		}

		bodyMLOG, err := statementToMLOG(forStatement.Body, options)
		if err != nil {
			return nil, err
		}

		results = append(results, bodyMLOG...)

		// TODO Support more statements
		if incDecStatement, ok := forStatement.Post.(*ast.IncDecStmt); ok {
			name := &NormalVariable{Name: incDecStatement.X.(*ast.Ident).Name}
			op := "add"
			if incDecStatement.Tok == token.DEC {
				op = "sub"
			}
			results = append(results, &MLOG{
				Comment: "Execute for loop post condition increment/decrement",
				Statement: [][]Resolvable{
					{
						&Value{Value: "op"},
						&Value{Value: op},
						name,
						name,
						&Value{Value: "1"},
					},
				},
			})
		} else {
			return nil, errors.New("for loop supports only increment or decrement post statements")
		}

		loopStartJump.JumpTarget = bodyMLOG[0]
		results = append(results, loopStartJump)

		break
	case *ast.ExprStmt:
		expressionStatement := statement.(*ast.ExprStmt)

		if callExpression, ok := expressionStatement.X.(*ast.CallExpr); ok {
			instructions, err := callExprToMLOG(callExpression, options)
			if err != nil {
				return nil, err
			}
			results = append(results, instructions...)
		} else {
			return nil, errors.New(fmt.Sprintf("unknown expression statement: %T", expressionStatement.X))
		}

		break
	case *ast.IfStmt:
		ifStmt := statement.(*ast.IfStmt)

		// TODO If statement init

		dVar := &DynamicVariable{}

		instructions, err := expressionToMLOG(dVar, ifStmt.Cond, options)
		if err != nil {
			return nil, err
		}

		results = append(results, instructions...)

		blockInstructions, err := statementToMLOG(ifStmt.Body, options)
		if err != nil {
			return nil, err
		}

		results = append(results, &MLOGJump{
			MLOG: MLOG{
				Comment: "Jump to if block if true",
			},
			Condition: []Resolvable{
				&Value{Value: "equal"},
				dVar,
				&Value{Value: "1"},
			},
			JumpTarget: &StatementJumpTarget{
				Statement: blockInstructions[0],
			},
		})

		afterIfTarget := &StatementJumpTarget{
			After:     true,
			Statement: blockInstructions[len(blockInstructions)-1],
		}
		results = append(results, &MLOGJump{
			MLOG: MLOG{
				Comment: "Jump to after if block",
			},
			Condition: []Resolvable{
				&Value{Value: "always"},
			},
			JumpTarget: afterIfTarget,
		})

		results = append(results, blockInstructions...)

		if ifStmt.Else != nil {
			elseInstructions, err := statementToMLOG(ifStmt.Else, options)
			if err != nil {
				return nil, err
			}

			afterElseJump := &MLOGJump{
				MLOG: MLOG{
					Comment: "Jump to after else block",
				},
				Condition: []Resolvable{
					&Value{Value: "always"},
				},
				JumpTarget: &StatementJumpTarget{
					After:     true,
					Statement: elseInstructions[len(elseInstructions)-1],
				},
			}
			results = append(results, afterElseJump)
			afterIfTarget.Statement = afterElseJump

			results = append(results, elseInstructions...)
		}

		break
	case *ast.AssignStmt:
		assignMlog, err := assignStmtToMLOG(statement.(*ast.AssignStmt), options)
		if err != nil {
			return nil, err
		}
		results = append(results, assignMlog...)
		break
	case *ast.ReturnStmt:
		returnStmt := statement.(*ast.ReturnStmt)

		if len(returnStmt.Results) > 1 {
			// TODO Multi-value returns
			return nil, errors.New("only single value returns are supported")
		}

		if len(returnStmt.Results) > 0 {
			returnValue := returnStmt.Results[0]

			var resultVar Resolvable
			if ident, ok := returnValue.(*ast.Ident); ok {
				resultVar = &NormalVariable{Name: ident.Name}
			} else if basicLit, ok := returnValue.(*ast.BasicLit); ok {
				resultVar = &Value{Value: basicLit.Value}
			} else if expr, ok := returnValue.(ast.Expr); ok {
				dVar := &DynamicVariable{}

				instructions, err := expressionToMLOG(dVar, expr, options)
				if err != nil {
					return nil, err
				}

				results = append(results, instructions...)
				resultVar = dVar
			} else {
				return nil, errors.New(fmt.Sprintf("unknown return value type: %T", returnValue))
			}

			results = append(results, &MLOG{
				Comment: "Set return data",
				Statement: [][]Resolvable{
					{
						&Value{Value: "set"},
						&Value{Value: functionReturnVariable},
						resultVar,
					},
				},
			})
		}

		results = append(results, &MLOGTrampolineBack{})
		break
	case *ast.BlockStmt:
		blockStmt := statement.(*ast.BlockStmt)
		for _, s := range blockStmt.List {
			instructions, err := statementToMLOG(s, options)
			if err != nil {
				return nil, err
			}
			results = append(results, instructions...)
		}
		break
	default:
		return nil, errors.New(fmt.Sprintf("statement type not supported: %T", statement))
	}

	return results, nil
}

func assignStmtToMLOG(statement *ast.AssignStmt, options Options) ([]MLOGStatement, error) {
	mlog := make([]MLOGStatement, 0)

	for i, expr := range statement.Lhs {
		if ident, ok := expr.(*ast.Ident); ok {
			if statement.Tok != token.ASSIGN && statement.Tok != token.DEFINE {
				return nil, errors.New("only direct assignment is supported")
			}

			exprMLOG, err := expressionToMLOG(&NormalVariable{Name: ident.Name}, statement.Rhs[i], options)
			if err != nil {
				return nil, err
			}
			mlog = append(mlog, exprMLOG...)
		} else {
			return nil, errors.New("left side variable assignment can only contain identifications")
		}
	}

	return mlog, nil
}
