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
  8: op add @stack @stack 1
  9: write _main_i bank1 @stack
 10: op add @stack @stack 1
 11: write 13 bank1 @stack
 12: jump 2 always
 13: op sub @stack @stack 2
 14: set _main_0 @return
 15: print _main_0
 16: print "\n"
 17: op add _main_i _main_i 1
 18: jump 8 lessThan _main_i 10`,
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
op add @stack @stack 1
write @counter cell2 0
write @stack cell2 1
write _main_i bank1 @stack
write @counter cell2 0
write @stack cell2 1
op add @stack @stack 1
write @counter cell2 0
write @stack cell2 1
write 35 bank1 @stack
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
jump 22 lessThan _main_i 10`,
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
op add @stack @stack 1                        // Update Stack Pointer
write _main_i bank1 @stack                    // Write argument to memory
op add @stack @stack 1                        // Update Stack Pointer
write 13 bank1 @stack                         // Set Trampoline Address
jump 2 always                                 // Jump to function: foo
op sub @stack @stack 2                        // Update Stack Pointer
set _main_0 @return                           // Set the variable to the value
print _main_0                                 // Call to native function
print "\n"                                    // Call to native function
op add _main_i _main_i 1                      // Execute for loop post condition increment/decrement
jump 8 lessThan _main_i 10                    // Jump to start of loop`,
		},
		{
			name:  "All",
			input: testInput,
			options: transpiler.Options{
				Numbers:  true,
				Comments: true,
				Debug:    true,
			},
			output: `  0: set @stack 0                                  // Reset Stack
  1: jump 19 always                                // Jump to start of main

     // Function: foo //
  2: write @counter cell2 0                        // Debug
  3: write @stack cell2 1                          // Debug
  4: op sub _foo_0 @stack 1                        // Calculate address of parameter
  5: write @counter cell2 0                        // Debug
  6: write @stack cell2 1                          // Debug
  7: read _foo_x bank1 _foo_0                      // Read parameter into variable
  8: write @counter cell2 0                        // Debug
  9: write @stack cell2 1                          // Debug
 10: op add _foo_1 _foo_x 20                       // Execute operation
 11: write @counter cell2 0                        // Debug
 12: write @stack cell2 1                          // Debug
 13: set @return _foo_1                            // Set return data
 14: write @counter cell2 0                        // Debug
 15: write @stack cell2 1                          // Debug
 16: read @counter bank1 @stack                    // Trampoline back

     // Function: main //
 17: write @counter cell2 0                        // Debug
 18: write @stack cell2 1                          // Debug
 19: set _main_i 0                                 // Set the variable to the value
 20: write @counter cell2 0                        // Debug
 21: write @stack cell2 1                          // Debug
 22: op add @stack @stack 1                        // Update Stack Pointer
 23: write @counter cell2 0                        // Debug
 24: write @stack cell2 1                          // Debug
 25: write _main_i bank1 @stack                    // Write argument to memory
 26: write @counter cell2 0                        // Debug
 27: write @stack cell2 1                          // Debug
 28: op add @stack @stack 1                        // Update Stack Pointer
 29: write @counter cell2 0                        // Debug
 30: write @stack cell2 1                          // Debug
 31: write 35 bank1 @stack                         // Set Trampoline Address
 32: write @counter cell2 0                        // Debug
 33: write @stack cell2 1                          // Debug
 34: jump 4 always                                 // Jump to function: foo
 35: write @counter cell2 0                        // Debug
 36: write @stack cell2 1                          // Debug
 37: op sub @stack @stack 2                        // Update Stack Pointer
 38: write @counter cell2 0                        // Debug
 39: write @stack cell2 1                          // Debug
 40: set _main_0 @return                           // Set the variable to the value
 41: write @counter cell2 0                        // Debug
 42: write @stack cell2 1                          // Debug
 43: print _main_0                                 // Call to native function
 44: print "\n"                                    // Call to native function
 45: write @counter cell2 0                        // Debug
 46: write @stack cell2 1                          // Debug
 47: op add _main_i _main_i 1                      // Execute for loop post condition increment/decrement
 48: write @counter cell2 0                        // Debug
 49: write @stack cell2 1                          // Debug
 50: jump 22 lessThan _main_i 10                   // Jump to start of loop`,
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
