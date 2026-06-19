package manufacturing

import "math"

// roundMm redondea a milímetros enteros: ≤0,4 baja; ≥0,5 sube.
func roundMm(value float64) float64 {
	return math.Round(value)
}
