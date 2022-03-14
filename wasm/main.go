//go:build wasm

package main

import (
	"encoding/json"
	"fmt"
	"github.com/Vilsol/go-mlog/checker"
	_ "github.com/Vilsol/go-mlog/m"
	_ "github.com/Vilsol/go-mlog/m/impl"
	"github.com/Vilsol/go-mlog/transpiler"
	_ "github.com/Vilsol/go-mlog/x"
	_ "github.com/Vilsol/go-mlog/x/impl"
	"runtime/debug"
	"syscall/js"
)

func transpileWrapper() js.Func {
	transpileFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 4 {
			return "Invalid no of arguments passed"
		}
		input := args[0].String()
		numbers := args[1].Bool()
		comments := args[2].Bool()
		source := args[3].Bool()

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic: %s\n", r)
				fmt.Println(string(debug.Stack()))
			}
		}()

		mlog, err := transpiler.GolangToMLOG(input, transpiler.Options{
			Numbers:  numbers,
			Comments: comments,
			Source:   source,
		})
		if err != nil {
			fmt.Printf("error transpiling: %s\n", err)
			return err.Error()
		}
		return mlog
	})
	return transpileFunc
}

func main() {
	fmt.Println("Transpiler Initialized")

	js.Global().Set("transpileGo", transpileWrapper())

	result := checker.GetSerializablePackages()
	marshal, _ := json.Marshal(result)

	var processed map[string]interface{}
	_ = json.Unmarshal(marshal, &processed)

	js.Global().Set("goTypings", processed)
	select {}
}
