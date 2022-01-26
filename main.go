package main

import (
	"github.com/Vilsol/go-mlog/cmd"
	_ "github.com/Vilsol/go-mlog/m"
	_ "github.com/Vilsol/go-mlog/m/impl"
	_ "github.com/Vilsol/go-mlog/x"
	_ "github.com/Vilsol/go-mlog/x/impl"
)

var (
	// The build version
	version = "dev"

	// The build commit
	commit = "none"

	// The build date
	date = "unknown"
)

func main() {
	cmd.Execute(version, date, commit)
}
