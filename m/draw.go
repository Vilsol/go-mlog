package m

// Fill the display with the provided color
func DrawClear[A integer, B integer, C integer](r A, g B, b C) {
}

// Set the drawing color for future statements
func DrawColor[A integer, B integer, C integer, D integer](r A, g B, b C, a D) {
}

// Set the line width for future line statements
//
// Affects DrawLine, DrawLineRect and DrawLinePoly
func DrawStroke[A integer](width A) {
}

// Draw a line between 2 points
func DrawLine[A integer, B integer, C integer, D integer](x1 A, y1 B, x2 C, y2 D) {
}

// Draw a filled rectangle from the provided point with the provided width and height
func DrawRect[A integer, B integer, C integer, D integer](x A, y B, width C, height D) {
}

// Draw an outlined rectangle from the provided point with the provided width and height
func DrawLineRect[A integer, B integer, C integer, D integer](x A, y B, width C, height D) {
}

// Draw a filled equilateral polygon centered around the provided point
func DrawPoly[A integer, B integer, C integer, D float, E float](x A, y B, sides C, radius D, rotation E) {
}

// Draw an outlined equilateral polygon centered around the provided point
func DrawLinePoly[A integer, B integer, C integer, D float, E float](x A, y B, sides C, radius D, rotation E) {
}

// Draw a filled triangle between provided points
func DrawTriangle[A integer, B integer, C integer, D integer, E integer, F integer](x1 A, y1 B, x2 C, y2 D, x3 E, y3 F) {
}

// Draw provided icon centered around the provided point
func DrawImage[A integer, B integer, C float, D float](x A, y B, image string, size C, rotation D) {
}

// Flush all draw statements to the provided display block
func DrawFlush(targetDisplay string) {
}
