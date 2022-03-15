package transpiler

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

func statementToMLOG(ctx context.Context, statement ast.Stmt) ([]MLOGStatement, []*VarReference, error) {
	subCtx := context.WithValue(ctx, contextStatement, statement)

	switch castStmt := statement.(type) {
	case *ast.ForStmt:
		mlog, err := forStmtToMLOG(subCtx, castStmt)
		return mlog, nil, err
	case *ast.ExprStmt:
		mlog, err := expressionToMLOG(subCtx, nil, castStmt.X)
		return mlog, nil, err
	case *ast.IfStmt:
		return ifStmtToMLOG(subCtx, castStmt)
	case *ast.AssignStmt:
		return assignStmtToMLOG(subCtx, castStmt)
	case *ast.ReturnStmt:
		mlog, err := returnStmtToMLOG(subCtx, castStmt)
		return mlog, nil, err
	case *ast.BlockStmt:
		return blockStmtToMLOG(subCtx, castStmt)
	case *ast.IncDecStmt:
		mlog, err := incDecStmtToMLOG(subCtx, castStmt)
		return mlog, nil, err
	case *ast.BranchStmt:
		mlog, err := branchStmtToMLOG(subCtx, castStmt)
		return mlog, nil, err
	case *ast.SwitchStmt:
		mlog, err := switchStmtToMLOG(subCtx, castStmt)
		return mlog, nil, err
	case *ast.LabeledStmt:
		mlog, err := labeledStmtToMLOG(subCtx, castStmt)
		return mlog, nil, err
	}

	return nil, nil, Err(subCtx, fmt.Sprintf("statement type not supported: %T", statement))
}

func assignStmtToMLOG(ctx context.Context, statement *ast.AssignStmt) ([]MLOGStatement, []*VarReference, error) {
	mlog := make([]MLOGStatement, 0)
	varReferences := make([]*VarReference, 0)

	if len(statement.Lhs) != len(statement.Rhs) {
		if len(statement.Rhs) == 1 {
			leftSide := make([]Resolvable, len(statement.Lhs))

			for i, lhs := range statement.Lhs {
				var nVar Resolvable

				if statement.Tok != token.DEFINE {
					nVar = contextOrVariable(ctx, lhs.(*ast.Ident).Name)
				}

				if nVar == nil {
					nVar = &NormalVariable{Name: lhs.(*ast.Ident).Name}
				}

				leftSide[i] = nVar

				if statement.Tok == token.DEFINE {
					varReferences = append(varReferences, &VarReference{
						Name:     lhs.(*ast.Ident).Name,
						Identity: nVar,
					})
				}
			}

			if callExpr, ok := statement.Rhs[0].(*ast.CallExpr); ok {
				count, err := getFunctionReturnCount(ctx, callExpr)
				if err != nil {
					return nil, nil, err
				}

				if count != len(statement.Lhs) {
					return nil, nil, Err(ctx, "mismatched variable assignment sides")
				}
			}

			exprMLOG, err := expressionToMLOG(ctx, leftSide, statement.Rhs[0])
			if err != nil {
				return nil, nil, err
			}
			mlog = append(mlog, exprMLOG...)
		} else {
			return nil, nil, Err(ctx, "mismatched variable assignment sides")
		}
	} else {
		for i, expr := range statement.Lhs {
			if ident, ok := expr.(*ast.Ident); ok {
				var nVar Resolvable

				if statement.Tok != token.DEFINE {
					nVar = contextOrVariable(ctx, ident.Name)
				}

				if nVar == nil {
					nVar = &NormalVariable{Name: ident.Name}
				}

				if opTranslated, ok := regularOperators[statement.Tok]; ok {
					instructions := make([]MLOGStatement, 0)

					rightSide, rightExprInstructions, err := exprToResolvable(ctx, statement.Rhs[i])
					if err != nil {
						return nil, nil, err
					}
					instructions = append(instructions, rightExprInstructions...)

					if len(rightSide) != 1 {
						return nil, nil, Err(ctx, "unknown error")
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
					}), varReferences, nil
				}

				if statement.Tok != token.ASSIGN && statement.Tok != token.DEFINE {
					return nil, nil, Err(ctx, "only direct assignment is supported")
				}

				if callExpr, ok := statement.Rhs[i].(*ast.CallExpr); ok {
					count, err := getFunctionReturnCount(ctx, callExpr)
					if err != nil {
						return nil, nil, err
					}

					if count != len(statement.Lhs) {
						return nil, nil, Err(ctx, "mismatched variable assignment sides")
					}
				}

				exprMLOG, err := expressionToMLOG(ctx, []Resolvable{nVar}, statement.Rhs[i])
				if err != nil {
					return nil, nil, err
				}
				mlog = append(mlog, exprMLOG...)

				if statement.Tok == token.DEFINE {
					varReferences = append(varReferences, &VarReference{
						Name:     ident.Name,
						Identity: nVar,
					})
				}
			} else {
				return nil, nil, Err(ctx, "left side variable assignment can only contain identifications")
			}
		}
	}

	return mlog, varReferences, nil
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

