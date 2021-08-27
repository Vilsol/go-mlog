package runtime

import "math"

func FastAtan2(x float64, y float64) float64 {
	if math.Abs(x) < 0.0000001 {
		if y > 0 {
			return math.Pi / 2
		}

		if y == 0 {
			return 0
		}

		return (math.Pi * -1) / 2
	}

	z := y / x
	if math.Abs(z) < 1 {
		atan := z / (1 + 0.28*z*z)
		if x < 0 {
			if y < 0 {
				return atan + math.Pi*-1
			}
			return atan + math.Pi
		}
		return atan
	}

	atan := math.Pi/2 - z/(z*z+0.28)

	if y < 0 {
		return atan - math.Pi
	}
	return atan
}
