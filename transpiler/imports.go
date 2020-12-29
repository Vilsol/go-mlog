package transpiler

var validImports = map[string]bool{
	`"github.com/Vilsol/go-mlog/m"`: true,
	`"github.com/Vilsol/go-mlog/x"`: true,
}

func RegisterValidImport(path string) {
	if _, ok := funcTranslations[path]; ok {
		panic("Import path already exists: " + path)
	}

	validImports[path] = true
}
