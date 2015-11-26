package main

// TranslatePositionToSector : translates the co-ordinates of a 2D vector to sector indices location
func TranslatePositionToSector(ed float32, n int, vec Vec2f) (int, int) {
	fn := float32(n)
	x := int((vec.x + ed) / (2 * ed) * fn)
	y := int((vec.y + ed) / (2 * ed) * fn)
	return x, y
}
