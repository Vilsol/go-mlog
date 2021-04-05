package transpiler

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

func statementToMLOG(ctx context.Context, statement ast.Stmt) ([]MLOGStatement, error) {
	subCtx := context.WithValue(ctx, contextStatement, statement)

	switch castStmt := statement.(type) {
	case *ast.ForStmt:
		return forStmtToMLOG(subCtx, castStmt)
	case *ast.ExprStmt:
		return expressionToMLOG(subCtx, nil, castStmt.X)
	case *ast.IfStmt:
		return ifStmtToMLOG(subCtx, castStmt)
	case *ast.AssignStmt:
		return assignStmtToMLOG(subCtx, castStmt)
	case *ast.ReturnStmt:
		return returnStmtToMLOG(subCtx, castStmt)
	case *ast.BlockStmt:
		return blockStmtToMLOG(subCtx, castStmt)
	case *ast.IncDecStmt:
		return incDecStmtToMLOG(subCtx, castStmt)
	case *ast.BranchStmt:
		return branchStmtToMLOG(subCtx, castStmt)
	case *ast.SwitchStmt:
		return switchStmtToMLOG(subCtx, castStmt)
	}

	return nil, Err(subCtx, fmt.Sprintf("statement type not supported: %T", statement))
}

func assignStmtToMLOG(ctx context.Context, statement *ast.AssignStmt) ([]MLOGStatement, error) {
	mlog := make([]MLOGStatement, 0)

	if len(statement.Lhs) != len(statement.Rhs) {
		if len(statement.Rhs) == 1 {
			leftSide := make([]Resolvable, len(statement.Lhs))

			for i, lhs := range statement.Lhs {
				leftSide[i] = &NormalVariable{Name: lhs.(*ast.Ident).Name}
			}

			if callExpr, ok := statement.Rhs[0].(*ast.CallExpr); ok {
				count, err := getFunctionReturnCount(ctx, callExpr)
				if err != nil {
					return nil, err
				}

				if count != len(statement.Lhs) {
					return nil, Err(ctx, "mismatched variable assignment sides")
				}
			}

			exprMLOG, err := expressionToMLOG(ctx, leftSide, statement.Rhs[0])
			if err != nil {
				return nil, err
			}
			mlog = append(mlog, exprMLOG...)
		} else {
			return nil, Err(ctx, "mismatched variable assignment sides")
		}
	} else {
		for i, expr := range statement.Lhs {
			if ident, ok := expr.(*ast.Ident); ok {
				nVar := &NormalVariable{Name: ident.Name}
				if opTranslated, ok := regularOperators[statement.Tok]; ok {
					instructions := make([]MLOGStatement, 0)

					rightSide, rightExprInstructions, err := exprToResolvable(ctx, statement.Rhs[i])
					if err != nil {
						return nil, err
					}
					instructions = append(instructions, rightExprInstructions...)

					if len(rightSide) != 1 {
						return nil, Err(ctx, "unknown error")
					}

					return append(instructions, &MLOG{
						Comment: "Execute operation",
						Statement: [][]Resolvable{
							{
								&Value{Value: "op"},
								&Value{Value: opTranslated},
								nVar,
								nVar,
								rightSide[0],
							},
						},
						SourcePos: statement,
					}), nil
				}

				if statement.Tok != token.ASSIGN && statement.Tok != token.DEFINE {
					return nil, Err(ctx, "only direct assignment is supported")
				}

				if callExpr, ok := statement.Rhs[i].(*ast.CallExpr); ok {
					count, err := getFunctionReturnCount(ctx, callExpr)
					if err != nil {
						return nil, err
					}

					if count != len(statement.Lhs) {
						return nil, Err(ctx, "mismatched variable assignment sides")
					}
				}

				exprMLOG, err := expressionToMLOG(ctx, []Resolvable{nVar}, statement.Rhs[i])
				if err != nil {
					return nil, err
				}
				mlog = append(mlog, exprMLOG...)
			} else {
				return nil, Err(ctx, "left side variable assignment can only contain identifications")
			}
		}
	}

	return mlog, nil
}

func returnStmtToMLOG(ctx context.Context, statement *ast.ReturnStmt) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)

	if len(statement.Results) > 0 {
		for i, returnValue := range statement.Results {
			resultVar, exprInstructions, err := exprToResolvable(ctx, returnValue)
			if err != nil {
				return nil, err
			}

			if len(resultVar) != 1 {
				return nil, Err(ctx, "unknown error")
			}

			results = append(results, exprInstructions...)

			results = append(results, &MLOG{
				Comment: "Set return data",
				Statement: [][]Resolvable{
					{
						&Value{Value: "set"},
						&Value{Value: FunctionReturnVariable + "_" + strconv.Itoa(i)},
						resultVar[0],
					},
				},
				SourcePos: statement,
			})
		}
	}

	return append(results, &MLOGTrampolineBack{
		Stacked:  ctx.Value(contextOptions).(Options).Stacked,
		Function: ctx.Value(contextFunction).(*ast.FuncDecl).Name.Name,
	}), nil
}

