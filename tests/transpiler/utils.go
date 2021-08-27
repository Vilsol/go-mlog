package transpiler

import (
	"fmt"
	_ "github.com/Vilsol/go-mlog/m"
	_ "github.com/Vilsol/go-mlog/m/impl"
	_ "github.com/Vilsol/go-mlog/x"
	_ "github.com/Vilsol/go-mlog/x/impl"
)

func TestMain(main string, useM bool, useX bool) string {
	result := "package main\n\n"

	if useM || useX {
		result += "import (\n"
		if useM {
			result += "\"github.com/Vilsol/go-mlog/m\""
		}
		if useX {
			result += "\"github.com/Vilsol/go-mlog/x\""
		}
		result += ")"
	}

	return fmt.Sprintf(`%s

func main() {
%s
}`, result, main)
}