func ifStmtToMLOG(ctx context.Context, statement *ast.IfStmt) ([]MLOGStatement, []*VarReference, error) {
	results := make([]MLOGStatement, 0)
	varReferences := make([]*VarReference, 0)

	if statement.Init != nil {
		instructions, references, err := statementToMLOG(ctx, statement.Init)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, instructions...)
		varReferences = append(varReferences, references...)
		ctx = addVariablesToContext(ctx, references)
	}

	var condVar Resolvable
	if condIdent, ok := statement.Cond.(*ast.Ident); ok {
		condVar = contextOrVariable(ctx, condIdent.Name)
	} else {
		condVar = &DynamicVariable{}

		instructions, err := expressionToMLOG(ctx, []Resolvable{condVar}, statement.Cond)
		if err != nil {
			return nil, nil, err
		}

		results = append(results, instructions...)
	}

	blockInstructions, references, err := statementToMLOG(ctx, statement.Body)
	if err != nil {
		return nil, nil, err
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

	varReferences = append(varReferences, references...)
	ctx = addVariablesToContext(ctx, references)

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
		elseInstructions, _, err := statementToMLOG(ctx, statement.Else)
		if err != nil {
			return nil, nil, err
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

	return results, varReferences, nil
}

func forStmtToMLOG(ctx context.Context, statement *ast.ForStmt) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)
	varReferences := make([]*VarReference, 0)

	if len(statement.Body.List) == 0 {
		return results, nil
	}

	if statement.Init != nil {
		initMlog, references, err := statementToMLOG(ctx, statement.Init)
		if err != nil {
			return nil, err
		}
		results = append(results, initMlog...)
		varReferences = append(varReferences, references...)
		ctx = addVariablesToContext(ctx, references)
	}

	var loopStartJump *MLOGJump
	var intoLoopJump *MLOGJump

	var loopStartOverride *MLOGStatement

	if statement.Cond != nil {
		loopStartJump = &MLOGJump{
			MLOG: MLOG{
				Comment: "Jump to start of loop",
			},
		}

		intoLoopJump = &MLOGJump{
			MLOG: MLOG{
				Comment: "Jump into the loop",
			},
		}

		if binaryExpr, ok := statement.Cond.(*ast.BinaryExpr); ok {
			// TODO Optimize jump instruction if possible

			expr, exprInstructions, err := exprToResolvable(ctx, binaryExpr)
			if err != nil {
				return nil, err
			}

			results = append(results, exprInstructions...)

			if len(expr) != 1 {
				return nil, Err(ctx, "unknown error")
			}

			loopStartJump.Condition = []Resolvable{
				&Value{Value: "always"},
			}

			intoLoopJump.Condition = []Resolvable{
				&Value{Value: jumpOperators[token.EQL]},
				expr[0],
				&Value{Value: "true"},
			}

			loopStartOverride = &exprInstructions[0]
		} else if unaryExpr, ok := statement.Cond.(*ast.UnaryExpr); ok {
			if unaryExpr.Op != token.NOT {
				return nil, Err(ctx, fmt.Sprintf("loop unary expresion cannot use this operation: %T", binaryExpr.Op))
			}

			expr, exprInstructions, err := exprToResolvable(ctx, unaryExpr.X)
			if err != nil {
				return nil, err
			}

			results = append(results, exprInstructions...)

			if len(expr) != 1 {
				return nil, Err(ctx, "unknown error")
			}

			loopStartJump.Condition = []Resolvable{
				&Value{Value: jumpOperators[token.NEQ]},
				expr[0],
				&Value{Value: "true"},
			}

			intoLoopJump.Condition = []Resolvable{
				&Value{Value: jumpOperators[token.NEQ]},
				expr[0],
				&Value{Value: "true"},
			}
		} else {
			return nil, Err(ctx, "for loop can only have unary or binary conditional expressions")
		}
	} else {
		loopStartJump = &MLOGJump{
			MLOG: MLOG{
				Comment: "Jump to start of loop",
			},
			Condition: []Resolvable{
				&Value{Value: "always"},
			},
		}

		intoLoopJump = &MLOGJump{
			MLOG: MLOG{
				Comment: "Jump into the loop",
			},
			Condition: []Resolvable{
				&Value{Value: "always"},
			},
		}
	}

	blockCtxStruct := &ContextBlock{}
	bodyMLOG, _, err := statementToMLOG(context.WithValue(ctx, contextBreakableBlock, blockCtxStruct), statement.Body)
	if err != nil {
		return nil, err
	}
	blockCtxStruct.Statements = bodyMLOG

	intoLoopJump.JumpTarget = bodyMLOG[0]
	results = append(results, intoLoopJump)

	if statement.Cond != nil {
		results = append(results, &MLOGJump{
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
		})
	}

	results = append(results, bodyMLOG...)

	if statement.Post != nil {
		instructions, _, err := statementToMLOG(ctx, statement.Post)
		if err != nil {
			return nil, err
		}
		results = append(results, instructions...)
		blockCtxStruct.Extra = append(blockCtxStruct.Extra, instructions...)
	}

	if loopStartOverride != nil {
		loopStartJump.JumpTarget = *loopStartOverride
	} else {
		loopStartJump.JumpTarget = bodyMLOG[0]
	}

	results = append(results, loopStartJump)
	blockCtxStruct.Extra = append(blockCtxStruct.Extra, loopStartJump)

	return results, nil
}

