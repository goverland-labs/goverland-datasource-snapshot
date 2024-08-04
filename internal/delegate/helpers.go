package delegate

import (
	"math"
)

func basisPointToPercentage(basisPoint float64) float64 {
	return math.Round(basisPoint) / 100
}
