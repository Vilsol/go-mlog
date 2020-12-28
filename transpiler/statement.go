package transpiler

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
)

func statementToMLOG(ctx context.Context, statement ast.Stmt) ([]MLOGStatement, error) {
	subCtx := context.WithValue(ctx, contextStatement, statement)

	results := make([]MLOGStatement, 0)

	switch castStmt := statement.(type) {
	case *ast.ForStmt:
		// TODO Switch from do while to while do

		if len(castStmt.Body.List) == 0 {
			break
		}

		initMlog, err := statementToMLOG(subCtx, castStmt.Init)
		if err != nil {
			return nil, err
		}
		results = append(results, initMlog...)

		var loopStartJump *MLOGJump
		var loopEndJump *MLOGJump
		if binaryExpr, ok := castStmt.Cond.(*ast.BinaryExpr); ok {
			if translatedOp, ok := jumpOperators[binaryExpr.Op]; ok {
				var leftSide Resolvable
				var rightSide Resolvable

				// TODO Convert to switch
				if basicLit, ok := binaryExpr.X.(*ast.BasicLit); ok {
					leftSide = &Value{Value: basicLit.Value}
				} else if ident, ok := binaryExpr.X.(*ast.Ident); ok {
					leftSide = &NormalVariable{Name: ident.Name}
				} else {
					return nil, Err(subCtx, fmt.Sprintf("unknown left side expression type: %T", binaryExpr.X))
				}

				// TODO Convert to switch
				if basicLit, ok := binaryExpr.Y.(*ast.BasicLit); ok {
					rightSide = &Value{Value: basicLit.Value}
				} else if ident, ok := binaryExpr.Y.(*ast.Ident); ok {
					rightSide = &NormalVariable{Name: ident.Name}
				} else {
					return nil, Err(subCtx, fmt.Sprintf("unknown right side expression type: %T", binaryExpr.Y))
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
				return nil, Err(subCtx, fmt.Sprintf("jump statement cannot use this operation: %T", binaryExpr.Op))
			}
		} else {
			return nil, Err(subCtx, "for loop can only have binary conditional expressions")
		}

		blockCtxStruct := &ContextBlock{}
		bodyMLOG, err := statementToMLOG(context.WithValue(subCtx, contextBreakableBlock, blockCtxStruct), castStmt.Body)
		if err != nil {
			return nil, err
		}
		blockCtxStruct.Statements = bodyMLOG

		results = append(results, loopEndJump)

		results = append(results, bodyMLOG...)

		instructions, err := statementToMLOG(subCtx, castStmt.Post)
		if err != nil {
			return nil, err
		}
		results = append(results, instructions...)
		blockCtxStruct.Extra = append(blockCtxStruct.Extra, instructions...)

		loopStartJump.JumpTarget = bodyMLOG[0]
		results = append(results, loopStartJump)
		blockCtxStruct.Extra = append(blockCtxStruct.Extra, loopStartJump)

		break
	case *ast.ExprStmt:
		instructions, err := expressionToMLOG(subCtx, nil, castStmt.X)
		if err != nil {
			return nil, err
		}

		results = append(results, instructions...)
		break
	case *ast.IfStmt:
		if castStmt.Init != nil {
			instructions, err := statementToMLOG(subCtx, castStmt.Init)
			if err != nil {
				return nil, err
			}
			results = append(results, instructions...)
		}

		var condVar Resolvable
		if condIdent, ok := castStmt.Cond.(*ast.Ident); ok {
			condVar = &NormalVariable{Name: condIdent.Name}
		} else {
			condVar = &DynamicVariable{}

			instructions, err := expressionToMLOG(subCtx, []Resolvable{condVar}, castStmt.Cond)
			if err != nil {
				return nil, err
			}

			results = append(results, instructions...)
		}

		blockInstructions, err := statementToMLOG(subCtx, castStmt.Body)
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

		if castStmt.Else != nil {
			elseInstructions, err := statementToMLOG(subCtx, castStmt.Else)
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
		assignMlog, err := assignStmtToMLOG(subCtx, castStmt)
		if err != nil {
			return nil, err
		}
		results = append(results, assignMlog...)
		break
	case *ast.ReturnStmt:
		if len(castStmt.Results) > 1 {
			// TODO Multi-value returns
			return nil, Err(subCtx, "only single value returns are supported")
		}

		if len(castStmt.Results) > 0 {
			returnValue := castStmt.Results[0]

			var resultVar Resolvable
			// TODO Convert to switch
			if ident, ok := returnValue.(*ast.Ident); ok {
				resultVar = &NormalVariable{Name: ident.Name}
			} else if basicLit, ok := returnValue.(*ast.BasicLit); ok {
				resultVar = &Value{Value: basicLit.Value}
			} else if expr, ok := returnValue.(ast.Expr); ok {
				dVar := &DynamicVariable{}

				instructions, err := expressionToMLOG(subCtx, []Resolvable{dVar}, expr)
				if err != nil {
					return nil, err
				}

				results = append(results, instructions...)
				resultVar = dVar
			} else {
				return nil, Err(subCtx, fmt.Sprintf("unknown return value type: %T", returnValue))
			}

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

		results = append(results, &MLOGTrampolineBack{})
		break
	case *ast.BlockStmt:
		blockCtxStruct := &ContextBlock{}
		statements := make([]MLOGStatement, 0)
		for _, s := range castStmt.List {
			instructions, err := statementToMLOG(context.WithValue(subCtx, contextBlock, blockCtxStruct), s)
			if err != nil {
				return nil, err
			}
			statements = append(statements, instructions...)
		}
		blockCtxStruct.Statements = statements
		results = append(results, statements...)
		break
	case *ast.IncDecStmt:
		name := &NormalVariable{Name: castStmt.X.(*ast.Ident).Name}
		op := "add"
		if castStmt.Tok == token.DEC {
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
		break
	case *ast.BranchStmt:
		switch castStmt.Tok {
		case token.BREAK:
			block := ctx.Value(contextBreakableBlock)
			if block == nil {
				return nil, Err(subCtx, fmt.Sprintf("branch statement outside any breakable block scope"))
			}
			results = append(results, &MLOGBreak{
				Block: block.(*ContextBlock),
			})
			break
		case token.CONTINUE:
			block := ctx.Value(contextBreakableBlock)
			if block == nil {
				return nil, Err(subCtx, fmt.Sprintf("branch statement outside any breakable block scope"))
			}
			results = append(results, &MLOGContinue{
				Block: block.(*ContextBlock),
			})
			break
		case token.FALLTHROUGH:
			block := ctx.Value(contextSwitchClauseBlock)
			if block == nil {
				return nil, Err(subCtx, fmt.Sprintf("fallthrough statement outside switch scope"))
			}
			results = append(results, &MLOGFallthrough{
				Block: block.(*ContextBlock),
			})
			break
		default:
			return nil, Err(subCtx, fmt.Sprintf("branch statement not supported: %s", castStmt.Tok))
		}
		break
	case *ast.SwitchStmt:
		if castStmt.Init != nil {
			instructions, err := statementToMLOG(subCtx, castStmt.Init)
			if err != nil {
				return nil, err
			}
			results = append(results, instructions...)
		}

		// TODO Convert to switch
		var tag Resolvable
		if tagBasic, ok := castStmt.Tag.(*ast.BasicLit); ok {
			tag = &Value{Value: tagBasic.Value}
		} else if tagIdent, ok := castStmt.Tag.(*ast.Ident); ok {
			tag = &NormalVariable{Name: tagIdent.Name}
		} else {
			return nil, Err(subCtx, fmt.Sprintf("unknown switch condition type: %T", castStmt.Tag))
		}

		blockCtxStruct := &ContextBlock{}
		blockCtx := context.WithValue(subCtx, contextBreakableBlock, blockCtxStruct)
		instructions := make([]MLOGStatement, 0)

		var previousSwitchClause *ContextBlock
		for _, switchStmt := range castStmt.Body.List {
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
						return nil, Err(subCtx, fmt.Sprintf("unknown switch case condition type: %T", caseExpr))
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
					endBreak := &MLOGBreak{
						Block: blockCtxStruct,
					}
					instructions = append(instructions, endBreak)

					if skipClause != nil {
						skipClause.JumpTarget.(*StatementJumpTarget).Statement = endBreak
					}
				}

				previousSwitchClause = switchClauseBlockCtxStruct
			} else {
				return nil, Err(subCtx, "switch statement may only contain case and default statements")
			}
		}

		blockCtxStruct.Statements = instructions

		results = append(results, instructions...)

		break
	default:
		return nil, Err(subCtx, fmt.Sprintf("statement type not supported: %T", statement))
	}

	return results, nil
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
				if statement.Tok != token.ASSIGN && statement.Tok != token.DEFINE {
					return nil, Err(ctx, "only direct assignment is supported")
				}

				exprMLOG, err := expressionToMLOG(ctx, []Resolvable{&NormalVariable{Name: ident.Name}}, statement.Rhs[i])
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
