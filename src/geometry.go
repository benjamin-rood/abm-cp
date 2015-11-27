package main

import "math"

// UnitAngle will map any floating-point value to its angle on a unit circle.
func UnitAngle(angle float64) float64 {
	twoPi := math.Pi * 2
	return angle - (twoPi * math.Floor(angle/twoPi))
}

// TranslatePositionToSector : translates the co-ordinates of a 2D vector to sector indices location
func TranslatePositionToSector(ed float32, n int, vec Vec2f) (int, int) {
	fn := float32(n)
	x := int((vec.x + ed) / (2 * ed) * fn)
	y := int((vec.y + ed) / (2 * ed) * fn)
	return x, y
}

// Vec2fMagniftude does the classic calculation
func Vec2fMagniftude(v Vec2f) float64 {
	return math.Sqrt(float64(v.x*v.x + v.y*v.y))
}

// Vec2fDistance calculates the distance between two positions
func Vec2fDistance(v1 Vec2f, v2 Vec2f) float64 {
	xDiff := (v1.x - v2.x)
	yDiff := (v1.y - v2.y)
	vd := Vec2f{xDiff, yDiff}
	return Vec2fMagniftude(vd)
}

/*ColourDistance quantifies the value difference between two ColRGB structs,
NOT the difference in magnitude between them as 3D vectors, as there would be 2
other ColRGB for any ColRGB with an identical magnitude.
e.g. [255 0 0] [0 255 0] [0 0 255] will all have the same magnitude, but are
pure Red, pure Blue, pure Green respectively! */
func ColourDistance(c1 ColRGB, c2 ColRGB) float32 {
	redDiff := float32(c1.red-c2.red) / 255
	greenDiff := float32(c1.green-c2.green) / 255
	blueDiff := float32(c1.blue-c1.blue) / 255
	return (redDiff + greenDiff + blueDiff) / 3
	/*	will return a floating-point ratio, not a percentage. Multiply the returned value by 100 for a percentage. */
}
