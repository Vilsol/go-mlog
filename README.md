# go-mlog

Go to MLOG transpiler.

## Web IDE

```
// TODO
```

## Supports

* Functions
* Multiple function parameters/arguments
* `return` from functions
* `for` loops
* `if`/`else if`/`else` statements
* Binary and Unary math
* Function level variable scopes

## Roadmap

* Documentation
* `break` out of loops
* `switch` statement
* Full variable block scoping
* `print`/`println` multiple arguments
* Errors with context
* Extensions
* Variable argument count functions
* Multiple function return values

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