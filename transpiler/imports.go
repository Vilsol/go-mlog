package transpiler

var validImports = make(map[string]bool)

func RegisterValidImport(path string) {
	if _, ok := funcTranslations[path]; ok {
		panic("Import path already exists: " + path)
	}

	validImports[path] = true
}
