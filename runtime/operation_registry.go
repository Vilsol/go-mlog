package runtime

var operationRegistry = make(map[string]OperationSetup)

func RegisterOperation(name string, setup OperationSetup) {
	if _, ok := operationRegistry[name]; ok {
		panic("operation with name " + name + "already registered")
	}

	operationRegistry[name] = setup
}
