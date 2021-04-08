package m

// Fill the display with the provided color
func DrawClear(r int, g int, b int) {
}

// Set the drawing color for future statements
func DrawColor(r int, g int, b int, a int) {
}

// Set the line width for future line statements
//
// Affects DrawLine, DrawLineRect and DrawLinePoly
func DrawStroke(width int) {
}

// Draw a line between 2 points
func DrawLine(x1 int, y1 int, x2 int, y2 int) {
}

// Draw a filled rectangle from the provided point with the provided width and height
func DrawRect(x int, y int, width int, height int) {
}

// Draw an outlined rectangle from the provided point with the provided width and height
func DrawLineRect(x int, y int, width int, height int) {
}

// Draw a filled equilateral polygon centered around the provided point
func DrawPoly(x int, y int, sides int, radius float64, rotation float64) {
}

// Draw an outlined equilateral polygon centered around the provided point
func DrawLinePoly(x int, y int, sides int, radius float64, rotation float64) {
}

// Draw a filled triangle between provided points
func DrawTriangle(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) {
}

// Draw provided icon centered around the provided point
func DrawImage(x int, y int, image string, size float64, rotation float64) {
}

// Flush all draw statements to the provided display block
func DrawFlush(targetDisplay string) {
}
