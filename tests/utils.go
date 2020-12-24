package tests

import (
	"fmt"
	_ "github.com/Vilsol/go-mlog/m"
	_ "github.com/Vilsol/go-mlog/x"
)

func TestMain(main string) string {
	return fmt.Sprintf(`package main

import (
	"github.com/Vilsol/go-mlog/m"
	"github.com/Vilsol/go-mlog/x"
)

func main() {
%s
}`, main)
}
