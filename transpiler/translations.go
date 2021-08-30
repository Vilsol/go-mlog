package transpiler

type TranslateFunc func(args []Resolvable, vars []Resolvable) ([]MLOGStatement, error)

type Translator struct {
	Count     func(args []Resolvable, vars []Resolvable) int
	Variables int

	Translate TranslateFunc
}

var funcTranslations = map[string]Translator{}

func RegisterFuncTranslation(name string, translator Translator) {
	if _, ok := funcTranslations[name]; ok {
		panic("Function translation already exists: " + name)
	}

	funcTranslations[name] = translator
}
