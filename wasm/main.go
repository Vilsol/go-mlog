package main

import (
	"fmt"
	_ "github.com/Vilsol/go-mlog/m"
	"github.com/Vilsol/go-mlog/transpiler"
	_ "github.com/Vilsol/go-mlog/x"
	"runtime/debug"
	"syscall/js"
)

func transpileWrapper() js.Func {
	transpileFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		input := args[0].String()

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic: %s\n", r)
				fmt.Println(string(debug.Stack()))
			}
		}()

		mlog, err := transpiler.GolangToMLOG(input, transpiler.Options{
			Numbers:  false,
			Comments: false,
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
	<-make(chan bool)
}
