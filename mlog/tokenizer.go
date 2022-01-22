package mlog

import (
	"strings"
)

type MLOGLine struct {
	Instruction []string
	Comment     string
	SourceLine  int
	Label       string
}

func Tokenize(input string) ([]MLOGLine, int) {
	count := strings.Count(input, "\n")
	result := make([]MLOGLine, count+1)

	operationLines := 0
	j := 0
	inString := false
	inComment := false
	currentLine := make([]string, 0)
	var nextLineLabel string
	var currentToken strings.Builder
	for _, c := range input {
		if c == '#' && !inString && !inComment {
			if currentToken.Len() > 0 {
				currentLine = append(currentLine, currentToken.String())
			}

			currentToken.Reset()

			inComment = true
		} else if c == '\n' {
			if !inComment && currentToken.Len() > 0 {
				currentLine = append(currentLine, currentToken.String())
			}

			if len(currentLine) == 0 && !inComment {
				continue
			}

			operationLines++

			result[j] = MLOGLine{
				Instruction: currentLine,
				SourceLine:  j,
			}

			currentLine = make([]string, 0)

			if inComment {
				result[j].Comment = currentToken.String()
			}

			if nextLineLabel != "" {
				result[j].Label = nextLineLabel
				nextLineLabel = ""
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
		} else if c == ':' && !inString && !inComment && len(currentLine) == 0 {
			nextLineLabel = currentToken.String()
			currentToken.Reset()
		} else if c == '\r' {
			// Ignored
		} else {
			currentToken.WriteRune(c)
		}
	}

	if currentToken.Len() != 0 && !inComment {
		currentLine = append(currentLine, currentToken.String())
	}

	result[j] = MLOGLine{
		Instruction: currentLine,
		SourceLine:  j,
	}

	if inComment {
		result[j].Comment = currentToken.String()
	}

	if len(currentLine) > 0 {
		operationLines++
	}

	return result, operationLines
}
