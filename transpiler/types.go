package transpiler

import (
	"context"
	"strconv"
	"strings"
)

func MLOGToString(ctx context.Context, statements [][]Resolvable, statement MLOGAble, lineNumber int, source string) [][]string {
	lines := make([][]string, 0)

	for _, line := range statements {
		currentLine := make([]string, 0)
		resultLine := ""
		for _, t := range line {
			if resultLine != "" {
				resultLine += " "
			}
			resultLine += t.GetValue()
		}

		currentLine = append(currentLine, resultLine)

		if ctx.Value(contextOptions).(Options).Numbers {
			currentLine = append(currentLine, "# "+strconv.Itoa(lineNumber))
		}

		if ctx.Value(contextOptions).(Options).Comments {
			currentLine = append(currentLine, "# "+statement.GetComment(lineNumber))
		}

		if ctx.Value(contextOptions).(Options).Source {
			sourcePos := statement.GetSourcePos(lineNumber)
			if sourcePos != nil {
				currentLine = append(currentLine, "# "+strings.ReplaceAll(source[sourcePos.Pos()-1:sourcePos.End()-1], "\n", ""))
			} else {
				currentLine = append(currentLine, "")
			}
		}

		lines = append(lines, currentLine)
		lineNumber++
	}

	return lines
}
