package main

type rgbColour struct {
	r byte
	g byte
	b byte
}

type vec2D struct {
	x float32
	y float32
}

type colourPolymorhicPrey struct {
	movS      float32
	movA      float32
	lifetime  int32
	fertility int32
	gravid    bool // i.e. pregnant
	colour    rgbColour
	pos       vec2D
}

type visualPredator struct {
	movS      float32
	movA      float32
	hunger    int32
	visRange  float32
	visAcuity float32
	imprint   rgbColour
	pos       vec2D
}
