package main

import (
	"github.com/Vilsol/go-mlog/cmd"
	_ "github.com/Vilsol/go-mlog/m"
	_ "github.com/Vilsol/go-mlog/m/impl"
	_ "github.com/Vilsol/go-mlog/x"
	_ "github.com/Vilsol/go-mlog/x/impl"
)

func main() {
	cmd.Execute()
}
