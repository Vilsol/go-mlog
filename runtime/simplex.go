package runtime

import "math"

var grad3 = [][]int{
	{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
	{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
	{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1},
}

// Translated to Go from https://github.com/Anuken/Arc/blob/f9bf8552f839d87fa561adafa3085d83a044be9f/arc-core/src/arc/util/noise/Simplex.java#L114-L187
// 2D raw Simplex noise
func raw2d(seed int, x float64, y float64) float64 {
	// Noise contributions from the three corners
	var n0, n1, n2 float64

	// Skew the input space to determine which simplex cell we're in
	F2 := 0.5 * (math.Sqrt(3.0) - 1.0)

	// Hairy factor for 2D
	s := (x + y) * F2
	i := math.Floor(x + s)
	j := math.Floor(y + s)

	G2 := (3.0 - math.Sqrt(3.0)) / 6.0
	t := (i + j) * G2

	// Unskew the cell origin back to (x,y) space
	X0 := i - t
	Y0 := j - t
	// The x,y distances from the cell origin
	x0 := x - X0
	y0 := y - Y0

	// For the 2D case, the simplex shape is an equilateral triangle.
	// Determine which simplex we are in.
	var i1, j1 int // Offsets for second (middle) corner of simplex in (i,j) coords
	if x0 > y0 {
		i1 = 1
		j1 = 0
	} else { // lower triangle, XY order: (0,0)->(1,0)->(1,1)
		i1 = 0
		j1 = 1
	} // upper triangle, YX order: (0,0)->(0,1)->(1,1)

	// A step of (1,0) in (i,j) means a step of (1-c,-c) in (x,y), and
	// a step of (0,1) in (i,j) means a step of (-c,1-c) in (x,y), where
	// c = (3-sqrt(3))/6
	x1 := x0 - float64(i1) + G2 // Offsets for middle corner in (x,y) unskewed coords
	y1 := y0 - float64(j1) + G2
	x2 := x0 - 1.0 + 2.0*G2 // Offsets for last corner in (x,y) unskewed coords
	y2 := y0 - 1.0 + 2.0*G2

	// Work out the hashed gradient indices of the three simplex corners
	ii := int(i) & 255
	jj := int(j) & 255
	gi0 := perm(seed, ii+perm(seed, jj)) % 12
	gi1 := perm(seed, ii+i1+perm(seed, jj+j1)) % 12
	gi2 := perm(seed, ii+1+perm(seed, jj+1)) % 12

	// Calculate the contribution from the three corners
	t0 := 0.5 - x0*x0 - y0*y0
	if t0 < 0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * dot(grad3[gi0], x0, y0) // (x,y) of grad3 used for 2D gradient
	}

	t1 := 0.5 - x1*x1 - y1*y1
	if t1 < 0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * dot(grad3[gi1], x1, y1)
	}

	t2 := 0.5 - x2*x2 - y2*y2
	if t2 < 0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * dot(grad3[gi2], x2, y2)
	}

	return 70.0 * (n0 + n1 + n2)
}

//hash function: seed (any) + x (0-512) -> 0-256
func perm(seed int, x int) int {
	r := ((uint(x) >> 16) ^ uint(x)) * uint(0x45d9f3b)
	r = ((r >> 16) ^ r) * (0x45d9f3b + uint(seed))
	r = (r >> 16) ^ r
	return x & 0xff
}

func dot(g []int, x float64, y float64) float64 {
	return float64(g[0])*x + float64(g[1])*y
}
