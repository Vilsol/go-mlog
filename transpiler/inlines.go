package transpiler

type InlineTranslator func(args []Resolvable) (Resolvable, error)

var inlineTranslations = map[string]InlineTranslator{}

func RegisterInlineTranslation(name string, translator InlineTranslator) {
	if _, ok := funcTranslations[name]; ok {
		panic("Inline translation already exists: " + name)
	}

	inlineTranslations[name] = translator
}
