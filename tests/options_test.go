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
			output: `  0: set @stack 0
  1: jump 7 always
  2: op sub _foo_0 @stack 1
  3: read _foo_x bank1 _foo_0
  4: op add _foo_1 _foo_x 20
  5: set @return _foo_1
  6: read @counter bank1 @stack
  7: set _main_i 0
  8: jump 20 lessThan _main_i 10
  9: op add @stack @stack 1
 10: write _main_i bank1 @stack
 11: op add @stack @stack 1
 12: write 14 bank1 @stack
 13: jump 2 always
 14: op sub @stack @stack 2
 15: set _main_0 @return
 16: print _main_0
 17: print "\n"
 18: op add _main_i _main_i 1
 19: jump 9 lessThan _main_i 10`,
		},
		{
			name:  "Debug",
			input: testInput,
			options: transpiler.Options{
				Debug: true,
			},
			output: `set @stack 0
jump 19 always
write @counter cell2 0
write @stack cell2 1
op sub _foo_0 @stack 1
write @counter cell2 0
write @stack cell2 1
read _foo_x bank1 _foo_0
write @counter cell2 0
write @stack cell2 1
op add _foo_1 _foo_x 20
write @counter cell2 0
write @stack cell2 1
set @return _foo_1
write @counter cell2 0
write @stack cell2 1
read @counter bank1 @stack
write @counter cell2 0
write @stack cell2 1
set _main_i 0
write @counter cell2 0
write @stack cell2 1
jump 54 lessThan _main_i 10
write @counter cell2 0
write @stack cell2 1
op add @stack @stack 1
write @counter cell2 0
write @stack cell2 1
write _main_i bank1 @stack
write @counter cell2 0
write @stack cell2 1
op add @stack @stack 1
write @counter cell2 0
write @stack cell2 1
write 38 bank1 @stack
write @counter cell2 0
write @stack cell2 1
jump 4 always
write @counter cell2 0
write @stack cell2 1
op sub @stack @stack 2
write @counter cell2 0
write @stack cell2 1
set _main_0 @return
write @counter cell2 0
write @stack cell2 1
print _main_0
print "\n"
write @counter cell2 0
write @stack cell2 1
op add _main_i _main_i 1
write @counter cell2 0
write @stack cell2 1
jump 25 lessThan _main_i 10`,
		},
		{
			name:  "Comments",
			input: testInput,
			options: transpiler.Options{
				Comments: true,
			},
			output: `set @stack 0                                  // Reset Stack
jump 7 always                                 // Jump to start of main

     // Function: foo //
op sub _foo_0 @stack 1                        // Calculate address of parameter
read _foo_x bank1 _foo_0                      // Read parameter into variable
op add _foo_1 _foo_x 20                       // Execute operation
set @return _foo_1                            // Set return data
read @counter bank1 @stack                    // Trampoline back

     // Function: main //
set _main_i 0                                 // Set the variable to the value
jump 20 lessThan _main_i 10                   // Jump to end of loop
op add @stack @stack 1                        // Update Stack Pointer
write _main_i bank1 @stack                    // Write argument to memory
op add @stack @stack 1                        // Update Stack Pointer
write 14 bank1 @stack                         // Set Trampoline Address
jump 2 always                                 // Jump to function: foo
op sub @stack @stack 2                        // Update Stack Pointer
set _main_0 @return                           // Set the variable to the value
print _main_0                                 // Call to native function
print "\n"                                    // Call to native function
op add _main_i _main_i 1                      // Execute for loop post condition increment/decrement
jump 9 lessThan _main_i 10                    // Jump to start of loop`,
		},
		{
			name:  "All",
			input: testInput,
			options: transpiler.Options{
				Numbers:  true,
				Comments: true,
			},
			output: `  0: set @stack 0                                  // Reset Stack
  1: jump 7 always                                 // Jump to start of main

     // Function: foo //
  2: op sub _foo_0 @stack 1                        // Calculate address of parameter
  3: read _foo_x bank1 _foo_0                      // Read parameter into variable
  4: op add _foo_1 _foo_x 20                       // Execute operation
  5: set @return _foo_1                            // Set return data
  6: read @counter bank1 @stack                    // Trampoline back

     // Function: main //
  7: set _main_i 0                                 // Set the variable to the value
  8: jump 20 lessThan _main_i 10                   // Jump to end of loop
  9: op add @stack @stack 1                        // Update Stack Pointer
 10: write _main_i bank1 @stack                    // Write argument to memory
 11: op add @stack @stack 1                        // Update Stack Pointer
 12: write 14 bank1 @stack                         // Set Trampoline Address
 13: jump 2 always                                 // Jump to function: foo
 14: op sub @stack @stack 2                        // Update Stack Pointer
 15: set _main_0 @return                           // Set the variable to the value
 16: print _main_0                                 // Call to native function
 17: print "\n"                                    // Call to native function
 18: op add _main_i _main_i 1                      // Execute for loop post condition increment/decrement
 19: jump 9 lessThan _main_i 10                    // Jump to start of loop`,
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
