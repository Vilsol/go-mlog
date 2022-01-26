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

## Examples

There are several example programs available on [the wiki](https://github.com/Vilsol/go-mlog/wiki/Examples).

## Supports

* Functions
* Multiple function parameters/arguments
* Multiple function return values
* `return` from functions
* `for` loops
* `if`/`else if`/`else` statements
* `switch` statement
* `break`/`continue`/`fallthrough` statements
* Statement labeling
* `goto` statements
* Binary and Unary math
* Function level variable scopes
* Contextual errors
* Tree-shaking unused functions
* Multi-pass pre/post-processing
* Stackless functions
* Comment generation including source mapping
* Sub-selector support
* Type checking

## In Progress

* MLOG Runtime
* MLOG to Go decompiler

## Roadmap

* Full variable block scoping
* Nested sub-selector support
* Merged compiler and decompiler registries
* Constant string and number slices

## Planned Optimizations

* Simple jump instructions

## Design Limitations

* Only hardcoded (translated) imports allowed
* Single file support only
* No recursion ([more info here](RECURSION.md))

## Endgame Roadmap

* Transpilation optimizations
* Support tail-recursion

## CLI Usage

```
Usage:
  go-mlog [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  decompile   Decompile MLOG to Go
  execute     Execute MLOG
  help        Help about any command
  transpile   Transpile Go to MLOG
  trex        Transpile Go to MLOG and execute it
  typings     Output typings as JSON
  version     Print current go-mlog version

Flags:
      --colors               Force log output with colors
      --comment-offset int   Comment offset from line start (default 60)
      --comments             Output comments
  -h, --help                 help for go-mlog
      --log string           The log level to output (default "info")
      --metrics              Output source metrics after execution
      --numbers              Output line numbers
      --output string        Output file. Outputs to stdout if unspecified
      --source               Output source code after comment
      --stacked string       Use a provided memory cell/bank as a stack
```