func blockStmtToMLOG(ctx context.Context, statement *ast.BlockStmt) ([]MLOGStatement, []*VarReference, error) {
	blockCtxStruct := &ContextBlock{}
	statements := make([]MLOGStatement, 0)
	varReferences := make([]*VarReference, 0)
	for _, s := range statement.List {
		instructions, references, err := statementToMLOG(context.WithValue(ctx, contextBlock, blockCtxStruct), s)
		if err != nil {
			return nil, nil, err
		}
		statements = append(statements, instructions...)
		varReferences = append(varReferences, references...)
		ctx = addVariablesToContext(ctx, references)
	}
	blockCtxStruct.Statements = statements

	return statements, varReferences, nil
}

func incDecStmtToMLOG(ctx context.Context, statement *ast.IncDecStmt) ([]MLOGStatement, error) {
	name := contextOrVariable(ctx, statement.X.(*ast.Ident).Name)
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
			return nil, Err(ctx, "branch statement outside any breakable block scope")
		}
		return []MLOGStatement{&MLOGBranch{
			Block: block.(*ContextBlock),
			Token: statement.Tok,
		}}, nil
	case token.FALLTHROUGH:
		// Requires no extra instructions
		return []MLOGStatement{}, nil
	case token.GOTO:
		return []MLOGStatement{
			&MLOG{
				Statement: [][]Resolvable{
					{
						&Value{Value: "jump"},
						&Value{Value: statement.Label.Name},
					},
				},
				Comment:   "Jump to label",
				SourcePos: statement,
			},
		}, nil
	}

	return nil, Err(ctx, fmt.Sprintf("branch statement not supported: %s", statement.Tok))
}

func switchStmtToMLOG(ctx context.Context, statement *ast.SwitchStmt) ([]MLOGStatement, error) {
	results := make([]MLOGStatement, 0)

	if statement.Init != nil {
		instructions, references, err := statementToMLOG(ctx, statement.Init)
		if err != nil {
			return nil, err
		}
		results = append(results, instructions...)
		ctx = addVariablesToContext(ctx, references)
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
				bodyInstructions, _, err := statementToMLOG(context.WithValue(blockCtx, contextSwitchClauseBlock, switchClauseBlockCtxStruct), s)
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
					caseTag = contextOrVariable(ctx, tagIdent.Name)
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

func labeledStmtToMLOG(ctx context.Context, statement *ast.LabeledStmt) ([]MLOGStatement, error) {
	subStmt, _, err := statementToMLOG(ctx, statement.Stmt)
	if err != nil {
		return nil, err
	}

	return append([]MLOGStatement{
		&MLOGLabel{
			MLOG: MLOG{
				SourcePos: statement,
			},
			Name: statement.Label.Name,
		},
	}, subStmt...), nil
}
