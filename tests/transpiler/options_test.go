package transpiler

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/transpiler"
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
			output: `jump 5 always                 	# 0 	
set _foo_x @funcArg_foo_0     	# 1 	
op add _foo_0 _foo_x 20       	# 2 	
set @return_0 _foo_0          	# 3 	
set @counter @funcTramp_foo   	# 4 	
set _main_i 0                 	# 5 	
op lessThan _main_0 _main_i 10	# 6 	
jump 9 equal _main_0 true     	# 7 	
jump 17 always                	# 8 	
set @funcArg_foo_0 _main_i    	# 9 	
set @funcTramp_foo 12         	# 10	
jump 1 always                 	# 11	
set _main_1 @return_0         	# 12	
print _main_1                 	# 13	
print "\n"                    	# 14	
op add _main_i _main_i 1      	# 15	
jump 6 always                 	# 16	
end                           	# 17	`,
		},
		{
			name:  "Comments",
			input: testInput,
			options: transpiler.Options{
				Comments:      true,
				CommentOffset: 45,
			},
			output: `jump 5 always                 	# Jump to start of main         	
#                             	
# Function: foo #             	
#                             	
set _foo_x @funcArg_foo_0     	# Read parameter into variable  	
op add _foo_0 _foo_x 20       	# Execute operation             	
set @return_0 _foo_0          	# Set return data               	
set @counter @funcTramp_foo   	# Trampoline back               	
#                             	
# Function: main #            	
#                             	
set _main_i 0                 	# Assign value to variable      	
op lessThan _main_0 _main_i 10	# Execute operation             	
jump 9 equal _main_0 true     	# Jump into the loop            	
jump 17 always                	# Jump to end of loop           	
set @funcArg_foo_0 _main_i    	# Set foo argument: 0           	
set @funcTramp_foo 12         	# Set Trampoline Address        	
jump 1 always                 	# Jump to function: foo         	
set _main_1 @return_0         	# Set variable to returned value	
print _main_1                 	# Call to native function       	
print "\n"                    	# Call to native function       	
op add _main_i _main_i 1      	# Execute increment/decrement   	
jump 6 always                 	# Jump to start of loop         	
end                           	# Trampoline back               	`,
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
op add _foo_0 _foo_x 20       	# x + 20         	
set @return_0 _foo_0          	# return x + 20  	
set @counter @funcTramp_foo   	                 	
set _main_i 0                 	# i := 0         	
op lessThan _main_0 _main_i 10	# i < 10         	
jump 9 equal _main_0 true     	                 	
jump 17 always                	                 	
set @funcArg_foo_0 _main_i    	                 	
set @funcTramp_foo 12         	                 	
jump 1 always                 	# foo(i)         	
set _main_1 @return_0         	                 	
print _main_1                 	# println(foo(i))	
print "\n"                    	# println(foo(i))	
op add _main_i _main_i 1      	# i++            	
jump 6 always                 	                 	
end                           	                 	`,
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
			output: `jump 5 always                 	# 0 	# Jump to start of main         	                 	
#                             	
# Function: foo #             	
#                             	
set _foo_x @funcArg_foo_0     	# 1 	# Read parameter into variable  	                 	
op add _foo_0 _foo_x 20       	# 2 	# Execute operation             	# x + 20         	
set @return_0 _foo_0          	# 3 	# Set return data               	# return x + 20  	
set @counter @funcTramp_foo   	# 4 	# Trampoline back               	                 	
#                             	
# Function: main #            	
#                             	
set _main_i 0                 	# 5 	# Assign value to variable      	# i := 0         	
op lessThan _main_0 _main_i 10	# 6 	# Execute operation             	# i < 10         	
jump 9 equal _main_0 true     	# 7 	# Jump into the loop            	                 	
jump 17 always                	# 8 	# Jump to end of loop           	                 	
set @funcArg_foo_0 _main_i    	# 9 	# Set foo argument: 0           	                 	
set @funcTramp_foo 12         	# 10	# Set Trampoline Address        	                 	
jump 1 always                 	# 11	# Jump to function: foo         	# foo(i)         	
set _main_1 @return_0         	# 12	# Set variable to returned value	                 	
print _main_1                 	# 13	# Call to native function       	# println(foo(i))	
print "\n"                    	# 14	# Call to native function       	# println(foo(i))	
op add _main_i _main_i 1      	# 15	# Execute increment/decrement   	# i++            	
jump 6 always                 	# 16	# Jump to start of loop         	                 	
end                           	# 17	# Trampoline back               	                 	`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mlog, err := transpiler.GolangToMLOG(test.input, test.options)

			if err != nil {
				t.Error(err)
				return
			}

			testza.AssertEqual(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
