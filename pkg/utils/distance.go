package utils

import "math"

func CalculateDistance(x1, y1, x2, y2 float64) float64 {
	dis := math.Abs(x2-x1) + math.Abs(y2-y1)
	return dis
}

func CalculatePrice(distance float64, isSurge bool) float64 {
	price := float64(0)
	if distance <= 10 {
		price = distance * 10
	} else {
		price = 100 + ((distance - 10) * 7)
	}

	if isSurge {
		price *= 1.2
	}

	return price
}
