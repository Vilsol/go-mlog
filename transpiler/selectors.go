package transpiler

var selectors = map[string]string{}

func RegisterSelector(name string, selector string) {
	if _, ok := selectors[name]; ok {
		panic("Selector already exists: " + name)
	}

	selectors[name] = selector
}
