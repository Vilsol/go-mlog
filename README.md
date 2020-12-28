# go-mlog

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/vilsol/go-mlog/build)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/vilsol/go-mlog)
[![codecov](https://codecov.io/gh/Vilsol/go-mlog/branch/master/graph/badge.svg?token=LFNKYWS0N2)](https://codecov.io/gh/Vilsol/go-mlog)
[![CodeFactor](https://www.codefactor.io/repository/github/vilsol/go-mlog/badge)](https://www.codefactor.io/repository/github/vilsol/go-mlog)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vilsol/go-mlog)
[![Go Reference](https://pkg.go.dev/badge/github.com/Vilsol/go-mlog.svg)](https://pkg.go.dev/github.com/Vilsol/go-mlog)

Go to MLOG transpiler.

## Web IDE

A Web IDE is available [here](https://vilsol.github.io/go-mlog-web/?1)

## Supports

* Functions
* Multiple function parameters/arguments
* `return` from functions
* `for` loops
* `if`/`else if`/`else` statements
* `switch` statement
* `break`/`continue`/`fallthrough` statements
* Binary and Unary math
* Function level variable scopes
* Contextual errors

## Roadmap

* Full variable block scoping
* `print`/`println` multiple arguments
* Variable argument count functions
* Multiple function return values
* Optimize simple jump instructions
* Multi-pass post-processing

## Design Limitations

* Only hardcoded (translated) imports allowed
* Single file support only
* No recursion ([more info here](RECURSION.md))

## Endgame Roadmap

* Transpilation optimizations
* MLOG Runtime

## Usage

```
Usage:
  go-mlog transpile [flags] <program>

Flags:
  -h, --help   help for transpile

Global Flags:
      --colors          Force log output with colors
      --comments        Output comments
      --debug           Write to debug memory cell
      --log string      The log level to output (default "info")
      --numbers         Output line numbers
      --output string   Output file. Outputs to stdout if unspecified
```