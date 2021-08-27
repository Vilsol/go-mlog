package checker

import (
	"go/ast"
	"go/parser"
	"go/token"
)

var packages = make(map[string][]*ast.File)
var registryFset = token.NewFileSet()

func RegisterPackages(path string, files map[string]string) {
	asts := make([]*ast.File, 0)

	for name, file := range files {
		f, err := parser.ParseFile(registryFset, name, file, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		if f.Imports != nil || len(f.Imports) > 0 {
			continue
		}
		asts = append(asts, f)
	}

	packages[path] = asts
}

func FindPackage(path string) []*ast.File {
	return packages[path]
}

func GetFset() *token.FileSet {
	return registryFset
}