func ifStmtToMLOG(ctx context.Context, statement *ast.IfStmt) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)

	if statement.Init != nil {
		instructions, err := statementToMLOG(ctx, statement.Init)
		if err != nil {
			return nil, err
		}
		results = append(results, instructions...)
	}

	var condVar Resolvable
	if condIdent, ok := statement.Cond.(*ast.Ident); ok {
		condVar = &NormalVariable{Name: condIdent.Name}
	} else {
		condVar = &DynamicVariable{}

		instructions, err := expressionToMLOG(ctx, []Resolvable{condVar}, statement.Cond)
		if err != nil {
			return nil, err
		}

		results = append(results, instructions...)
	}

	blockInstructions, err := statementToMLOG(ctx, statement.Body)
	if err != nil {
		return nil, err
	}

	results = append(results, &MLOGJump{
		MLOG: MLOG{
			Comment: "Jump to if block if true",
		},
		Condition: []Resolvable{
			&Value{Value: "equal"},
			condVar,
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

	if statement.Else != nil {
		elseInstructions, err := statementToMLOG(ctx, statement.Else)
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

	return results, nil
}

func forStmtToMLOG(ctx context.Context, statement *ast.ForStmt) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)

	if len(statement.Body.List) == 0 {
		return results, nil
	}

	initMlog, err := statementToMLOG(ctx, statement.Init)
	if err != nil {
		return nil, err
	}
	results = append(results, initMlog...)

	var loopStartJump *MLOGJump
	var intoLoopJump *MLOGJump
	var loopEndJump *MLOGJump
	if binaryExpr, ok := statement.Cond.(*ast.BinaryExpr); ok {
		if translatedOp, ok := jumpOperators[binaryExpr.Op]; ok {

			leftSide, leftExprInstructions, err := exprToResolvable(ctx, binaryExpr.X)
			if err != nil {
				return nil, err
			}
			results = append(results, leftExprInstructions...)

			if len(leftSide) != 1 {
				return nil, Err(ctx, "unknown error")
			}

			rightSide, rightExprInstructions, err := exprToResolvable(ctx, binaryExpr.Y)
			if err != nil {
				return nil, err
			}
			results = append(results, rightExprInstructions...)

			if len(rightSide) != 1 {
				return nil, Err(ctx, "unknown error")
			}

			loopStartJump = &MLOGJump{
				MLOG: MLOG{
					Comment: "Jump to start of loop",
				},
				Condition: []Resolvable{
					&Value{Value: translatedOp},
					leftSide[0],
					rightSide[0],
				},
			}

			intoLoopJump = &MLOGJump{
				MLOG: MLOG{
					Comment: "Jump into the loop",
				},
				Condition: []Resolvable{
					&Value{Value: translatedOp},
					leftSide[0],
					rightSide[0],
				},
			}

			loopEndJump = &MLOGJump{
				MLOG: MLOG{
					Comment: "Jump to end of loop",
				},
				Condition: []Resolvable{
					&Value{Value: "always"},
				},
				JumpTarget: &StatementJumpTarget{
					Statement: loopStartJump,
					After:     true,
				},
			}
		} else {
			return nil, Err(ctx, fmt.Sprintf("jump statement cannot use this operation: %T", binaryExpr.Op))
		}
	} else {
		return nil, Err(ctx, "for loop can only have binary conditional expressions")
	}

	blockCtxStruct := &ContextBlock{}
	bodyMLOG, err := statementToMLOG(context.WithValue(ctx, contextBreakableBlock, blockCtxStruct), statement.Body)
	if err != nil {
		return nil, err
	}
	blockCtxStruct.Statements = bodyMLOG

	intoLoopJump.JumpTarget = bodyMLOG[0]
	results = append(results, intoLoopJump)

	results = append(results, loopEndJump)

	results = append(results, bodyMLOG...)

	instructions, err := statementToMLOG(ctx, statement.Post)
	if err != nil {
		return nil, err
	}
	results = append(results, instructions...)
	blockCtxStruct.Extra = append(blockCtxStruct.Extra, instructions...)

	loopStartJump.JumpTarget = bodyMLOG[0]
	results = append(results, loopStartJump)
	blockCtxStruct.Extra = append(blockCtxStruct.Extra, loopStartJump)

	return results, nil
}

func blockStmtToMLOG(ctx context.Context, statement *ast.BlockStmt) ([]MLOGStatement, error) {
	blockCtxStruct := &ContextBlock{}
	statements := make([]MLOGStatement, 0)
	for _, s := range statement.List {
		instructions, err := statementToMLOG(context.WithValue(ctx, contextBlock, blockCtxStruct), s)
		if err != nil {
			return nil, err
		}
		statements = append(statements, instructions...)
	}
	blockCtxStruct.Statements = statements

	return statements, nil
}

func incDecStmtToMLOG(_ context.Context, statement *ast.IncDecStmt) ([]MLOGStatement, error) {
	name := &NormalVariable{Name: statement.X.(*ast.Ident).Name}
	op := "add"
	if statement.Tok == token.DEC {
		op = "sub"
	}
	return []MLOGStatement{&MLOG{
		Comment: "Execute increment/decrement",
		Statement: [][]Resolvable{
			{
				&Value{Value: "op"},
				&Value{Value: op},
				name,
				name,
				&Value{Value: "1"},
			},
		},
		SourcePos: statement,
	}}, nil
}

func branchStmtToMLOG(ctx context.Context, statement *ast.BranchStmt) ([]MLOGStatement, error) {
	switch statement.Tok {
	case token.BREAK:
		fallthrough
	case token.CONTINUE:
		block := ctx.Value(contextBreakableBlock)
		if block == nil {
			return nil, Err(ctx, fmt.Sprintf("branch statement outside any breakable block scope"))
		}
		return []MLOGStatement{&MLOGBranch{
			Block: block.(*ContextBlock),
			Token: statement.Tok,
		}}, nil
	case token.FALLTHROUGH:
		// Requires no extra instructions
		return []MLOGStatement{}, nil
	}

	return nil, Err(ctx, fmt.Sprintf("branch statement not supported: %s", statement.Tok))
}

func switchStmtToMLOG(ctx context.Context, statement *ast.SwitchStmt) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)

	if statement.Init != nil {
		instructions, err := statementToMLOG(ctx, statement.Init)
		if err != nil {
			return nil, err
		}
		results = append(results, instructions...)
	}

	tag, leftExprInstructions, err := exprToResolvable(ctx, statement.Tag)
	if err != nil {
		return nil, err
	}
	results = append(results, leftExprInstructions...)

	if len(tag) != 1 {
		return nil, Err(ctx, "unknown error")
	}

	blockCtxStruct := &ContextBlock{}
	blockCtx := context.WithValue(ctx, contextBreakableBlock, blockCtxStruct)

	jumpInstructions := make([]MLOGStatement, 0)
	instructions := make([]MLOGStatement, 0)

	var previousSwitchClause *ContextBlock
	for _, switchStmt := range statement.Body.List {
		if caseStmt, ok := switchStmt.(*ast.CaseClause); ok {
			statements := make([]MLOGStatement, 0)
			switchClauseBlockCtxStruct := &ContextBlock{}
			for _, s := range caseStmt.Body {
				bodyInstructions, err := statementToMLOG(context.WithValue(blockCtx, contextSwitchClauseBlock, switchClauseBlockCtxStruct), s)
				if err != nil {
					return nil, err
				}
				statements = append(statements, bodyInstructions...)
			}
			switchClauseBlockCtxStruct.Statements = statements

			for _, caseExpr := range caseStmt.List {
				var caseTag Resolvable
				if tagBasic, ok := caseExpr.(*ast.BasicLit); ok {
					caseTag = &Value{Value: tagBasic.Value}
				} else if tagIdent, ok := caseExpr.(*ast.Ident); ok {
					caseTag = &NormalVariable{Name: tagIdent.Name}
				} else {
					return nil, Err(ctx, fmt.Sprintf("unknown switch case condition type: %T", caseExpr))
				}

				jumpIn := &MLOGJump{
					MLOG: MLOG{
						Comment: "Jump in if match",
					},
					Condition: []Resolvable{
						&Value{Value: "equal"},
						tag[0],
						caseTag,
					},
					JumpTarget: &StatementJumpTarget{
						Statement: statements[0],
					},
				}
				jumpInstructions = append(jumpInstructions, jumpIn)
				if previousSwitchClause != nil {
					previousSwitchClause.Extra = append(previousSwitchClause.Extra, jumpIn)
				}
			}

			instructions = append(instructions, statements...)

			addJump := len(caseStmt.Body) == 0
			if len(caseStmt.Body) >= 0 {
				if _, ok := caseStmt.Body[len(caseStmt.Body)-1].(*ast.BranchStmt); !ok {
					addJump = true
				}
			}

			if addJump {
				endBreak := &MLOGBranch{
					Block: blockCtxStruct,
				}
				instructions = append(instructions, endBreak)
			}

			previousSwitchClause = switchClauseBlockCtxStruct
		} else {
			return nil, Err(ctx, "switch statement may only contain case and default statements")
		}
	}

	combined := append(
		jumpInstructions,
		append(
			[]MLOGStatement{&MLOGBranch{
				Block: blockCtxStruct,
			}},
			instructions...,
		)...,
	)

	blockCtxStruct.Statements = combined

	return append(results, combined...), nil
}
