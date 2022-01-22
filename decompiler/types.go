package decompiler

import (
	"context"
	"errors"
	"github.com/Vilsol/go-mlog/mlog"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type Global struct {
	Lines       []mlog.MLOGLine
	Labels      map[string]*mlog.MLOGLine
	MappedLines map[int]*mlog.MLOGLine
	Variables   map[string]string
}

func (g Global) AssignOrDefine(variable string, variableType string) (token.Token, error) {
	if storedType, ok := g.Variables[variable]; ok {
		if variableType != storedType {
			return 0, errors.New("attempting to assign type " + variableType + " to " + variable + " [" + storedType + "]")
		}

		return token.ASSIGN, nil
	}

	g.Variables[variable] = variableType

	return token.DEFINE, nil
}

func (g Global) Exists(variable string, variableType string) (bool, error) {
	if storedType, ok := g.Variables[variable]; ok {
		if variableType != storedType {
			return true, errors.New("attempting to use " + variable + " [" + storedType + "] as " + variableType)
		}

		return true, nil
	}

	return false, nil
}

func (g Global) Resolve(str string, variableType string) (ast.Expr, error) {
	if str == "@this" {
		return &ast.SelectorExpr{
			X:   ast.NewIdent("m"),
			Sel: ast.NewIdent("This"),
		}, nil
	}

	if variableType == "string" && strings.HasPrefix(str, "@") {
		return &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("m"),
				Sel: ast.NewIdent("Const"),
			},
			Args: []ast.Expr{
				ast.NewIdent("\"" + str + "\""),
			},
		}, nil
	}

	exists, err := g.Exists(str, variableType)
	if exists {
		if err != nil {
			return nil, err
		}

		return ast.NewIdent(str), nil
	}

	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("m"),
			Sel: ast.NewIdent("B"),
		},
		Args: []ast.Expr{
			ast.NewIdent("\"" + str + "\""),
		},
	}, nil
}

func (g Global) GetType(variable string) (string, error) {
	if storedType, ok := g.Variables[variable]; ok {
		return storedType, nil
	}

	return "", errors.New("variable not defined: " + variable)
}

func (g Global) Assert(variable string, variableType string) error {
	if variableType == "int" {
		if _, err := strconv.ParseInt(variable, 10, 64); err == nil {
			return nil
		}
	}

	if variableType == "float64" {
		if _, err := strconv.ParseFloat(variable, 64); err == nil {
			return nil
		}
	}

	storedType, err := g.GetType(variable)
	if err != nil {
		return err
	}

	if storedType != variableType {
		return errors.New("provided variable " + variable + " [" + storedType + "] cannot be used as " + variableType)
	}

	return nil
}

type Resolvable interface {
	PostProcess(context.Context, *Global) error
	GetValue() string
}
