package runtime

import (
	"fmt"
	"github.com/Vilsol/go-mlog/mlog"
	"github.com/pkg/errors"
	"strings"
)

func Parse(input string) ([]Operation, error) {
	lines, operationLines := mlog.Tokenize(input)

	instructions := make([]Operation, operationLines)
	i := 0
	for _, line := range lines {
		if len(line.Instruction) == 0 {
			continue
		}

		if op, ok := operationRegistry[line.Instruction[0]]; ok {
			var err error
			executor, err := op(line.Instruction[1:])

			if err != nil {
				// TODO Contextual errors
				return nil, errors.Wrap(err, fmt.Sprintf("error on line %d: '%s'", line.SourceLine, strings.Join(line.Instruction, " ")))
			}

			instructions[i] = Operation{
				Line:     line,
				Executor: executor,
			}

			i++
		} else {
			return nil, errors.New("unknown operation: " + line.Instruction[0])
		}
	}

	return instructions, nil
}
