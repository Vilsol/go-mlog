package transpiler

import (
	"testing"
)

func TestDraw(t *testing.T) {
	tests := []Test{
		{
			name:   "DrawClear",
			input:  TestMain(`m.DrawClear(1, 2, 3)`, true, false),
			output: `draw clear 1 2 3`,
		},
		{
			name:   "DrawColor",
			input:  TestMain(`m.DrawColor(1, 2, 3, 4)`, true, false),
			output: `draw color 1 2 3 4`,
		},
		{
			name:   "DrawStroke",
			input:  TestMain(`m.DrawStroke(1)`, true, false),
			output: `draw stroke 1`,
		},
		{
			name:   "DrawLine",
			input:  TestMain(`m.DrawLine(1, 2, 3, 4)`, true, false),
			output: `draw line 1 2 3 4`,
		},
		{
			name:   "DrawRect",
			input:  TestMain(`m.DrawRect(1, 2, 3, 4)`, true, false),
			output: `draw rect 1 2 3 4`,
		},
		{
			name:   "DrawLineRect",
			input:  TestMain(`m.DrawLineRect(1, 2, 3, 4)`, true, false),
			output: `draw lineRect 1 2 3 4`,
		},
		{
			name:   "DrawPoly",
			input:  TestMain(`m.DrawPoly(1, 2, 3, 4, 5)`, true, false),
			output: `draw poly 1 2 3 4 5`,
		},
		{
			name:   "DrawLinePoly",
			input:  TestMain(`m.DrawLinePoly(1, 2, 3, 4, 5)`, true, false),
			output: `draw linePoly 1 2 3 4 5`,
		},
		{
			name:   "DrawTriangle",
			input:  TestMain(`m.DrawTriangle(1, 2, 3, 4, 5, 6)`, true, false),
			output: `draw triangle 1 2 3 4 5 6`,
		},
		{
			name:   "DrawImage",
			input:  TestMain(`m.DrawImage(1, 2, "A", 4, 5)`, true, false),
			output: `draw image 1 2 A 4 5`,
		},
		{
			name:   "DrawFlush",
			input:  TestMain(`m.DrawFlush("display1")`, true, false),
			output: `drawflush display1`,
		},
	}
	RunTests(t, tests)
}
