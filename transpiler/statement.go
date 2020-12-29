package transpiler

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
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

					return append(instructions, &MLOG{
						Comment: "Execute operation",
						Statement: [][]Resolvable{
							{
								&Value{Value: "op"},
								&Value{Value: opTranslated},
								nVar,
								nVar,
								rightSide,
							},
						},
					}), nil
				}

				if statement.Tok != token.ASSIGN && statement.Tok != token.DEFINE {
					return nil, Err(ctx, "only direct assignment is supported")
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
	if len(statement.Results) > 1 {
		// TODO Multi-value returns
		return nil, Err(ctx, "only single value returns are supported")
	}

	results := make([]MLOGStatement, 0)

	if len(statement.Results) > 0 {
		returnValue := statement.Results[0]

		resultVar, exprInstructions, err := exprToResolvable(ctx, returnValue)
		if err != nil {
			return nil, err
		}

		results = append(results, exprInstructions...)

		results = append(results, &MLOG{
			Comment: "Set return data",
			Statement: [][]Resolvable{
				{
					&Value{Value: "set"},
					&Value{Value: FunctionReturnVariable},
					resultVar,
				},
			},
		})
	}

	return append(results, &MLOGTrampolineBack{}), nil
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
	var loopEndJump *MLOGJump
	if binaryExpr, ok := statement.Cond.(*ast.BinaryExpr); ok {
		if translatedOp, ok := jumpOperators[binaryExpr.Op]; ok {

			leftSide, leftExprInstructions, err := exprToResolvable(ctx, binaryExpr.X)
			if err != nil {
				return nil, err
			}
			results = append(results, leftExprInstructions...)

			rightSide, rightExprInstructions, err := exprToResolvable(ctx, binaryExpr.Y)
			if err != nil {
				return nil, err
			}
			results = append(results, rightExprInstructions...)

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

			loopEndJump = &MLOGJump{
				MLOG: MLOG{
					Comment: "Jump to end of loop",
				},
				Condition: []Resolvable{
					&Value{Value: translatedOp},
					leftSide,
					rightSide,
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
		block := ctx.Value(contextSwitchClauseBlock)
		if block == nil {
			return nil, Err(ctx, fmt.Sprintf("fallthrough statement outside switch scope"))
		}
		return []MLOGStatement{&MLOGBranch{
			Block: block.(*ContextBlock),
			Token: statement.Tok,
		}}, nil
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

	blockCtxStruct := &ContextBlock{}
	blockCtx := context.WithValue(ctx, contextBreakableBlock, blockCtxStruct)
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
						tag,
						caseTag,
					},
					JumpTarget: &StatementJumpTarget{
						Statement: statements[0],
					},
				}
				instructions = append(instructions, jumpIn)
				if previousSwitchClause != nil {
					previousSwitchClause.Extra = append(previousSwitchClause.Extra, jumpIn)
				}
			}

			var skipClause *MLOGJump
			if len(caseStmt.List) > 0 {
				skipClause = &MLOGJump{
					MLOG: MLOG{
						Comment: "Otherwise skip clause",
					},
					Condition: []Resolvable{
						&Value{Value: "always"},
					},
					JumpTarget: &StatementJumpTarget{
						Statement: statements[len(statements)-1],
						After:     true,
					},
				}
				instructions = append(instructions, skipClause)
				if previousSwitchClause != nil {
					previousSwitchClause.Extra = append(previousSwitchClause.Extra, skipClause)
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

				if skipClause != nil {
					skipClause.JumpTarget.(*StatementJumpTarget).Statement = endBreak
				}
			}

			previousSwitchClause = switchClauseBlockCtxStruct
		} else {
			return nil, Err(ctx, "switch statement may only contain case and default statements")
		}
	}

	blockCtxStruct.Statements = instructions

	return append(results, instructions...), nil
}
