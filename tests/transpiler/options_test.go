package transpiler

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
			output: `jump 5 always              	# 0 	
set _foo_x @funcArg_foo_0  	# 1 	
op add _foo_0 _foo_x 20    	# 2 	
set @return_0 _foo_0       	# 3 	
set @counter @funcTramp_foo	# 4 	
set _main_i 0              	# 5 	
jump 8 lessThan _main_i 10 	# 6 	
jump 16 always             	# 7 	
set @funcArg_foo_0 _main_i 	# 8 	
set @funcTramp_foo 11      	# 9 	
jump 1 always              	# 10	
set _main_0 @return_0      	# 11	
print _main_0              	# 12	
print "\n"                 	# 13	
op add _main_i _main_i 1   	# 14	
jump 8 lessThan _main_i 10 	# 15	
end                        	# 16	`,
		},
		{
			name:  "Comments",
			input: testInput,
			options: transpiler.Options{
				Comments:      true,
				CommentOffset: 45,
			},
			output: `jump 5 always              	# Jump to start of main         	
#                          	
# Function: foo #          	
#                          	
set _foo_x @funcArg_foo_0  	# Read parameter into variable  	
op add _foo_0 _foo_x 20    	# Execute operation             	
set @return_0 _foo_0       	# Set return data               	
set @counter @funcTramp_foo	# Trampoline back               	
#                          	
# Function: main #         	
#                          	
set _main_i 0              	# Assign value to variable      	
jump 8 lessThan _main_i 10 	# Jump into the loop            	
jump 16 always             	# Jump to end of loop           	
set @funcArg_foo_0 _main_i 	# Set foo argument: 0           	
set @funcTramp_foo 11      	# Set Trampoline Address        	
jump 1 always              	# Jump to function: foo         	
set _main_0 @return_0      	# Set variable to returned value	
print _main_0              	# Call to native function       	
print "\n"                 	# Call to native function       	
op add _main_i _main_i 1   	# Execute increment/decrement   	
jump 8 lessThan _main_i 10 	# Jump to start of loop         	
end                        	# Trampoline back               	`,
		},
		{
			name:  "Comments",
			input: testInput,
			options: transpiler.Options{
				Source:        true,
				CommentOffset: 45,
			},
			output: `jump 5 always              	                 	
set _foo_x @funcArg_foo_0  	                 	
op add _foo_0 _foo_x 20    	# x + 20         	
set @return_0 _foo_0       	# return x + 20  	
set @counter @funcTramp_foo	                 	
set _main_i 0              	# i := 0         	
jump 8 lessThan _main_i 10 	                 	
jump 16 always             	                 	
set @funcArg_foo_0 _main_i 	                 	
set @funcTramp_foo 11      	                 	
jump 1 always              	# foo(i)         	
set _main_0 @return_0      	                 	
print _main_0              	# println(foo(i))	
print "\n"                 	# println(foo(i))	
op add _main_i _main_i 1   	# i++            	
jump 8 lessThan _main_i 10 	                 	
end                        	                 	`,
		},
		{
			name:  "All",
			input: testInput,
			options: transpiler.Options{
				Numbers:       true,
				Comments:      true,
				Source:        true,
				CommentOffset: 45,
			},
			output: `jump 5 always              	# 0 	# Jump to start of main         	                 	
#                          	
# Function: foo #          	
#                          	
set _foo_x @funcArg_foo_0  	# 1 	# Read parameter into variable  	                 	
op add _foo_0 _foo_x 20    	# 2 	# Execute operation             	# x + 20         	
set @return_0 _foo_0       	# 3 	# Set return data               	# return x + 20  	
set @counter @funcTramp_foo	# 4 	# Trampoline back               	                 	
#                          	
# Function: main #         	
#                          	
set _main_i 0              	# 5 	# Assign value to variable      	# i := 0         	
jump 8 lessThan _main_i 10 	# 6 	# Jump into the loop            	                 	
jump 16 always             	# 7 	# Jump to end of loop           	                 	
set @funcArg_foo_0 _main_i 	# 8 	# Set foo argument: 0           	                 	
set @funcTramp_foo 11      	# 9 	# Set Trampoline Address        	                 	
jump 1 always              	# 10	# Jump to function: foo         	# foo(i)         	
set _main_0 @return_0      	# 11	# Set variable to returned value	                 	
print _main_0              	# 12	# Call to native function       	# println(foo(i))	
print "\n"                 	# 13	# Call to native function       	# println(foo(i))	
op add _main_i _main_i 1   	# 14	# Execute increment/decrement   	# i++            	
jump 8 lessThan _main_i 10 	# 15	# Jump to start of loop         	                 	
end                        	# 16	# Trampoline back               	                 	`,
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
