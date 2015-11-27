package main

/*ColRGB stores a standard 8-bit per channel Red Green Blue colour representation. */
type ColRGB struct {
	red   byte
	green byte
	blue  byte
}

/*Environment specifies the boundary / dimensions of the working model. They extend in both positive and negative directions, oriented at the center. Setting any field (eg. zBounds) to zero will reduce the dimensionality of the model. For most cases, a 2D environment will be sufficient. */
type Environment struct {
	xBounds float32
	yBounds float32
	zBounds float32
}

// Vec2f – 2D co-ordinates
type Vec2f struct {
	x float32
	y float32
}

// Vec2fCalc - 2D vector arithmetic
type Vec2fCalc interface {
	vecAddition(Vec2f, Vec2f) Vec2f
	vecScalarMultiply(Vec2f, float32) Vec2f
	dotProduct(Vec2f, Vec2f) float32
	magnitude(Vec2f) float32
	angleFromOrigin(Vec2f) float32
	relativeAngle(Vec2f, Vec2f) float32
}

type colourPolymorhicPrey struct {
	pos       Vec2f
	movS      float32
	movA      float32
	heading   float32
	direction Vec2f
	lifetime  int32
	hunger    int32
	fertility int32 //	interval measurement between birth and sex
	gravid    bool  //	i.e. pregnant
	colour    ColRGB
}

type visualPredator struct {
	pos       Vec2f
	movS      float32
	movA      float32
	heading   float32
	direction Vec2f
	lifetime  int32
	hunger    int32
	fertility int32 //	interval measurement between birth and sex
	gravid    bool  //	i.e. pregnant
	visRange  float32
	visAcuity float32
	imprint   ColRGB
}

// AgentActions interface for general agent behaviours
type AgentActions interface {
	UpdateSectorLocation()
	Turn(𝝧 float32)
	Move()
	Death()
}
