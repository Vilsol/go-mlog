package decompiler

import "go/ast"

type PreprocessFunc func(args []string) ([]int, error)
type TranslateFunc func(args []string, global *Global) ([]ast.Stmt, []string, error)

type Translator struct {
	Preprocess PreprocessFunc
	Translate  TranslateFunc
}

var funcTranslations = map[string]Translator{}

func RegisterFuncTranslation(name string, translator Translator) {
	if _, ok := funcTranslations[name]; ok {
		panic("Function translation already exists: " + name)
	}

	funcTranslations[name] = translator
}
