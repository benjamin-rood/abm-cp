package colour

import "github.com/benjamin-rood/abm-colour-polymorphism/calc"

/*
RGB stores a standard 8-bit per channel Red Green Blue colour
representation. Part of pkg geometry colour lives in a form of
vector space.
*/
type RGB struct {
	Red   byte
	Green byte
	Blue  byte
}

/*RGBDistance quantifies the value difference between two RGB structs,
returning a floating-point ratio from 0.0 to 1.0.
Multiply the returned value by100 for a percentage.
NOTE: this is a distinct concept from the distance between them as 3D vectors,
as there would be 2 other RGB for any RGB with an identical magnitude.
e.g. [255 0 0] [0 255 0] [0 0 255] will all have the same magnitude, but are
pure Red, pure Blue, pure Green respectively! */
func RGBDistance(c1 RGB, c2 RGB) float64 {
	red := float64(c1.Red-c2.Red) / 255
	green := float64(c1.Green-c2.Green) / 255
	blue := float64(c1.Blue-c2.Blue) / 255
	return calc.ToFixed(((red + blue + green) / 3.0), 3) // returns to 3 d.p. only
}
