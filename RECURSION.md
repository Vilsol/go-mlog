# Issues with recursion

We cannot use variables to store data, as they would be overwritten on the next depth pass of the function. 
That forces us to use memory cells/banks.

Currently, the game does not support writing objects to memory cells/banks.
That means we would be unable to support any functions that return any object types.

It is planned in the future to implement recursion support, but if any function would call any function
that returns an object, it would instantly panic out of the transpiler.