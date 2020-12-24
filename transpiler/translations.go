package transpiler

type Translator struct {
	Count int

	// TODO Pass return variable as argument
	Translate func(args []Resolvable) []MLOGStatement
}

var funcTranslations = map[string]Translator{}

func RegisterFuncTranslation(name string, translator Translator) {
	if _, ok := funcTranslations[name]; ok {
		panic("Function translation already exists: " + name)
	}

	funcTranslations[name] = translator
}
