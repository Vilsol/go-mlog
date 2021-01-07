package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestOptions(t *testing.T) {
	const testInput = `package main

func main() {
	for i := 0; i < 10; i++ {
		println(foo(i))
	}
}

func foo(x int) int {
	return x + 20
}`

	tests := []struct {
		name    string
		input   string
		output  string
		options transpiler.Options
	}{
		{
			name:  "Numbers",
			input: testInput,
			options: transpiler.Options{
				Numbers: true,
			},
			output: `  0: jump 5 always
  1: set _foo_x @funcArg_foo_0
  2: op add _foo_0 _foo_x 20
  3: set @return _foo_0
  4: set @counter @funcTramp_foo
  5: set _main_i 0
  6: jump 8 lessThan _main_i 10
  7: jump 16 always
  8: set @funcArg_foo_0 _main_i
  9: set @funcTramp_foo 11
 10: jump 1 always
 11: set _main_0 @return
 12: print _main_0
 13: print "\n"
 14: op add _main_i _main_i 1
 15: jump 8 lessThan _main_i 10`,
		},
		{
			name:  "Comments",
			input: testInput,
			options: transpiler.Options{
				Comments: true,
			},
			output: `jump 5 always                                 // Jump to start of main

     // Function: foo //
set _foo_x @funcArg_foo_0                     // Read parameter into variable
op add _foo_0 _foo_x 20                       // Execute operation
set @return _foo_0                            // Set return data
set @counter @funcTramp_foo                   // Trampoline back

     // Function: main //
set _main_i 0                                 // Set the variable to the value
jump 8 lessThan _main_i 10                    // Jump into the loop
jump 16 always                                // Jump to end of loop
set @funcArg_foo_0 _main_i                    // Set foo argument: 0
set @funcTramp_foo 11                         // Set Trampoline Address
jump 1 always                                 // Jump to function: foo
set _main_0 @return                           // Set variable to returned value
print _main_0                                 // Call to native function
print "\n"                                    // Call to native function
op add _main_i _main_i 1                      // Execute increment/decrement
jump 8 lessThan _main_i 10                    // Jump to start of loop`,
		},
		{
			name:  "All",
			input: testInput,
			options: transpiler.Options{
				Numbers:  true,
				Comments: true,
			},
			output: `  0: jump 5 always                                 // Jump to start of main

     // Function: foo //
  1: set _foo_x @funcArg_foo_0                     // Read parameter into variable
  2: op add _foo_0 _foo_x 20                       // Execute operation
  3: set @return _foo_0                            // Set return data
  4: set @counter @funcTramp_foo                   // Trampoline back

     // Function: main //
  5: set _main_i 0                                 // Set the variable to the value
  6: jump 8 lessThan _main_i 10                    // Jump into the loop
  7: jump 16 always                                // Jump to end of loop
  8: set @funcArg_foo_0 _main_i                    // Set foo argument: 0
  9: set @funcTramp_foo 11                         // Set Trampoline Address
 10: jump 1 always                                 // Jump to function: foo
 11: set _main_0 @return                           // Set variable to returned value
 12: print _main_0                                 // Call to native function
 13: print "\n"                                    // Call to native function
 14: op add _main_i _main_i 1                      // Execute increment/decrement
 15: jump 8 lessThan _main_i 10                    // Jump to start of loop`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mlog, err := transpiler.GolangToMLOG(test.input, test.options)

			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
