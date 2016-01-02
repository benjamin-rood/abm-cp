package colour

import (
	"math"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
)

/*
RGB stores a floating-point representation of 3 colour channel
RGB colour (Red Green Blue), where each value is in range [0,1]
Colours live in a form of vector space.
*/
type RGB struct {
	Red   float64
	Green float64
	Blue  float64
}

// provided global "constants"
// Black White Red Green Blue Yellow Magenta Cyan Orange

var (
	Black   = RGB{Red: 0.0, Green: 0.0, Blue: 0.0}
	White   = RGB{Red: 1.0, Green: 1.0, Blue: 1.0}
	Red     = RGB{Red: 1.0, Green: 0.0, Blue: 0.0}
	Green   = RGB{Red: 0.0, Green: 1.0, Blue: 0.0}
	Blue    = RGB{Red: 0.0, Green: 0.0, Blue: 1.0}
	Yellow  = RGB{Red: 1.0, Green: 1.0, Blue: 0.0}
	Magenta = RGB{Red: 1.0, Green: 0.0, Blue: 1.0}
	Cyan    = RGB{Red: 0.0, Green: 1.0, Blue: 1.0}
	Orange  = RGB{Red: 1.0, Green: 0.5, Blue: 0.0}
)

/*RGBDistance quantifies the value difference between two RGB structs,
returning a floating-point ratio from 0.0 to 1.0.
Multiply the returned value by100 for a percentage.
NOTE: this is a distinct concept from the distance between them as 3D vectors,
as there would be 2 other RGB for any RGB with an identical magnitude.
e.g. [1.0 0 0] [0 1.0 0] [0 0 1.0] will all have the same magnitude, but are
pure Red, pure Blue, pure Green respectively! */
func RGBDistance(c1 RGB, c2 RGB) float64 {
	red := c1.Red - c2.Red
	green := c1.Green - c2.Green
	blue := c1.Blue - c2.Blue
	return calc.ToFixed(((red + blue + green) / 3.0), 3) // returns to 3 d.p. only
}

// RandRGB will return a random valid RGB object within the complete range of all possible RGB values.
func RandRGB() RGB {
	red := calc.RandFloatIn(0, math.Nextafter(1.0, 2.0))
	green := calc.RandFloatIn(0, math.Nextafter(1.0, 2.0))
	blue := calc.RandFloatIn(0, math.Nextafter(1.0, 2.0))
	return RGB{red, green, blue}
}

// RandRGBClamped will return a random valid RGB object within some differential of `col`
func RandRGBClamped(col RGB, diff float64) RGB {
	red := col.Red + calc.RandFloatIn(-diff, diff)
	green := col.Green + calc.RandFloatIn(-diff, diff)
	blue := col.Blue + calc.RandFloatIn(-diff, diff)
	red = calc.ClampFloatIn(red, 0.0, 1.0)
	green = calc.ClampFloatIn(green, 0.0, 1.0)
	blue = calc.ClampFloatIn(blue, 0.0, 1.0)
	return RGB{red, green, blue}
}
