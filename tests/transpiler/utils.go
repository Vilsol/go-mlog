package transpiler

import (
	"fmt"
	"github.com/MarvinJWendt/testza"
	_ "github.com/Vilsol/go-mlog/m"
	_ "github.com/Vilsol/go-mlog/m/impl"
	"github.com/Vilsol/go-mlog/transpiler"
	_ "github.com/Vilsol/go-mlog/x"
	_ "github.com/Vilsol/go-mlog/x/impl"
	"strings"
	"testing"
)

type Test struct {
	name   string
	input  string
	output string
}

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

func RunTests(t *testing.T, tests []Test) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mlog, err := transpiler.GolangToMLOG(test.input, transpiler.Options{
				NoStartup: true,
			})

			if err != nil {
				t.Error(err)
				return
			}

			test.output = test.output + "\nend"
			testza.AssertEqual(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
