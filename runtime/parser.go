package runtime

import (
	"github.com/pkg/errors"
	"strings"
)

func Tokenize(input string) ([]MLOGLine, int) {
	count := strings.Count(input, "\n")
	result := make([]MLOGLine, count+1)

	operationLines := 0
	j := 0
	inString := false
	inComment := false
	currentLine := make([]string, 0)
	var currentToken strings.Builder
	for _, c := range input {
		if c == '#' && !inString && !inComment {
			if currentToken.Len() > 0 {
				currentLine = append(currentLine, currentToken.String())
			}

			currentToken.Reset()

			inComment = true
		} else if c == '\n' {
			if !inComment {
				currentLine = append(currentLine, currentToken.String())
			}

			if len(currentLine) > 0 {
				operationLines++
			}

			result[j] = MLOGLine{
				Instruction: currentLine,
			}

			currentLine = make([]string, 0)

			if inComment {
				result[j].Comment = currentToken.String()
			}

			currentToken.Reset()
			inString = false
			inComment = false
			j++
		} else if (c == ' ' || c == '\t') && !inString && !inComment {
			if currentToken.Len() > 0 {
				currentLine = append(currentLine, currentToken.String())
			}
			currentToken.Reset()
		} else if c == '"' {
			currentToken.WriteRune(c)
			inString = !inString
		} else if c == '\r' {
			// Ignored
		} else {
			currentToken.WriteRune(c)
		}
	}

	if currentToken.Len() != 0 {
		currentLine = append(currentLine, currentToken.String())
	}

	result[j] = MLOGLine{
		Instruction: currentLine,
	}

	if len(currentLine) > 0 {
		operationLines++
	}

	if inComment {
		result[j].Comment = currentToken.String()
	}

	return result, operationLines
}

func Parse(input string) ([]OperationExecutor, error) {
	lines, operationLines := Tokenize(input)

	instructions := make([]OperationExecutor, operationLines)
	i := 0
	for _, line := range lines {
		if len(line.Instruction) == 0 {
			continue
		}

		if op, ok := operationRegistry[line.Instruction[0]]; ok {
			var err error
			instructions[i], err = op(line.Instruction[1:])
			if err != nil {
				// TODO Contextual errors
				return nil, errors.Wrap(err, "error on line '"+strings.Join(line.Instruction, " ")+"'")
			}
			i++
		} else {
			return nil, errors.New("unknown operation: " + line.Instruction[0])
		}
	}

	return instructions, nil
}
